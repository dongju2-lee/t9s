package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CommitDialog represents a commit message input dialog
type CommitDialog struct {
	*tview.Flex
	form     *tview.Form
	onCommit func(message string)
	onCancel func()
	message  string
}

// NewCommitDialog creates a new commit dialog
func NewCommitDialog(onCommit func(message string), onCancel func()) *CommitDialog {
	cd := &CommitDialog{
		Flex:     tview.NewFlex(),
		onCommit: onCommit,
		onCancel: onCancel,
	}

	// Create form
	cd.form = tview.NewForm().
		AddInputField("Commit Message", "", 60, nil, func(text string) {
			cd.message = text
		}).
		AddButton("Commit & Switch", func() {
			if cd.onCommit != nil && cd.message != "" {
				cd.onCommit(cd.message)
			}
		}).
		AddButton("Cancel", func() {
			if cd.onCancel != nil {
				cd.onCancel()
			}
		})

	cd.form.SetBorder(true).
		SetTitle(" 📝 Commit Changes ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(tcell.NewRGBColor(0, 255, 255))
	
	cd.form.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tcell.NewRGBColor(0, 100, 100)).
		SetButtonTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.NewRGBColor(30, 30, 30))

	// Set up key bindings
	cd.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if cd.onCancel != nil {
				cd.onCancel()
			}
			return nil
		}
		return event
	})

	// Create layout
	cd.SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false).
			AddItem(cd.form, 70, 1, true).
			AddItem(nil, 0, 1, false), 10, 1, true).
		AddItem(nil, 0, 1, false)

	cd.SetBackgroundColor(tcell.ColorDefault)

	return cd
}

// GetForm returns the form component
func (cd *CommitDialog) GetForm() *tview.Form {
	return cd.form
}

