package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App represents the main application
type App struct {
	tviewApp *tview.Application
	pages    *tview.Pages
	layout   *tview.Flex
}

// NewApp creates a new P9s application
func NewApp() *App {
	app := &App{
		tviewApp: tview.NewApplication(),
		pages:    tview.NewPages(),
	}

	app.setupUI()
	app.setupKeyBindings()

	return app
}

// setupUI initializes the UI components
func (a *App) setupUI() {
	// Create header
	header := a.createHeader()

	// Create main content area
	mainContent := a.createMainContent()

	// Create footer with shortcuts
	footer := a.createFooter()

	// Layout
	a.layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(mainContent, 0, 1, true).
		AddItem(footer, 3, 0, false)

	a.tviewApp.SetRoot(a.layout, true)
}

// createHeader creates the application header
func (a *App) createHeader() *tview.TextView {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	fmt.Fprintf(header, "[::b][cyan]╔═══════════════════════════════════════════════════════════════════╗\n")
	fmt.Fprintf(header, "[::b][cyan]║[white]  P9s - Terraform Infrastructure Manager  [cyan]║\n")
	fmt.Fprintf(header, "[::b][cyan]╚═══════════════════════════════════════════════════════════════════╝")

	return header
}

// createMainContent creates the main content area with directory list and details
func (a *App) createMainContent() *tview.Flex {
	// Directory list (left panel)
	dirList := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	dirList.SetTitle(" 📁 Terraform Directories ").SetBorder(true)

	// Color definitions
	colorYellow := tcell.NewHexColor(0xFFD700)
	colorWhite := tcell.NewHexColor(0xFFFFFF)
	colorGreen := tcell.NewHexColor(0x00FF00)
	colorCyan := tcell.NewHexColor(0x00FFFF)
	colorGray := tcell.NewHexColor(0x808080)

	// Add sample data
	dirList.SetCell(0, 0, tview.NewTableCell("Directory").
		SetTextColor(colorYellow).
		SetAlign(tview.AlignCenter).
		SetSelectable(false))
	dirList.SetCell(0, 1, tview.NewTableCell("Status").
		SetTextColor(colorYellow).
		SetAlign(tview.AlignCenter).
		SetSelectable(false))
	dirList.SetCell(0, 2, tview.NewTableCell("Last Apply").
		SetTextColor(colorYellow).
		SetAlign(tview.AlignCenter).
		SetSelectable(false))

	// Sample directories
	directories := []struct {
		name      string
		status    string
		lastApply string
	}{
		{"s3-buckets", "✓ Synced", "2h ago"},
		{"eks-cluster", "⚠ Drift", "1d ago"},
		{"vpc-network", "✓ Synced", "3h ago"},
		{"rds-databases", "⚠ Drift", "2d ago"},
		{"lambda-functions", "✓ Synced", "5h ago"},
		{"helm-charts", "? Unknown", "N/A"},
	}

	for i, dir := range directories {
		row := i + 1
		dirList.SetCell(row, 0, tview.NewTableCell(dir.name).
			SetTextColor(colorWhite))

		statusColor := colorGreen
		if dir.status == "⚠ Drift" {
			statusColor = colorYellow
		} else if dir.status == "? Unknown" {
			statusColor = colorGray
		}

		dirList.SetCell(row, 1, tview.NewTableCell(dir.status).
			SetTextColor(statusColor).
			SetAlign(tview.AlignCenter))
		dirList.SetCell(row, 2, tview.NewTableCell(dir.lastApply).
			SetTextColor(colorCyan).
			SetAlign(tview.AlignCenter))
	}

	// Details panel (right panel)
	detailsView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)

	detailsView.SetTitle(" 📋 Details ").SetBorder(true)
	fmt.Fprintf(detailsView, "[yellow]Select a directory to view details[white]\n\n")
	fmt.Fprintf(detailsView, "[cyan]Available Actions:[white]\n")
	fmt.Fprintf(detailsView, "  • [green]p[white] - Plan\n")
	fmt.Fprintf(detailsView, "  • [green]a[white] - Apply\n")
	fmt.Fprintf(detailsView, "  • [green]e[white] - Edit tfvars\n")
	fmt.Fprintf(detailsView, "  • [green]g[white] - Git diff\n")
	fmt.Fprintf(detailsView, "  • [green]s[white] - State info\n")
	fmt.Fprintf(detailsView, "  • [green]h[white] - Helm charts\n")

	// Update details when selection changes
	dirList.SetSelectedFunc(func(row, col int) {
		if row > 0 {
			dir := directories[row-1]
			detailsView.Clear()
			fmt.Fprintf(detailsView, "[yellow]Directory:[white] %s\n\n", dir.name)
			fmt.Fprintf(detailsView, "[cyan]Status:[white] %s\n", dir.status)
			fmt.Fprintf(detailsView, "[cyan]Last Apply:[white] %s\n\n", dir.lastApply)
			fmt.Fprintf(detailsView, "[cyan]Configuration:[white]\n")
			fmt.Fprintf(detailsView, "  • Config path: ./config/\n")
			fmt.Fprintf(detailsView, "  • tfvars: %s.tfvars\n", dir.name)
			fmt.Fprintf(detailsView, "  • Backend: s3://terraform-state/%s\n\n", dir.name)
			fmt.Fprintf(detailsView, "[cyan]Recent Changes:[white]\n")
			fmt.Fprintf(detailsView, "  • No recent changes detected\n")
		}
	})

	// Main content layout
	mainFlex := tview.NewFlex().
		AddItem(dirList, 0, 1, true).
		AddItem(detailsView, 0, 1, false)

	return mainFlex
}

// createFooter creates the footer with keyboard shortcuts
func (a *App) createFooter() *tview.TextView {
	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	fmt.Fprintf(footer, "[cyan]═══════════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(footer, "[yellow]↑↓[white] Navigate  [yellow]Enter[white] Select  [yellow]p[white] Plan  [yellow]a[white] Apply  [yellow]e[white] Edit  [yellow]g[white] Git  [yellow]q[white] Quit\n")
	fmt.Fprintf(footer, "[cyan]═══════════════════════════════════════════════════════════════════")

	return footer
}

// setupKeyBindings sets up global key bindings
func (a *App) setupKeyBindings() {
	a.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			a.tviewApp.Stop()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				a.tviewApp.Stop()
				return nil
			}
		}
		return event
	})
}

// Run starts the application
func (a *App) Run() error {
	return a.tviewApp.Run()
}
