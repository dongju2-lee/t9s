package ui

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/idongju/t9s/internal/config"
	"github.com/idongju/t9s/internal/db"
	"github.com/idongju/t9s/internal/git"
	"github.com/idongju/t9s/internal/ui/components"
	"github.com/idongju/t9s/internal/ui/dialog"
	"github.com/idongju/t9s/internal/view"
	"github.com/rivo/tview"
)

// AppNew represents the main application with new structure
type AppNew struct {
	tviewApp *tview.Application
	pages    *tview.Pages

	// Views
	headerView  *view.HeaderView
	treeView    *view.TreeView
	contentView *view.ContentView
	statusBar   *view.StatusBar
	helpView    *view.HelpView
	commandView *view.CommandView
	historyView *view.HistoryView

	// Components
	executor   *components.CommandExecutor
	historyDB  *db.HistoryDB
	gitManager *git.Manager

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

	// Initialize history DB in terraform root directory (shared by all users)
	historyDB, err := db.NewHistoryDB(currentDir)
	if err != nil {
		// Fallback if DB fails - app still works
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize history DB: %v\n", err)
	}

	// Initialize git manager
	gitManager := git.NewManager()

	app := &AppNew{
		tviewApp:    tview.NewApplication(),
		currentDir:  currentDir,
		config:      cfg,
		pages:       tview.NewPages(),
		gitManager:  gitManager,
		historyDB:   historyDB,
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

	// Update git branch info in header
	if status, err := a.gitManager.GetStatus(a.currentDir); err == nil {
		a.headerView.SetGitBranch(status.Branch, status.IsDirty)
	}

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
		AddItem(a.headerView, 7, 0, false).
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
			case 'B':
				// Shift+B: Branch switch
				path := a.treeView.GetCurrentPath()
				if path != "" {
					a.showBranchSwitch(path)
				}
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
				// If focus is on content view, let it handle 'd' (scroll down)
				if !a.focusOnTree {
					return event
				}
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

				// Update git branch info in header
				if status, err := a.gitManager.GetStatus(a.currentDir); err == nil {
					a.headerView.SetGitBranch(status.Branch, status.IsDirty)
				}

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
		AddItem(a.headerView, 7, 0, false).
		AddItem(mainFlex, 0, 1, true).
		AddItem(a.statusBar, 1, 0, false)
	mainPage.SetBackgroundColor(tcell.ColorBlack)

	a.pages.RemovePage("main")
	a.pages.AddPage("main", mainPage, true, true)
}

// showInitConfirmation shows file selection for init config
func (a *AppNew) showInitConfirmation(path string) {
	// Get directory path
	info, err := os.Stat(path)
	workDir := path
	if err == nil && !info.IsDir() {
		workDir = filepath.Dir(path)
	}

	// Check if config directory exists
	configDir := filepath.Join(workDir, "config")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// No config directory, use default
		a.showInitConfirmationWithFile(path, "")
		return
	}

	// Show file selection dialog
	fileDialog := dialog.NewFileSelectionDialog(
		configDir,
		"*.conf",
		"Select Init Config File",
		func(filePath, content string) {
			// File selected, show confirmation
			a.pages.RemovePage("file_selection")
			a.showInitConfirmationWithFile(path, filePath)
		},
		func() {
			// Cancelled
			a.pages.RemovePage("file_selection")
			a.focusOnTree = true
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("file_selection", fileDialog, true, true)
	a.tviewApp.SetFocus(fileDialog.GetList())
}

// showInitConfirmationWithFile shows confirmation dialog with selected file
func (a *AppNew) showInitConfirmationWithFile(path, configFile string) {
	info := components.GetTerraformCommandInfo(path, a.config.Commands.InitTemplate, configFile, a.config)

	confirmDialog := dialog.NewTerraformConfirmDialog(
		"terraform init",
		info.WorkDir,
		info.ConfigFile,
		info.Content,
		// Execute: normal execution (no auto-approve)
		func() {
			a.pages.RemovePage("confirm_tf")
			a.executeTerraformCommand("Init", info.WorkDir, info.Command, info.ConfigFile, info.Content)
		},
		// Auto Approve: add -auto-approve flag
		func() {
			a.pages.RemovePage("confirm_tf")
			cmdWithAutoApprove := info.Command + " -auto-approve"
			a.executeTerraformCommand("Init", info.WorkDir, cmdWithAutoApprove, info.ConfigFile, info.Content)
		},
		// Cancel
		func() {
			a.pages.RemovePage("confirm_tf")
			a.focusOnTree = true
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("confirm_tf", confirmDialog, true, true)
	if form := confirmDialog.GetForm(); form != nil {
		a.tviewApp.SetFocus(form)
	}
}

// showPlanConfirmation shows file selection for plan tfvars
func (a *AppNew) showPlanConfirmation(path string) {
	// Get directory path
	info, err := os.Stat(path)
	workDir := path
	if err == nil && !info.IsDir() {
		workDir = filepath.Dir(path)
	}

	// Check if config directory exists
	configDir := filepath.Join(workDir, "config")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// No config directory, use default
		a.showPlanConfirmationWithFile(path, "")
		return
	}

	// Show file selection dialog
	fileDialog := dialog.NewFileSelectionDialog(
		configDir,
		"*.tfvars",
		"Select Terraform Variables File for Plan",
		func(filePath, content string) {
			// File selected, show confirmation
			a.pages.RemovePage("file_selection")
			a.showPlanConfirmationWithFile(path, filePath)
		},
		func() {
			// Cancelled
			a.pages.RemovePage("file_selection")
			a.focusOnTree = true
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("file_selection", fileDialog, true, true)
	a.tviewApp.SetFocus(fileDialog.GetList())
}

// showPlanConfirmationWithFile shows confirmation dialog with selected file
func (a *AppNew) showPlanConfirmationWithFile(path, configFile string) {
	info := components.GetTerraformCommandInfo(path, a.config.Commands.PlanTemplate, configFile, a.config)

	confirmDialog := dialog.NewTerraformConfirmDialog(
		"terraform plan",
		info.WorkDir,
		info.ConfigFile,
		info.Content,
		// Execute: normal execution (no auto-approve, plan doesn't need it anyway)
		func() {
			a.pages.RemovePage("confirm_tf")
			a.executeTerraformCommand("Plan", info.WorkDir, info.Command, info.ConfigFile, info.Content)
		},
		// Auto Approve: not applicable for plan, but still provide the option
		func() {
			a.pages.RemovePage("confirm_tf")
			a.executeTerraformCommand("Plan", info.WorkDir, info.Command, info.ConfigFile, info.Content)
		},
		// Cancel
		func() {
			a.pages.RemovePage("confirm_tf")
			a.focusOnTree = true
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("confirm_tf", confirmDialog, true, true)
	if form := confirmDialog.GetForm(); form != nil {
		a.tviewApp.SetFocus(form)
	}
}

// showApplyConfirmation shows file selection for apply tfvars
func (a *AppNew) showApplyConfirmation() {
	path := a.treeView.GetCurrentPath()
	if path == "" {
		return
	}

	// Get directory path
	info, err := os.Stat(path)
	workDir := path
	if err == nil && !info.IsDir() {
		workDir = filepath.Dir(path)
	}

	// Check if config directory exists
	configDir := filepath.Join(workDir, "config")

	// DEBUG LOGGING
	f, _ := os.OpenFile("/tmp/t9s_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if f != nil {
		fmt.Fprintf(f, "Apply Path: %s\nWorkDir: %s\nConfigDir: %s\n", path, workDir, configDir)
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			fmt.Fprintf(f, "ConfigDir does NOT exist: %v\n", err)
		} else {
			fmt.Fprintf(f, "ConfigDir exists\n")
		}
		f.Close()
	}

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// No config directory, use default
		a.showApplyConfirmationWithFile(path, "")
		return
	}

	// Show file selection dialog
	fileDialog := dialog.NewFileSelectionDialog(
		configDir,
		"*.tfvars",
		"Select Terraform Variables File for Apply",
		func(filePath, content string) {
			// File selected, show confirmation
			a.pages.RemovePage("file_selection")
			a.showApplyConfirmationWithFile(path, filePath)
		},
		func() {
			// Cancelled
			a.pages.RemovePage("file_selection")
			a.focusOnTree = true
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("file_selection", fileDialog, true, true)
	a.tviewApp.SetFocus(fileDialog.GetList())
}

// showApplyConfirmationWithFile shows confirmation dialog with selected file
func (a *AppNew) showApplyConfirmationWithFile(path, configFile string) {
	info := components.GetTerraformCommandInfo(path, a.config.Commands.ApplyTemplate, configFile, a.config)

	confirmDialog := dialog.NewTerraformConfirmDialog(
		"terraform apply",
		info.WorkDir,
		info.ConfigFile,
		info.Content,
		// Execute: normal execution (terraform will ask for 'yes')
		func() {
			a.pages.RemovePage("confirm_tf")
			a.executeTerraformCommand("Apply", info.WorkDir, info.Command, info.ConfigFile, info.Content)
		},
		// Auto Approve: add -auto-approve flag
		func() {
			a.pages.RemovePage("confirm_tf")
			cmdWithAutoApprove := info.Command + " -auto-approve"
			a.executeTerraformCommand("Apply", info.WorkDir, cmdWithAutoApprove, info.ConfigFile, info.Content)
		},
		// Cancel
		func() {
			a.pages.RemovePage("confirm_tf")
			a.focusOnTree = true
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("confirm_tf", confirmDialog, true, true)
	if form := confirmDialog.GetForm(); form != nil {
		a.tviewApp.SetFocus(form)
	}
}

// showDestroyConfirmation shows file selection for destroy tfvars
func (a *AppNew) showDestroyConfirmation(path string) {
	// Get directory path
	info, err := os.Stat(path)
	workDir := path
	if err == nil && !info.IsDir() {
		workDir = filepath.Dir(path)
	}

	// Check if config directory exists
	configDir := filepath.Join(workDir, "config")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// No config directory, use default
		a.showDestroyConfirmationWithFile(path, "")
		return
	}

	// Show file selection dialog
	fileDialog := dialog.NewFileSelectionDialog(
		configDir,
		"*.tfvars",
		"Select Terraform Variables File for Destroy",
		func(filePath, content string) {
			// File selected, show confirmation
			a.pages.RemovePage("file_selection")
			a.showDestroyConfirmationWithFile(path, filePath)
		},
		func() {
			// Cancelled
			a.pages.RemovePage("file_selection")
			a.focusOnTree = true
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("file_selection", fileDialog, true, true)
	a.tviewApp.SetFocus(fileDialog.GetList())
}

// showDestroyConfirmationWithFile shows confirmation dialog with selected file
func (a *AppNew) showDestroyConfirmationWithFile(path, configFile string) {
	info := components.GetTerraformCommandInfo(path, a.config.Commands.DestroyTemplate, configFile, a.config)

	confirmDialog := dialog.NewTerraformConfirmDialog(
		"terraform destroy",
		info.WorkDir,
		info.ConfigFile,
		info.Content,
		// Execute: normal execution (terraform will ask for 'yes')
		func() {
			a.pages.RemovePage("confirm_tf")
			a.executeTerraformCommand("Destroy", info.WorkDir, info.Command, info.ConfigFile, info.Content)
		},
		// Auto Approve: add -auto-approve flag
		func() {
			a.pages.RemovePage("confirm_tf")
			cmdWithAutoApprove := info.Command + " -auto-approve"
			a.executeTerraformCommand("Destroy", info.WorkDir, cmdWithAutoApprove, info.ConfigFile, info.Content)
		},
		// Cancel
		func() {
			a.pages.RemovePage("confirm_tf")
			a.focusOnTree = true
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("confirm_tf", confirmDialog, true, true)
	if form := confirmDialog.GetForm(); form != nil {
		a.tviewApp.SetFocus(form)
	}
}

// executeTerraformCommand executes a terraform command with real-time streaming output
func (a *AppNew) executeTerraformCommand(action, workDir, cmdStr, configFile, configData string) {
	// Auto-switch focus to content view for Apply/Destroy to prevent accidental input
	if action == "Apply" || action == "Destroy" {
		a.focusOnTree = false
		a.tviewApp.SetFocus(a.contentView)
	}

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

		// Check if command has -auto-approve flag
		hasAutoApprove := strings.Contains(cmdStr, "-auto-approve")

		// Create stdin pipe for interactive input (only if no auto-approve)
		var stdinPipe io.WriteCloser
		var stdinErr error
		if !hasAutoApprove && (action == "Apply" || action == "Destroy") {
			stdinPipe, stdinErr = cmd.StdinPipe()
			if stdinErr != nil {
				a.tviewApp.QueueUpdateDraw(func() {
					fmt.Fprintf(a.contentView, "[red]Error creating stdin pipe:[white] %v\n", stdinErr)
				})
				return
			}
		}

		// Create pipes for real-time output streaming
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.contentView, "[red]Error creating stdout pipe:[white] %v\n", err)
			})
			return
		}

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.contentView, "[red]Error creating stderr pipe:[white] %v\n", err)
			})
			return
		}

		// Start the command
		startTime := time.Now()
		if err := cmd.Start(); err != nil {
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.contentView, "[red]Error starting command:[white] %v\n", err)
			})
			return
		}

		// Buffer to save full output for history
		var outputBuf bytes.Buffer

		// Channel to signal when "Enter a value:" is detected
		confirmChan := make(chan bool)
		userResponseChan := make(chan string)

		// Stream stdout in real-time and detect "Enter a value:"
		go func() {
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				line := scanner.Text()
				outputBuf.WriteString(line + "\n")

				// Detect terraform asking for confirmation
				if !hasAutoApprove && (strings.Contains(line, "Enter a value:") || strings.Contains(line, "Only 'yes' will be accepted")) {
					// Show confirmation dialog
					a.tviewApp.QueueUpdateDraw(func() {
						w := tview.ANSIWriter(a.contentView)
						w.Write([]byte(line + "\n"))
						a.contentView.ScrollToEnd()

						// Show Yes/No dialog
						a.showApplyConfirmDialog(func() {
							// User selected Yes
							userResponseChan <- "yes\n"
						}, func() {
							// User selected No
							userResponseChan <- "no\n"
						})
					})

					// Wait for user response
					response := <-userResponseChan
					if stdinPipe != nil {
						stdinPipe.Write([]byte(response))
						stdinPipe.Close()
					}
					confirmChan <- true
				} else {
					a.tviewApp.QueueUpdateDraw(func() {
						w := tview.ANSIWriter(a.contentView)
						w.Write([]byte(line + "\n"))
						a.contentView.ScrollToEnd()
					})
				}
			}
		}()

		// Stream stderr in real-time
		go func() {
			scanner := bufio.NewScanner(stderrPipe)
			for scanner.Scan() {
				line := scanner.Text()
				outputBuf.WriteString(line + "\n")
				a.tviewApp.QueueUpdateDraw(func() {
					w := tview.ANSIWriter(a.contentView)
					w.Write([]byte(line + "\n"))
					a.contentView.ScrollToEnd()
				})
			}
		}()

		// Wait for command to complete
		cmdErr := cmd.Wait()

		// Save to history if it's apply or destroy
		if (action == "Apply" || action == "Destroy") && a.historyDB != nil {
			// Get user and branch info
			user := os.Getenv("USER")
			if user == "" {
				user = "unknown"
			}

			branch := ""
			if status, gitErr := a.gitManager.GetStatus(workDir); gitErr == nil {
				branch = status.Branch
			}

			entry := &db.HistoryEntry{
				Directory:  workDir,
				Action:     strings.ToLower(action),
				Timestamp:  startTime,
				User:       user,
				Branch:     branch,
				ConfigFile: configFile,
				ConfigData: configData,
				Success:    cmdErr == nil,
				ErrorMsg:   "",
			}
			if cmdErr != nil {
				entry.ErrorMsg = cmdErr.Error()
			}

			if saveErr := a.historyDB.AddEntry(entry); saveErr != nil {
				fmt.Fprintf(os.Stderr, "Failed to save history: %v\n", saveErr)
			}
		}

		a.tviewApp.QueueUpdateDraw(func() {
			if cmdErr != nil {
				fmt.Fprintf(a.contentView, "\n[red]Error:[white] %v\n", cmdErr)
			}

			fmt.Fprintf(a.contentView, "\n[green]Done.[white]")

			// Show saved to history message
			if action == "Apply" || action == "Destroy" {
				fmt.Fprintf(a.contentView, "\n[gray](Saved to history)[white]")
			}

			// Scroll to end
			a.contentView.ScrollToEnd()
		})
	}()
}

// showApplyConfirmDialog shows a Yes/No dialog for terraform confirmation
func (a *AppNew) showApplyConfirmDialog(onYes, onNo func()) {
	confirmDialog := tview.NewModal().
		SetText("Do you want to perform these actions?\nTerraform will perform the actions described above.").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage("apply_confirm")
			a.tviewApp.SetFocus(a.contentView)

			if buttonLabel == "Yes" {
				onYes()
			} else {
				onNo()
			}
		})

	confirmDialog.SetBackgroundColor(tcell.ColorBlack)
	confirmDialog.SetBorderColor(tcell.NewRGBColor(255, 165, 0))
	confirmDialog.SetButtonBackgroundColor(tcell.NewRGBColor(50, 50, 50))
	confirmDialog.SetButtonTextColor(tcell.ColorWhite)

	a.pages.AddPage("apply_confirm", confirmDialog, true, true)
	a.tviewApp.SetFocus(confirmDialog)
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
			case 'd', 'D':
				// d: Load more
				a.historyView.LoadMore()
				return nil
			case 'u', 'U':
				// u: Load less (go back)
				a.historyView.LoadLess()
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

// showBranchSwitch shows branch selection dialog
func (a *AppNew) showBranchSwitch(path string) {
	// Get directory path
	info, err := os.Stat(path)
	workDir := path
	if err == nil && !info.IsDir() {
		workDir = filepath.Dir(path)
	}

	// Get branches
	branches, currentBranch, err := a.gitManager.GetBranches(workDir)
	if err != nil {
		a.contentView.DisplayText("Error", fmt.Sprintf("[red]Failed to get branches: %v[white]", err))
		return
	}

	if len(branches) == 0 {
		a.contentView.DisplayText("Info", "[yellow]No git branches found in this directory[white]")
		return
	}

	// Show branch selection dialog
	branchDialog := dialog.NewBranchDialog(
		branches,
		currentBranch,
		func(selectedBranch string) {
			a.pages.RemovePage("branch_select")
			a.switchBranch(workDir, currentBranch, selectedBranch)
		},
		func() {
			a.pages.RemovePage("branch_select")
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("branch_select", branchDialog, true, true)
	a.tviewApp.SetFocus(branchDialog.GetList())
}

// switchBranch handles branch switching with dirty working tree check
func (a *AppNew) switchBranch(workDir, currentBranch, targetBranch string) {
	// Don't switch if same branch
	if currentBranch == targetBranch {
		a.tviewApp.SetFocus(a.treeView)
		return
	}

	// Check for local changes
	status, err := a.gitManager.GetStatus(workDir)
	if err != nil {
		a.contentView.DisplayText("Error", fmt.Sprintf("[red]Failed to check git status: %v[white]", err))
		a.tviewApp.SetFocus(a.treeView)
		return
	}

	// If working tree is clean, switch directly
	if !status.IsDirty {
		a.performBranchSwitch(workDir, targetBranch, "")
		return
	}

	// Show dirty working tree dialog
	dirtyDialog := dialog.NewDirtyBranchDialog(
		currentBranch,
		targetBranch,
		status.ModifiedFiles,
		func() {
			// Stash & Switch
			a.pages.RemovePage("dirty_branch")
			a.performBranchSwitch(workDir, targetBranch, "stash")
		},
		func() {
			// Commit & Switch
			a.pages.RemovePage("dirty_branch")
			a.showCommitDialog(workDir, targetBranch)
		},
		func() {
			// Force Switch (Discard)
			a.pages.RemovePage("dirty_branch")
			a.performBranchSwitch(workDir, targetBranch, "force")
		},
		func() {
			// Cancel
			a.pages.RemovePage("dirty_branch")
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("dirty_branch", dirtyDialog, true, true)
	a.tviewApp.SetFocus(dirtyDialog.GetForm())
}

// showCommitDialog shows dialog to input commit message
func (a *AppNew) showCommitDialog(workDir, targetBranch string) {
	commitDialog := dialog.NewCommitDialog(
		func(message string) {
			a.pages.RemovePage("commit_dialog")
			a.performBranchSwitch(workDir, targetBranch, "commit:"+message)
		},
		func() {
			a.pages.RemovePage("commit_dialog")
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("commit_dialog", commitDialog, true, true)
	a.tviewApp.SetFocus(commitDialog.GetForm())
}

// performBranchSwitch performs the actual branch switch
func (a *AppNew) performBranchSwitch(workDir, targetBranch, mode string) {
	a.contentView.Clear()
	a.contentView.SetTitle(" 🌿 Branch Switch ")

	fmt.Fprintf(a.contentView, "[yellow]Switching Branch[white]\n")
	fmt.Fprintf(a.contentView, "[cyan]Directory:[white] %s\n", workDir)
	fmt.Fprintf(a.contentView, "[cyan]Target Branch:[white] %s\n", targetBranch)
	fmt.Fprintf(a.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))

	go func() {
		var err error

		// Handle different modes
		if mode == "stash" {
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.contentView, "[yellow]Stashing changes...[white]\n")
			})
			err = a.gitManager.StashChanges(workDir, fmt.Sprintf("T9s: Auto-stash before switching to %s", targetBranch))
			if err != nil {
				a.tviewApp.QueueUpdateDraw(func() {
					fmt.Fprintf(a.contentView, "[red]Failed to stash: %v[white]\n", err)
				})
				return
			}
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.contentView, "[green]Changes stashed successfully[white]\n\n")
			})
		} else if strings.HasPrefix(mode, "commit:") {
			message := strings.TrimPrefix(mode, "commit:")
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.contentView, "[yellow]Committing changes...[white]\n")
			})
			err = a.gitManager.CommitAll(workDir, message)
			if err != nil {
				a.tviewApp.QueueUpdateDraw(func() {
					fmt.Fprintf(a.contentView, "[red]Failed to commit: %v[white]\n", err)
				})
				return
			}
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.contentView, "[green]Changes committed successfully[white]\n\n")
			})
		} else if mode == "force" {
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.contentView, "[yellow]Force switching (discarding changes)...[white]\n")
			})
			err = a.gitManager.CheckoutBranchForce(workDir, targetBranch)
		} else {
			// Clean switch
			a.tviewApp.QueueUpdateDraw(func() {
				fmt.Fprintf(a.contentView, "[yellow]Switching branch...[white]\n")
			})
			err = a.gitManager.CheckoutBranch(workDir, targetBranch)
		}

		// Perform checkout for non-force modes
		if mode != "force" && err == nil {
			err = a.gitManager.CheckoutBranch(workDir, targetBranch)
		}

		a.tviewApp.QueueUpdateDraw(func() {
			if err != nil {
				fmt.Fprintf(a.contentView, "[red]Failed to switch branch: %v[white]\n", err)
			} else {
				fmt.Fprintf(a.contentView, "[green]Successfully switched to branch: %s[white]\n", targetBranch)

				// Update header with new branch info
				if status, err := a.gitManager.GetStatus(workDir); err == nil {
					a.headerView.SetGitBranch(status.Branch, status.IsDirty)
				}
			}
			a.tviewApp.SetFocus(a.treeView)
		})
	}()
}

// Run starts the application
func (a *AppNew) Run() error {
	return a.tviewApp.Run()
}
