package ui

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/idongju/t9s/internal/config"
	"github.com/rivo/tview"
)

// App represents the main application
type App struct {
	tviewApp    *tview.Application
	tree        *tview.TreeView
	contentView *tview.TextView
	statusBar   *tview.TextView
	pages       *tview.Pages
	currentDir  string
	currentFile string
	config      *config.Config
}

// NewApp creates a new T9s application
func NewApp() *App {
	// Set global tview theme to use black background
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack
	tview.Styles.ContrastBackgroundColor = tcell.ColorBlack
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorBlack
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.BorderColor = tcell.NewRGBColor(0, 255, 255) // Cyan
	tview.Styles.TitleColor = tcell.NewRGBColor(255, 215, 0)  // Yellow

	// Load config
	cfg, err := config.Load()
	if err != nil {
		// Fallback to current directory
		currentDir, _ := os.Getwd()
		if currentDir == "" {
			currentDir = "."
		}
		cfg = &config.Config{
			TerraformRoot: currentDir,
			Commands: config.CommandsConfig{
				InitTemplate:    "terraform init -backend-config={initconf}",
				PlanTemplate:    "terraform plan -var-file={varfile}",
				ApplyTemplate:   "terraform apply -var-file={varfile}",
				DestroyTemplate: "terraform destroy -var-file={varfile}",
				TfvarsFile:      "config/env.tfvars",
				InitConfFile:    "config/env.conf",
			},
		}
	}

	// Use TerraformRoot from config as the starting directory
	currentDir := cfg.TerraformRoot
	if currentDir == "" {
		currentDir, _ = os.Getwd()
		if currentDir == "" {
			currentDir = "."
		}
	}

	app := &App{
		tviewApp:   tview.NewApplication(),
		currentDir: currentDir,
		config:     cfg,
		pages:      tview.NewPages(),
	}

	app.setupUI()
	app.setupKeyBindings()

	return app
}

// setupUI initializes the UI components
func (a *App) setupUI() {
	// Create header
	header := a.createHeader()

	// Create tree view (left panel)
	a.tree = a.createTreeView()

	// Create content view (right panel)
	a.contentView = a.createContentView()

	// Create status bar
	a.statusBar = a.createStatusBar()

	// Main content layout
	mainFlex := tview.NewFlex().
		AddItem(a.tree, 0, 1, true).
		AddItem(a.contentView, 0, 2, false)
	mainFlex.SetBackgroundColor(tcell.ColorBlack)

	// Root layout (main page)
	mainPage := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 6, 0, false).
		AddItem(mainFlex, 0, 1, true).
		AddItem(a.statusBar, 1, 0, false)
	mainPage.SetBackgroundColor(tcell.ColorBlack)

	// Add pages
	a.pages.AddPage("main", mainPage, true, true)

	a.tviewApp.SetRoot(a.pages, true)
}

// createHeader creates the application header
func (a *App) createHeader() *tview.Flex {
	// Left side: Info and Shortcuts
	infoText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	infoText.SetBackgroundColor(tcell.ColorBlack)

	// Get user and hostname
	user := os.Getenv("USER")
	host, _ := os.Hostname()

	// Get terraform workspace
	workspace := "default"
	cmd := exec.Command("terraform", "workspace", "show")
	if out, err := cmd.Output(); err == nil {
		workspace = strings.TrimSpace(string(out))
	}

	// Info section
	fmt.Fprintf(infoText, "[cyan]Context:[white]  %s\n", workspace)
	fmt.Fprintf(infoText, "[cyan]Path:[white]     %s\n", a.currentDir)
	fmt.Fprintf(infoText, "[cyan]User:[white]     %s@%s\n", user, host)
	fmt.Fprintf(infoText, "[cyan]Version:[white]  v0.1.0\n")

	// Shortcuts section (2 columns)
	shortcuts := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	shortcuts.SetBackgroundColor(tcell.ColorBlack)

	fmt.Fprintf(shortcuts, "[yellow]<s>[white] Settings    [yellow]<p>[white] Plan      [yellow]<?>[white] Help\n")
	fmt.Fprintf(shortcuts, "[yellow]<h>[white] History     [yellow]<a>[white] Apply     [yellow]</>[white] Command\n")
	fmt.Fprintf(shortcuts, "[yellow]<e>[white] Edit        [yellow]<C>[white] Home      [yellow]<q>[white] Quit\n")
	fmt.Fprintf(shortcuts, "[yellow]<Enter>[white] Select")

	// Right side: Logo
	logo := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)
	logo.SetBackgroundColor(tcell.ColorBlack)

	// ASCII Art Logo
	fmt.Fprintf(logo, "[cyan] _______ ___      [orange] ______[white]\n")
	fmt.Fprintf(logo, "[cyan]|_   _/ _ \\___    [orange]/ ____|[white]\n")
	fmt.Fprintf(logo, "[cyan]  | | \\_, /(_-<   [orange]`--. \\ [white]\n")
	fmt.Fprintf(logo, "[cyan]  |_|  /_//__ /   [orange]/\\__/ /[white]\n")
	fmt.Fprintf(logo, "[cyan]                 [orange]\\____/ [white]")

	// Combine layouts
	leftFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(infoText, 40, 0, false).
		AddItem(shortcuts, 0, 1, false)
	leftFlex.SetBackgroundColor(tcell.ColorBlack)

	headerFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftFlex, 0, 1, false).
		AddItem(logo, 26, 0, false)
	headerFlex.SetBackgroundColor(tcell.ColorBlack)

	return headerFlex
}

// createTreeView creates the file tree view
func (a *App) createTreeView() *tview.TreeView {
	rootDir := a.currentDir
	root := tview.NewTreeNode(filepath.Base(rootDir)).
		SetColor(tcell.NewRGBColor(255, 215, 0)). // Yellow
		SetReference(rootDir).
		SetSelectable(true)

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	tree.SetBackgroundColor(tcell.ColorBlack)
	tree.SetBorderColor(tcell.NewRGBColor(0, 255, 255)) // Cyan
	tree.SetTitle(" 📂 File Tree ").SetBorder(true)
	tree.SetGraphicsColor(tcell.NewRGBColor(0, 255, 255)) // Cyan

	// Add initial children
	a.addTreeChildren(root, rootDir)

	// Handle selection
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return
		}

		path := reference.(string)
		children := node.GetChildren()

		if len(children) == 0 {
			// If it's a file, display its content
			info, err := os.Stat(path)
			if err == nil && !info.IsDir() {
				a.displayFile(path)
			} else if err == nil && info.IsDir() {
				// If it's a directory, expand it
				a.addTreeChildren(node, path)
			}
		} else {
			// Collapse if already expanded
			node.SetChildren([]*tview.TreeNode{})
		}
	})

	// Handle changed selection
	tree.SetChangedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return
		}
		path := reference.(string)
		a.updateStatusBar(path)
	})

	return tree
}

// addTreeChildren adds children to a tree node
func (a *App) addTreeChildren(target *tview.TreeNode, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, file := range files {
		// Skip hidden files except for .terraform
		if strings.HasPrefix(file.Name(), ".") && file.Name() != ".terraform" {
			continue
		}

		fullPath := filepath.Join(path, file.Name())
		node := tview.NewTreeNode(file.Name()).
			SetReference(fullPath).
			SetSelectable(true)

		if file.IsDir() {
			node.SetColor(tcell.NewRGBColor(100, 200, 255)) // Light blue for directories
		} else if strings.HasSuffix(file.Name(), ".tfvars") {
			node.SetColor(tcell.NewRGBColor(255, 100, 255)) // Magenta for tfvars
		} else if strings.HasSuffix(file.Name(), ".tf") {
			node.SetColor(tcell.NewRGBColor(100, 255, 100)) // Green for .tf files
		} else {
			node.SetColor(tcell.NewRGBColor(200, 200, 200)) // Gray for other files
		}

		target.AddChild(node)
	}
}

// createContentView creates the content display area
func (a *App) createContentView() *tview.TextView {
	contentView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)

	contentView.SetBackgroundColor(tcell.ColorBlack)
	contentView.SetTextColor(tcell.NewRGBColor(255, 255, 255))
	contentView.SetBorderColor(tcell.NewRGBColor(0, 255, 255)) // Cyan
	contentView.SetTitle(" 📄 Content ").SetBorder(true)

	fmt.Fprintf(contentView, "[yellow]Welcome to T9s![white]\n\n")
	fmt.Fprintf(contentView, "Select a file from the tree to view its content.\n\n")
	fmt.Fprintf(contentView, "[cyan]Available Commands:[white]\n")
	fmt.Fprintf(contentView, "  • [green]h[white] - View terraform history\n")
	fmt.Fprintf(contentView, "  • [green]H[white] - View helm list (helm list -A)\n")
	fmt.Fprintf(contentView, "  • [green]p[white] - Terraform plan\n")
	fmt.Fprintf(contentView, "  • [green]a[white] - Terraform apply\n")
	fmt.Fprintf(contentView, "  • [green]e[white] - Edit current file\n")
	fmt.Fprintf(contentView, "  • [green]q[white] - Quit\n")

	return contentView
}

// createStatusBar creates the status bar
func (a *App) createStatusBar() *tview.TextView {
	statusBar := tview.NewTextView().
		SetDynamicColors(true)

	statusBar.SetBackgroundColor(tcell.ColorBlack)
	statusBar.SetTextColor(tcell.ColorWhite)

	// Simple one-line status bar
	fmt.Fprintf(statusBar, "[yellow]<↑↓>[white] Navigate  [yellow]<Enter>[white] Expand/View  [yellow]<q>[white] Quit")

	return statusBar
}

// displayFile displays the content of a file
func (a *App) displayFile(path string) {
	a.currentFile = path
	content, err := ioutil.ReadFile(path)
	if err != nil {
		a.contentView.Clear()
		fmt.Fprintf(a.contentView, "[red]Error reading file: %v[white]", err)
		return
	}

	a.contentView.Clear()
	a.contentView.SetTitle(fmt.Sprintf(" 📄 %s ", filepath.Base(path)))

	// Check if it's a tfvars file
	if strings.HasSuffix(path, ".tfvars") {
		fmt.Fprintf(a.contentView, "[yellow]File:[white] %s\n", path)
		fmt.Fprintf(a.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))
		fmt.Fprintf(a.contentView, "%s", string(content))
		fmt.Fprintf(a.contentView, "\n\n[gray]Press 'e' to edit this file[white]")
	} else {
		fmt.Fprintf(a.contentView, "[yellow]File:[white] %s\n", path)
		fmt.Fprintf(a.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))
		fmt.Fprintf(a.contentView, "%s", string(content))
	}
}

// updateStatusBar updates the status bar with current path
func (a *App) updateStatusBar(path string) {
	relPath, _ := filepath.Rel(a.currentDir, path)
	a.statusBar.Clear()
	fmt.Fprintf(a.statusBar, "[cyan]═══════════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(a.statusBar, "[yellow]Current:[white] %s  [yellow]|[white]  [yellow]h[white] History  [yellow]H[white] Helm  [yellow]e[white] Edit  [yellow]q[white] Quit\n", relPath)
	fmt.Fprintf(a.statusBar, "[cyan]═══════════════════════════════════════════════════════════════════")
}

// setupKeyBindings sets up global key bindings
func (a *App) setupKeyBindings() {
	a.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Get current page name
		currentPage, _ := a.pages.GetFrontPage()

		// Only handle global shortcuts when on main page
		if currentPage != "main" {
			return event
		}

		switch event.Key() {
		case tcell.KeyCtrlC:
			a.tviewApp.Stop()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				a.tviewApp.Stop()
				return nil
			case 's', 'S':
				a.showSettings()
				return nil
			case 'h':
				a.showHistory()
				return nil
			case 'H':
				a.showHelm()
				return nil
			case 'e', 'E':
				a.editCurrentFile()
				return nil
			case 'p':
				a.showPlan()
				return nil
			case 'a':
				a.showApply()
				return nil
			}
		}
		return event
	})
}

// showHelm shows helm list -A output
func (a *App) showHelm() {
	a.contentView.Clear()
	a.contentView.SetTitle(" ⎈ Helm Releases ")
	fmt.Fprintf(a.contentView, "[yellow]Helm List -A[white]\n\n")
	fmt.Fprintf(a.contentView, "[cyan]Executing:[white] helm list -A\n")
	fmt.Fprintf(a.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))

	cmd := exec.Command("helm", "list", "-A")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(a.contentView, "[red]Error executing helm:[white] %v\n\n", err)
		fmt.Fprintf(a.contentView, "Make sure 'helm' is installed and available in your PATH.\n")
		if len(output) > 0 {
			fmt.Fprintf(a.contentView, "\n[red]Output:[white]\n%s", string(output))
		}
		return
	}

	fmt.Fprintf(a.contentView, "%s", string(output))
}

// showHistory shows terraform history based on state file
func (a *App) showHistory() {
	node := a.tree.GetCurrentNode()
	if node == nil {
		return
	}
	path := node.GetReference().(string)

	// If file is selected, get parent dir
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		path = filepath.Dir(path)
	}

	a.contentView.Clear()
	a.contentView.SetTitle(" ⏰ Terraform History ")
	fmt.Fprintf(a.contentView, "[yellow]Terraform History[white]\n")
	fmt.Fprintf(a.contentView, "[cyan]Directory:[white] %s\n", path)
	fmt.Fprintf(a.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))

	// Check terraform.tfstate
	statePath := filepath.Join(path, "terraform.tfstate")
	stateInfo, err := os.Stat(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(a.contentView, "[yellow]No local terraform.tfstate found.[white]\n")
			fmt.Fprintf(a.contentView, "This might be using a remote backend or hasn't been initialized.\n\n")
		} else {
			fmt.Fprintf(a.contentView, "[red]Error checking state file:[white] %v\n", err)
		}
	} else {
		fmt.Fprintf(a.contentView, "[green]Last Local Apply:[white] %s\n", stateInfo.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Fprintf(a.contentView, "[gray](Based on terraform.tfstate modification time)[white]\n\n")
	}

	// Run terraform show
	fmt.Fprintf(a.contentView, "[cyan]Current State Summary (terraform show):[white]\n")
	go func() {
		cmd := exec.Command("terraform", "show", "-no-color")
		cmd.Dir = path
		output, err := cmd.CombinedOutput()

		a.tviewApp.QueueUpdateDraw(func() {
			if err != nil {
				fmt.Fprintf(a.contentView, "[red]Error executing terraform show:[white] %v\n", err)
			}
			fmt.Fprintf(a.contentView, "%s", string(output))
		})
	}()
}

// showPlan executes terraform plan
func (a *App) showPlan() {
	a.runTerraformCommand("Plan", a.config.Commands.PlanTemplate)
}

// showApply executes terraform apply
func (a *App) showApply() {
	modal := tview.NewModal().
		SetText("Are you sure you want to execute terraform apply?").
		AddButtons([]string{"Yes", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("confirm_apply")
			if buttonLabel == "Yes" {
				// Run apply if confirmed
				a.runTerraformCommand("Apply", a.config.Commands.ApplyTemplate)
			} else {
				// Focus back to tree if cancelled
				a.tviewApp.SetFocus(a.tree)
			}
		})

	// Style the modal
	modal.SetBackgroundColor(tcell.ColorBlack)
	modal.SetTextColor(tcell.ColorWhite)
	modal.SetButtonBackgroundColor(tcell.NewRGBColor(50, 50, 50))
	modal.SetButtonTextColor(tcell.ColorWhite)

	a.pages.AddPage("confirm_apply", modal, true, true)
}

// runTerraformCommand executes a terraform command using the template
func (a *App) runTerraformCommand(action string, template string) {
	node := a.tree.GetCurrentNode()
	if node == nil {
		return
	}
	path := node.GetReference().(string)

	var workDir string
	var varFile string

	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		// If file selected
		workDir = filepath.Dir(path)
		if strings.HasSuffix(path, ".tfvars") {
			varFile = path
		}
	} else {
		// If directory selected
		workDir = path
		// Try to find default var file
		defaultVar := filepath.Join(workDir, a.config.Commands.TfvarsFile)
		if _, err := os.Stat(defaultVar); err == nil {
			varFile = defaultVar
		}
	}

	a.contentView.Clear()
	a.contentView.SetTitle(fmt.Sprintf(" 🚀 Terraform %s ", action))
	fmt.Fprintf(a.contentView, "[yellow]Executing Terraform %s[white]\n", action)
	fmt.Fprintf(a.contentView, "[cyan]Directory:[white] %s\n", workDir)
	if varFile != "" {
		fmt.Fprintf(a.contentView, "[cyan]Var File:[white] %s\n", varFile)
	} else {
		fmt.Fprintf(a.contentView, "[yellow]Warning:[white] No .tfvars file found or selected. Running without -var-file.\n")
	}
	fmt.Fprintf(a.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))

	// Prepare command
	cmdStr := template
	if varFile != "" {
		cmdStr = strings.ReplaceAll(cmdStr, "{varfile}", varFile)
	} else {
		// If no var file, remove the argument part if it exists
		// This is a simple replacement, might be fragile
		cmdStr = strings.ReplaceAll(cmdStr, "-var-file={varfile}", "")
	}

	fmt.Fprintf(a.contentView, "[gray]Command: %s[white]\n\n", cmdStr)

	// Execute
	go func() {
		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			return
		}

		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Dir = workDir

		// Capture output
		output, err := cmd.CombinedOutput()

		a.tviewApp.QueueUpdateDraw(func() {
			if err != nil {
				fmt.Fprintf(a.contentView, "[red]Error executing command:[white] %v\n", err)
			}
			fmt.Fprintf(a.contentView, "%s", string(output))
			fmt.Fprintf(a.contentView, "\n\n[green]Done.[white]")
		})
	}()
}

// editCurrentFile opens the current file in an editor
// editCurrentFile opens the current file in an editor
func (a *App) editCurrentFile() {
	if a.currentFile == "" {
		a.contentView.Clear()
		fmt.Fprintf(a.contentView, "[yellow]No file selected[white]\n\n")
		fmt.Fprintf(a.contentView, "Please select a file from the tree first.")
		return
	}

	// Get editor from env or default to vim
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	// Suspend tview to run external command
	a.tviewApp.Suspend(func() {
		cmd := exec.Command(editor, a.currentFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running editor: %v\nPress Enter to continue...", err)
			fmt.Scanln()
		}
	})

	// Refresh file content after editing
	a.displayFile(a.currentFile)
}

// Run starts the application
func (a *App) Run() error {
	return a.tviewApp.Run()
}

// showSettings displays the settings page
func (a *App) showSettings() {
	form := tview.NewForm().
		AddInputField("Terraform Root Directory", a.config.TerraformRoot, 60, nil, func(text string) {
			a.config.TerraformRoot = text
		}).
		AddInputField("Terraform Plan Template", a.config.Commands.PlanTemplate, 60, nil, func(text string) {
			a.config.Commands.PlanTemplate = text
		}).
		AddInputField("Terraform Apply Template", a.config.Commands.ApplyTemplate, 60, nil, func(text string) {
			a.config.Commands.ApplyTemplate = text
		}).
		AddInputField("Default tfvars File", a.config.Commands.TfvarsFile, 60, nil, func(text string) {
			a.config.Commands.TfvarsFile = text
		}).
		AddButton("Save", func() {
			if err := a.config.Save(); err != nil {
				a.contentView.Clear()
				fmt.Fprintf(a.contentView, "[red]Error saving config: %v[white]", err)
			} else {
				// If TerraformRoot changed, rebuild the UI
				if a.config.TerraformRoot != a.currentDir {
					a.currentDir = a.config.TerraformRoot
					a.setupUI()
				}
				a.pages.SwitchToPage("main")
				a.tviewApp.SetFocus(a.tree)
			}
		}).
		AddButton("Cancel", func() {
			// Reload config to discard changes
			cfg, _ := config.Load()
			if cfg != nil {
				a.config = cfg
			}
			a.pages.SwitchToPage("main")
			a.tviewApp.SetFocus(a.tree)
		})

	form.SetBackgroundColor(tcell.ColorBlack)
	form.SetBorderColor(tcell.NewRGBColor(0, 255, 255)) // Cyan
	form.SetTitle(" ⚙️  Settings ").SetBorder(true)
	form.SetButtonBackgroundColor(tcell.NewRGBColor(50, 50, 50))
	form.SetButtonTextColor(tcell.NewRGBColor(255, 255, 255))

	// Settings header
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	header.SetBackgroundColor(tcell.ColorBlack)
	fmt.Fprintf(header, "[::b][cyan]╔═══════════════════════════════════════════════════════════════════╗\n")
	fmt.Fprintf(header, "[::b][cyan]║[white]  T9s - Settings  [cyan]║\n")
	fmt.Fprintf(header, "[::b][cyan]╚═══════════════════════════════════════════════════════════════════╝")

	help := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	help.SetBackgroundColor(tcell.ColorBlack)
	fmt.Fprintf(help, "\n[yellow]Terraform Root Directory:[white]\n")
	fmt.Fprintf(help, "  Directory where your Terraform code is located (e.g., /home/user/terraform)\n")
	fmt.Fprintf(help, "\n[yellow]Template Variables:[white]\n")
	fmt.Fprintf(help, "  [cyan]{varfile}[white] - Will be replaced with the var file path\n")
	fmt.Fprintf(help, "\n[yellow]Examples:[white]\n")
	fmt.Fprintf(help, "  [green]terraform plan -var-file={varfile}[white]\n")
	fmt.Fprintf(help, "  [green]terraform apply -var-file={varfile} -auto-approve[white]\n")

	settingsLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(form, 0, 1, true).
		AddItem(help, 9, 0, false)
	settingsLayout.SetBackgroundColor(tcell.ColorBlack)

	a.pages.AddPage("settings", settingsLayout, true, false)
	a.pages.SwitchToPage("settings")
	a.tviewApp.SetFocus(form)
}
