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
	currentDir string
}

// NewStatusBar creates a new status bar
func NewStatusBar(currentDir string) *StatusBar {
	sb := &StatusBar{
		TextView:   tview.NewTextView().SetDynamicColors(true),
		currentDir: currentDir,
	}

	sb.SetBackgroundColor(tcell.ColorBlack)

	sb.ShowDefault()
	return sb
}

// ShowDefault shows the default status message
func (sb *StatusBar) ShowDefault() {
	sb.Clear()
	fmt.Fprintf(sb, "[yellow]<↑↓>[white] Navigate  [yellow]<Enter>[white] Expand/View  [yellow]<q>[white] Quit")
}

// UpdatePath updates the status bar with current path
func (sb *StatusBar) UpdatePath(path string) {
	relPath, _ := filepath.Rel(sb.currentDir, path)
	sb.Clear()
	fmt.Fprintf(sb, "[cyan]═══════════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(sb, "[yellow]Current:[white] %s  [yellow]|[white]  [yellow]h[white] History  [yellow]H[white] Helm  [yellow]e[white] Edit  [yellow]q[white] Quit\n", relPath)
	fmt.Fprintf(sb, "[cyan]═══════════════════════════════════════════════════════════════════")
}

// ShowMessage displays a message in the status bar
func (sb *StatusBar) ShowMessage(msg string) {
	sb.Clear()
	fmt.Fprintf(sb, "%s", msg)
}

