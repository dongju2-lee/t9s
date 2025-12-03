package view

import (
	"fmt"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// StatusBar represents the status bar
type StatusBar struct {
	*tview.TextView
	currentDir      string
	focusIndicator  string
}

// NewStatusBar creates a new status bar
func NewStatusBar(currentDir string) *StatusBar {
	sb := &StatusBar{
		TextView:       tview.NewTextView().SetDynamicColors(true),
		currentDir:     currentDir,
		focusIndicator: "File Tree", // Default focus on tree
	}

	sb.SetBackgroundColor(tcell.ColorBlack)

	sb.ShowDefault()
	return sb
}

// ShowDefault shows the default status message
func (sb *StatusBar) ShowDefault() {
	sb.Clear()
	fmt.Fprintf(sb, "[green]● %s[white]  [yellow]|[white]  [yellow]Tab[white] Switch Focus  [yellow]|[white]  [yellow]↑↓[white] Navigate  [yellow]Enter[white] Expand/View  [yellow]q[white] Quit", sb.focusIndicator)
}

// UpdatePath updates the status bar with current path
func (sb *StatusBar) UpdatePath(path string) {
	relPath, _ := filepath.Rel(sb.currentDir, path)
	sb.Clear()
	fmt.Fprintf(sb, "[yellow]Current:[white] %s  [yellow]|[white]  [green]● %s[white]  [yellow]|[white]  [yellow]Tab[white] Switch Focus  [yellow]|[white]  [yellow]h[white] History  [yellow]e[white] Edit  [yellow]q[white] Quit", relPath, sb.focusIndicator)
}

// ShowMessage displays a message in the status bar
func (sb *StatusBar) ShowMessage(msg string) {
	sb.Clear()
	fmt.Fprintf(sb, "%s", msg)
}

// SetFocusIndicator updates the focus indicator
func (sb *StatusBar) SetFocusIndicator(focus string) {
	sb.focusIndicator = focus
	sb.ShowDefault()
}

