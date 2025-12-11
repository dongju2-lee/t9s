package dialog

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// BranchDialog represents a branch selection dialog
type BranchDialog struct {
	*tview.Flex
	list          *tview.List
	onSelect      func(branch string)
	onCancel      func()
	currentBranch string
}

// NewBranchDialog creates a new branch selection dialog
func NewBranchDialog(branches []string, currentBranch string, onSelect func(branch string), onCancel func()) *BranchDialog {
	bd := &BranchDialog{
		Flex:          tview.NewFlex(),
		onSelect:      onSelect,
		onCancel:      onCancel,
		currentBranch: currentBranch,
	}

	// Create list
	bd.list = tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)

	// Add branches to list
	for _, branch := range branches {
		branchName := branch
		if branch == currentBranch {
			bd.list.AddItem(fmt.Sprintf("* %s (current)", branch), "", 0, func() {
				if bd.onSelect != nil {
					bd.onSelect(branchName)
				}
			})
		} else {
			bd.list.AddItem(fmt.Sprintf("  %s", branch), "", 0, func() {
				if bd.onSelect != nil {
					bd.onSelect(branchName)
				}
			})
		}
	}

	bd.list.SetBorder(true).
		SetTitle(" 🌿 Select Branch ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.NewRGBColor(0, 255, 255))

	// Set up key bindings
	bd.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if bd.onCancel != nil {
				bd.onCancel()
			}
			return nil
		case tcell.KeyRune:
			if event.Rune() == 'q' || event.Rune() == 'Q' {
				if bd.onCancel != nil {
					bd.onCancel()
				}
				return nil
			}
		}
		return event
	})

	// Create layout
	bd.SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false).
			AddItem(bd.list, 60, 1, true).
			AddItem(nil, 0, 1, false), 20, 1, true).
		AddItem(nil, 0, 1, false)

	bd.SetBackgroundColor(tcell.ColorDefault)

	return bd
}

// GetList returns the list component
func (bd *BranchDialog) GetList() *tview.List {
	return bd.list
}


