package ui

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/idongju/t9s/internal/config"
	"github.com/idongju/t9s/internal/ui/components"
	"github.com/idongju/t9s/internal/ui/dialog"
	"github.com/idongju/t9s/internal/view"
)

// AppNew represents the main application with new structure
type AppNew struct {
	tviewApp    *tview.Application
	pages       *tview.Pages
	
	// Views
	headerView  *view.HeaderView
	treeView    *view.TreeView
	contentView *view.ContentView
	statusBar   *view.StatusBar
	
	// Components
	executor    *components.CommandExecutor
	
	// State
	currentDir  string
	currentFile string
	config      *config.Config
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
				PlanTemplate:  "terraform plan -var-file={varfile}",
				ApplyTemplate: "terraform apply -var-file={varfile}",
				VarFile:       "config/prod.tfvars",
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

	app := &AppNew{
		tviewApp:   tview.NewApplication(),
		currentDir: currentDir,
		config:     cfg,
		pages:      tview.NewPages(),
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
	a.executor = components.NewCommandExecutor(a.tviewApp, a.contentView, a.config)
	
	// Setup tree view handler
	a.treeView.SetFileSelectHandler(func(path string) {
		a.currentFile = path
		a.contentView.DisplayFile(path)
	})
	
	// Setup tree view change handler
	a.treeView.SetChangedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference != nil {
			a.statusBar.UpdatePath(reference.(string))
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
					a.executor.ShowHistory(path)
				}
				return nil
			case 'H':
				a.executor.ShowHelm()
				return nil
			case 'e', 'E':
				if a.currentFile != "" {
					a.executor.EditFile(a.currentFile)
					a.contentView.DisplayFile(a.currentFile)
				}
				return nil
			case 'p':
				path := a.treeView.GetCurrentPath()
				if path != "" {
					a.executor.ExecutePlan(path)
				}
				return nil
			case 'a':
				a.showApplyConfirmation()
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

// showApplyConfirmation shows confirmation dialog for terraform apply
func (a *AppNew) showApplyConfirmation() {
	path := a.treeView.GetCurrentPath()
	if path == "" {
		return
	}

	confirmDialog := dialog.NewConfirmDialog(
		"Are you sure you want to execute terraform apply?",
		func() {
			a.pages.RemovePage("confirm_apply")
			a.executor.ExecuteApply(path)
		},
		func() {
			a.pages.RemovePage("confirm_apply")
			a.tviewApp.SetFocus(a.treeView)
		},
	)

	a.pages.AddPage("confirm_apply", confirmDialog, true, true)
}

// Run starts the application
func (a *AppNew) Run() error {
	return a.tviewApp.Run()
}

