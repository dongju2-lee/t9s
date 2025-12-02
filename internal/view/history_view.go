package view

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/idongju/t9s/internal/db"
)

// HistoryView displays terraform execution history
type HistoryView struct {
	*tview.TextView
	entries     []*db.HistoryEntry
	displayFrom int
	displaySize int
	showDetails bool
	directory   string
}

// NewHistoryView creates a new history view
func NewHistoryView(directory string, entries []*db.HistoryEntry) *HistoryView {
	hv := &HistoryView{
		TextView:    tview.NewTextView().SetDynamicColors(true),
		entries:     entries,
		displayFrom: 0,
		displaySize: 3,
		showDetails: false,
		directory:   directory,
	}

	hv.SetBorder(true)
	hv.SetTitle(" ⏰ Terraform History ")
	hv.SetBackgroundColor(tcell.ColorBlack)
	hv.SetTextColor(tcell.ColorWhite)

	hv.render()
	return hv
}

// render displays the history
func (hv *HistoryView) render() {
	hv.Clear()
	
	fmt.Fprintf(hv.TextView, "[yellow]Terraform History[white]\n")
	fmt.Fprintf(hv.TextView, "[cyan]Directory:[white] %s\n", hv.directory)
	fmt.Fprintf(hv.TextView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))
	
	// Show keyboard shortcuts at the top
	fmt.Fprintf(hv.TextView, "[yellow]Shortcuts:[white] ")
	if !hv.showDetails {
		fmt.Fprintf(hv.TextView, "[green]<Shift+M>[white] Show Details  ")
	} else {
		fmt.Fprintf(hv.TextView, "[green]<Shift+M>[white] Hide Details  ")
	}
	fmt.Fprintf(hv.TextView, "[green]<Shift+↓>[white] Load More  ")
	fmt.Fprintf(hv.TextView, "[green]<Shift+↑>[white] Load Less  ")
	fmt.Fprintf(hv.TextView, "[green]<Esc>[white] Back\n")
	fmt.Fprintf(hv.TextView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))

	if len(hv.entries) == 0 {
		fmt.Fprintf(hv.TextView, "[gray]No execution history found for this directory.[white]\n")
		return
	}

	// Calculate range to display
	displayEnd := hv.displayFrom + hv.displaySize
	if displayEnd > len(hv.entries) {
		displayEnd = len(hv.entries)
	}

	fmt.Fprintf(hv.TextView, "[yellow]📜 Execution History (Showing %d-%d of %d):[white]\n\n", 
		hv.displayFrom+1, displayEnd, len(hv.entries))

	// Display entries
	for i := hv.displayFrom; i < displayEnd; i++ {
		entry := hv.entries[i]
		hv.renderEntry(entry, i+1)
	}

	// Show additional info at the bottom
	fmt.Fprintf(hv.TextView, "\n[cyan]%s[white]\n", strings.Repeat("─", 60))
	
	if displayEnd < len(hv.entries) {
		fmt.Fprintf(hv.TextView, "[gray]%d more entries available (Press Shift+↓)[white]\n", 
			len(hv.entries)-displayEnd)
	} else if len(hv.entries) > hv.displaySize {
		fmt.Fprintf(hv.TextView, "[gray]End of history[white]\n")
	}
}

// renderEntry renders a single history entry
func (hv *HistoryView) renderEntry(entry *db.HistoryEntry, index int) {
	statusIcon := "✅"
	statusColor := "green"
	if !entry.Success {
		statusIcon = "❌"
		statusColor = "red"
	}

	fmt.Fprintf(hv.TextView, "[gray]#%d[white] [%s]%s %s[white] - %s\n",
		index,
		statusColor, statusIcon, strings.ToUpper(entry.Action),
		entry.Timestamp.Format("2006-01-02 15:04:05"))

	if entry.ConfigFile != "" {
		fmt.Fprintf(hv.TextView, "     [gray]Config:[white] %s\n", entry.ConfigFile)
	}

	if !entry.Success && entry.ErrorMsg != "" {
		fmt.Fprintf(hv.TextView, "     [red]Error:[white] %s\n", entry.ErrorMsg)
	}

	// Show details if enabled
	if hv.showDetails && entry.ConfigData != "" {
		fmt.Fprintf(hv.TextView, "     [cyan]Config Content:[white]\n")
		lines := strings.Split(entry.ConfigData, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				fmt.Fprintf(hv.TextView, "       [green]%s[white]\n", line)
			}
		}
	}

	fmt.Fprintf(hv.TextView, "\n")
}

// ToggleDetails toggles showing config details
func (hv *HistoryView) ToggleDetails() {
	hv.showDetails = !hv.showDetails
	hv.render()
}

// LoadMore loads more entries
func (hv *HistoryView) LoadMore() bool {
	if hv.displayFrom+hv.displaySize >= len(hv.entries) {
		return false // No more to load
	}
	hv.displayFrom += hv.displaySize
	hv.render()
	return true
}

// LoadLess goes back to previous entries
func (hv *HistoryView) LoadLess() bool {
	if hv.displayFrom <= 0 {
		return false // Already at the start
	}
	hv.displayFrom -= hv.displaySize
	if hv.displayFrom < 0 {
		hv.displayFrom = 0
	}
	hv.render()
	return true
}

// GetShowDetails returns whether details are shown
func (hv *HistoryView) GetShowDetails() bool {
	return hv.showDetails
}

