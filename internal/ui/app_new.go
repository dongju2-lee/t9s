package ui

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/idongju/t9s/internal/config"
	"github.com/idongju/t9s/internal/db"
	"github.com/idongju/t9s/internal/ui/components"
	"github.com/idongju/t9s/internal/ui/dialog"
	"github.com/idongju/t9s/internal/view"
)

// AppNew represents the main application with new structure
type AppNew struct {
	tviewApp    *tview.Application
	pages       *tview.Pages
	
	// Views
	headerView   *view.HeaderView
	treeView     *view.TreeView
	contentView  *view.ContentView
	statusBar    *view.StatusBar
	helpView     *view.HelpView
	commandView  *view.CommandView
	historyView  *view.HistoryView
	
	// Components
	executor    *components.CommandExecutor
	historyDB   *db.HistoryDB
	
	// State
	currentDir  string
	currentFile string
	config      *config.Config
	focusOnTree bool // true if tree is focused, false if content is focused
}

// NewAppNew creates a new T9s application with improved structure
func NewAppNew() *AppNew {
	// Set global tview theme
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack
	tview.Styles.ContrastBackgroundColor = tcell.ColorBlack
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorBlack
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.BorderColor = tcell.NewRGBColor(0, 255, 255)
	tview.Styles.TitleColor = tcell.NewRGBColor(255, 215, 0)

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

	// Initialize history DB
	historyDB, err := db.NewHistoryDB()
	if err != nil {
		// Fallback if DB fails - app still works
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize history DB: %v\n", err)
	}

	app := &AppNew{
		tviewApp:   tview.NewApplication(),
		currentDir: currentDir,
		config:     cfg,
		pages:      tview.NewPages(),
		historyDB:  historyDB,
		focusOnTree: true, // Start with tree focused
	}

	app.setupViews()
	app.setupKeyBindings()

	return app
}

// setupViews initializes all views
func (a *AppNew) setupViews() {
	// Create views
	a.headerView = view.NewHeaderView(a.currentDir)
	a.treeView = view.NewTreeView(a.currentDir)
	a.contentView = view.NewContentView()
	a.statusBar = view.NewStatusBar(a.currentDir)
	
	// Create executor
	a.executor = components.NewCommandExecutor(a.tviewApp, a.contentView, a.config, a.historyDB)
	
	// Setup tree view handler
	a.treeView.SetFileSelectHandler(func(path string) {
		a.currentFile = path
		a.contentView.DisplayFile(path)
	})
	
	// Setup tree view change handler
	a.treeView.SetChangedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference != nil {
			path := reference.(string)
			a.statusBar.UpdatePath(path)
			
			// Check if it's a directory
			info, err := os.Stat(path)
			if err == nil && info.IsDir() {
				// Show welcome screen when moving to a directory
				a.contentView.ShowWelcome()
			}
		}
	})

	// Main content layout
	mainFlex := tview.NewFlex().
		AddItem(a.treeView, 0, 1, true).
		AddItem(a.contentView, 0, 2, false)
	mainFlex.SetBackgroundColor(tcell.ColorBlack)

	// Root layout
	mainPage := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.headerView, 6, 0, false).
		AddItem(mainFlex, 0, 1, true).
		AddItem(a.statusBar, 1, 0, false)
	mainPage.SetBackgroundColor(tcell.ColorBlack)

	a.pages.AddPage("main", mainPage, true, true)
	a.tviewApp.SetRoot(a.pages, true)
}

// setupKeyBindings sets up global key bindings
func (a *AppNew) setupKeyBindings() {
	a.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Get current page name
		currentPage, _ := a.pages.GetFrontPage()
		
		// Only handle global shortcuts when on main page
		if currentPage != "main" {
			return event
		}
		
		switch event.Key() {
		case tcell.KeyTab:
			// Toggle focus between tree and content view
			a.focusOnTree = !a.focusOnTree
			if a.focusOnTree {
				a.tviewApp.SetFocus(a.treeView)
				a.statusBar.SetFocusIndicator("File Tree")
			} else {
				a.tviewApp.SetFocus(a.contentView)
				a.statusBar.SetFocusIndicator("Content View")
			}
			return nil
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
				path := a.treeView.GetCurrentPath()
				if path != "" {
					a.showHistory(path)
				}
				return nil
			case 'H':
				// Shift+H: Show help
				a.showHelp()
				return nil
			case '?':
				// ?: Also show help
				a.showHelp()
				return nil
			case '/':
				// /: Command mode
				a.showCommandInput()
				return nil
			case 'C':
				// Shift+C: Show available commands
				a.contentView.ShowWelcome()
				return nil
			case 'e', 'E':
				if a.currentFile != "" {
					a.executor.EditFile(a.currentFile)
					a.contentView.DisplayFile(a.currentFile)
				}
				return nil
			case 'i', 'I':
				// i: Terraform init
				path := a.treeView.GetCurrentPath()
				if path != "" {
					a.showInitConfirmation(path)
				}
				return nil
			case 'p':
				// p: Terraform plan
				path := a.treeView.GetCurrentPath()
				if path != "" {
					a.showPlanConfirmation(path)
				}
				return nil
			case 'a':
				// a: Terraform apply
				a.showApplyConfirmation()
				return nil
			case 'd', 'D':
				// d: Terraform destroy
				path := a.treeView.GetCurrentPath()
				if path != "" {
					a.showDestroyConfirmation(path)
				}
				return nil
			}
		}
		return event
	})
}

// showSettings displays the settings dialog
func (a *AppNew) showSettings() {
	settingsDialog := dialog.NewSettingsDialog(
		a.config,
		func() {
			// Check if TerraformRoot changed
			if a.config.TerraformRoot != a.currentDir {
				// Rebuild tree view with new root directory
				a.currentDir = a.config.TerraformRoot
				a.treeView = view.NewTreeView(a.currentDir)
				a.treeView.SetFileSelectHandler(func(path string) {
					a.currentFile = path
					a.contentView.DisplayFile(path)
				})
				a.treeView.SetChangedFunc(func(node *tview.TreeNode) {
					reference := node.GetReference()
					if reference != nil {
						a.statusBar.UpdatePath(reference.(string))
					}
				})
				
				// Rebuild header and status bar
				a.headerView = view.NewHeaderView(a.currentDir)
				a.statusBar = view.NewStatusBar(a.currentDir)
				
				// Rebuild main layout
				a.rebuildMainPage()
			}
			a.pages.SwitchToPage("main")
			a.tviewApp.SetFocus(a.treeView)
		},
		func() {
			cfg, _ := config.Load()
			if cfg != nil {
				a.config = cfg
			}
			a.pages.SwitchToPage("main")
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("settings", settingsDialog, true, false)
	a.pages.SwitchToPage("settings")
	a.tviewApp.SetFocus(settingsDialog.GetForm())
}

// rebuildMainPage rebuilds the main page with updated views
func (a *AppNew) rebuildMainPage() {
	// Main content layout
	mainFlex := tview.NewFlex().
		AddItem(a.treeView, 0, 1, true).
		AddItem(a.contentView, 0, 2, false)
	mainFlex.SetBackgroundColor(tcell.ColorBlack)

	// Root layout
	mainPage := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.headerView, 6, 0, false).
		AddItem(mainFlex, 0, 1, true).
		AddItem(a.statusBar, 1, 0, false)
	mainPage.SetBackgroundColor(tcell.ColorBlack)

	a.pages.RemovePage("main")
	a.pages.AddPage("main", mainPage, true, true)
}

// showInitConfirmation shows confirmation dialog for terraform init
func (a *AppNew) showInitConfirmation(path string) {
	info := components.GetTerraformCommandInfo(path, a.config.Commands.InitTemplate, a.config.Commands.InitConfFile, a.config)
	
	confirmDialog := dialog.NewTerraformConfirmDialog(
		"terraform init",
		info.WorkDir,
		info.ConfigFile,
		info.Content,
		func() {
			a.pages.RemovePage("confirm_tf")
			a.executeTerraformCommand("Init", info.WorkDir, info.Command)
		},
		func() {
			a.pages.RemovePage("confirm_tf")
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("confirm_tf", confirmDialog, true, true)
	if form := confirmDialog.GetForm(); form != nil {
		a.tviewApp.SetFocus(form)
	}
}

// showPlanConfirmation shows confirmation dialog for terraform plan
func (a *AppNew) showPlanConfirmation(path string) {
	info := components.GetTerraformCommandInfo(path, a.config.Commands.PlanTemplate, a.config.Commands.TfvarsFile, a.config)
	
	confirmDialog := dialog.NewTerraformConfirmDialog(
		"terraform plan",
		info.WorkDir,
		info.ConfigFile,
		info.Content,
		func() {
			a.pages.RemovePage("confirm_tf")
			a.executeTerraformCommand("Plan", info.WorkDir, info.Command)
		},
		func() {
			a.pages.RemovePage("confirm_tf")
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("confirm_tf", confirmDialog, true, true)
	if form := confirmDialog.GetForm(); form != nil {
		a.tviewApp.SetFocus(form)
	}
}

// showApplyConfirmation shows confirmation dialog for terraform apply
func (a *AppNew) showApplyConfirmation() {
	path := a.treeView.GetCurrentPath()
	if path == "" {
		return
	}

	info := components.GetTerraformCommandInfo(path, a.config.Commands.ApplyTemplate, a.config.Commands.TfvarsFile, a.config)
	
	confirmDialog := dialog.NewTerraformConfirmDialog(
		"terraform apply",
		info.WorkDir,
		info.ConfigFile,
		info.Content,
		func() {
			a.pages.RemovePage("confirm_tf")
			a.executeTerraformCommand("Apply", info.WorkDir, info.Command)
		},
		func() {
			a.pages.RemovePage("confirm_tf")
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("confirm_tf", confirmDialog, true, true)
	if form := confirmDialog.GetForm(); form != nil {
		a.tviewApp.SetFocus(form)
	}
}

// showDestroyConfirmation shows confirmation dialog for terraform destroy
func (a *AppNew) showDestroyConfirmation(path string) {
	info := components.GetTerraformCommandInfo(path, a.config.Commands.DestroyTemplate, a.config.Commands.TfvarsFile, a.config)
	
	confirmDialog := dialog.NewTerraformConfirmDialog(
		"terraform destroy",
		info.WorkDir,
		info.ConfigFile,
		info.Content,
		func() {
			a.pages.RemovePage("confirm_tf")
			a.executeTerraformCommand("Destroy", info.WorkDir, info.Command)
		},
		func() {
			a.pages.RemovePage("confirm_tf")
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("confirm_tf", confirmDialog, true, true)
	if form := confirmDialog.GetForm(); form != nil {
		a.tviewApp.SetFocus(form)
	}
}

// executeTerraformCommand executes a terraform command
func (a *AppNew) executeTerraformCommand(action, workDir, cmdStr string) {
	a.contentView.Clear()
	a.contentView.SetTitle(fmt.Sprintf(" 🚀 Terraform %s ", action))
	fmt.Fprintf(a.contentView, "[yellow]Executing Terraform %s[white]\n", action)
	fmt.Fprintf(a.contentView, "[cyan]Directory:[white] %s\n", workDir)
	fmt.Fprintf(a.contentView, "[cyan]Command:[white] %s\n", cmdStr)
	fmt.Fprintf(a.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))

	go func() {
		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			return
		}
		
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Dir = workDir
		
		startTime := time.Now()
		output, err := cmd.CombinedOutput()
		
		// Save to history if it's apply or destroy
		if (action == "Apply" || action == "Destroy") && a.historyDB != nil {
			// Extract config file from command
			configFile := ""
			configData := ""
			for i, part := range parts {
				if strings.HasPrefix(part, "-var-file=") {
					configFile = strings.TrimPrefix(part, "-var-file=")
					// Read config file content
					if data, readErr := ioutil.ReadFile(configFile); readErr == nil {
						configData = string(data)
					}
					break
				} else if part == "-var-file" && i+1 < len(parts) {
					configFile = parts[i+1]
					if data, readErr := ioutil.ReadFile(configFile); readErr == nil {
						configData = string(data)
					}
					break
				}
			}
			
			entry := &db.HistoryEntry{
				Directory:  workDir,
				Action:     strings.ToLower(action),
				Timestamp:  startTime,
				ConfigFile: configFile,
				ConfigData: configData,
				Success:    err == nil,
				ErrorMsg:   "",
			}
			if err != nil {
				entry.ErrorMsg = err.Error()
			}
			
			if saveErr := a.historyDB.AddEntry(entry); saveErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to save history: %v\n", saveErr)
			}
		}
		
		a.tviewApp.QueueUpdateDraw(func() {
			if err != nil {
				fmt.Fprintf(a.contentView, "[red]Error:[white] %v\n\n", err)
			}
			a.contentView.AppendText(string(output))
			a.contentView.AppendText("\n\n[green]Done.[white]")
			
			// Show saved to history message
			if action == "Apply" || action == "Destroy" {
				a.contentView.AppendText("\n[gray](Saved to history)[white]")
			}
		})
	}()
}

// showHelp displays the help screen
func (a *AppNew) showHelp() {
	if a.helpView == nil {
		a.helpView = view.NewHelpView()
	}
	
	// Create modal with help view
	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(a.helpView, 30, 0, true).
			AddItem(nil, 0, 1, false), 100, 0, true).
		AddItem(nil, 0, 1, false)
	modal.SetBackgroundColor(tcell.ColorBlack)
	
	// Setup input capture for help view
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Rune() == 'q' || event.Rune() == '?' || event.Rune() == 'H' {
			a.pages.RemovePage("help")
			a.tviewApp.SetFocus(a.treeView)
			return nil
		}
		return event
	})
	
	a.pages.AddPage("help", modal, true, true)
}

// showHistory shows the history screen
func (a *AppNew) showHistory(path string) {
	// Get directory path
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		path = filepath.Dir(path)
	}

	// Get history from DB
	var entries []*db.HistoryEntry
	if a.historyDB != nil {
		entries, err = a.historyDB.GetByDirectory(path, 100) // Get up to 100 entries
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading history: %v\n", err)
			entries = []*db.HistoryEntry{}
		}
	}

	// Create history view
	a.historyView = view.NewHistoryView(path, entries)
	
	// Set up key handler for history view
	a.historyView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			a.pages.RemovePage("history")
			a.tviewApp.SetFocus(a.treeView)
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'M':
				// Shift+M: Toggle details
				a.historyView.ToggleDetails()
				return nil
			}
		case tcell.KeyDown:
			if event.Modifiers()&tcell.ModShift != 0 {
				// Shift+Down: Load more
				a.historyView.LoadMore()
				return nil
			}
		case tcell.KeyUp:
			if event.Modifiers()&tcell.ModShift != 0 {
				// Shift+Up: Load less (go back)
				a.historyView.LoadLess()
				return nil
			}
		}
		return event
	})

	a.pages.AddPage("history", a.historyView, true, true)
	a.tviewApp.SetFocus(a.historyView)
}

// showCommandInput displays the command input bar
func (a *AppNew) showCommandInput() {
	// Get current selected path
	path := a.treeView.GetCurrentPath()
	if path == "" {
		path = a.currentDir
	}
	
	// Check if it's a file, if so use parent directory
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		path = filepath.Dir(path)
	}
	
	// Create or update command view
	if a.commandView == nil {
		a.commandView = view.NewCommandView(path)
		a.commandView.SetExecuteHandler(func(cmd string) {
			a.executeCommand(cmd)
			a.pages.RemovePage("command")
			a.tviewApp.SetFocus(a.treeView)
		})
	} else {
		a.commandView.UpdatePath(path)
		a.commandView.Clear()
	}
	
	// Setup input capture
	a.commandView.GetInput().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("command")
			a.tviewApp.SetFocus(a.treeView)
			return nil
		}
		return event
	})
	
	a.pages.AddPage("command", a.commandView, true, true)
	a.tviewApp.SetFocus(a.commandView.GetInput())
}

// executeCommand executes a command in the current directory
func (a *AppNew) executeCommand(cmd string) {
	if cmd == "" {
		return
	}
	
	workDir := a.commandView.GetCurrentDir()
	
	a.contentView.Clear()
	a.contentView.SetTitle(" 🚀 Command Execution ")
	fmt.Fprintf(a.contentView, "[yellow]Executing Command[white]\n")
	fmt.Fprintf(a.contentView, "[cyan]Directory:[white] %s\n", workDir)
	fmt.Fprintf(a.contentView, "[cyan]Command:[white] %s\n", cmd)
	fmt.Fprintf(a.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))
	
	go func() {
		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			return
		}
		
		execCmd := exec.Command(parts[0], parts[1:]...)
		execCmd.Dir = workDir
		
		output, err := execCmd.CombinedOutput()
		
		a.tviewApp.QueueUpdateDraw(func() {
			if err != nil {
				fmt.Fprintf(a.contentView, "[red]Error:[white] %v\n\n", err)
			}
			a.contentView.AppendText(string(output))
			a.contentView.AppendText("\n\n[green]Done.[white]")
		})
	}()
}

// Run starts the application
func (a *AppNew) Run() error {
	return a.tviewApp.Run()
}

