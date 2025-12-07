package view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HelpView represents the help screen
type HelpView struct {
	*tview.Flex
}

// NewHelpView creates a new help view
func NewHelpView() *HelpView {
	hv := &HelpView{
		Flex: tview.NewFlex(),
	}

	// Create help sections for T9s (실제 구현된 단축키)
	resourceSection := hv.createSection("RESOURCE", []HelpItem{
		{"<enter>", "Expand/View"},
		{"<e>", "Edit File"},
		{"<i>", "Terraform Init"},
		{"<p>", "Terraform Plan"},
		{"<a>", "Terraform Apply"},
		{"<d>", "Terraform Destroy"},
		{"<h>", "History"},
		{"<s>", "Settings"},
	})

	generalSection := hv.createSection("GENERAL", []HelpItem{
		{"<esc>", "Back/Clear"},
		{"</>", "Command Mode"},
		{"<shift-c>", "Home (Commands)"},
		{"<shift-b>", "Branch Switch"},
		{"<q>", "Quit"},
		{"<ctrl-c>", "Quit"},
		{"<?>", "Help"},
		{"<shift-h>", "Help"},
	})
	
	historySection := hv.createSection("HISTORY VIEW", []HelpItem{
		{"<h>", "Show History"},
		{"<shift-m>", "Toggle Details"},
		{"<d>", "Load More (3)"},
		{"<u>", "Load Less"},
	})
	
	gitSection := hv.createSection("GIT", []HelpItem{
		{"<shift-b>", "Branch Switch"},
		{"Stash", "Stash & Switch"},
		{"Commit", "Commit & Switch"},
		{"Force", "Discard & Switch"},
	})

	navigationSection := hv.createSection("NAVIGATION", []HelpItem{
		{"<↑/↓>", "Navigate"},
		{"<enter>", "Expand/Collapse"},
		{"<tab>", "Switch Focus"},
	})

	tfSection := hv.createSection("TERRAFORM", []HelpItem{
		{"<i>", "Init"},
		{"<p>", "Plan"},
		{"<a>", "Apply"},
		{"<d>", "Destroy"},
		{"<h>", "Show History"},
	})

	// Combine all sections
	topRow := tview.NewFlex().
		AddItem(resourceSection, 0, 1, false).
		AddItem(generalSection, 0, 1, false)
	topRow.SetBackgroundColor(tcell.ColorBlack)

	bottomRow := tview.NewFlex().
		AddItem(navigationSection, 0, 1, false).
		AddItem(tfSection, 0, 1, false)
	bottomRow.SetBackgroundColor(tcell.ColorBlack)
	
	thirdRow := tview.NewFlex().
		AddItem(historySection, 0, 1, false).
		AddItem(gitSection, 0, 1, false)
	thirdRow.SetBackgroundColor(tcell.ColorBlack)

	hv.SetDirection(tview.FlexRow).
		AddItem(topRow, 0, 1, false).
		AddItem(bottomRow, 0, 1, false).
		AddItem(thirdRow, 0, 1, false)
	hv.SetBackgroundColor(tcell.ColorBlack)
	hv.SetBorder(true)
	hv.SetBorderColor(tcell.NewRGBColor(0, 255, 255))
	hv.SetTitle(" Help ")

	return hv
}

// HelpItem represents a single help item
type HelpItem struct {
	Key         string
	Description string
}

// createSection creates a help section
func (hv *HelpView) createSection(title string, items []HelpItem) *tview.TextView {
	section := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	section.SetBackgroundColor(tcell.ColorBlack)

	fmt.Fprintf(section, "[green::b]%s[white]\n", title)
	for _, item := range items {
		fmt.Fprintf(section, "[blue]%-15s[white] %s\n", item.Key, item.Description)
	}

	return section
}

