package dialog

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ConfirmDialog creates a confirmation dialog
type ConfirmDialog struct {
	*tview.Modal
}

// NewConfirmDialog creates a new confirmation dialog
func NewConfirmDialog(text string, onConfirm, onCancel func()) *ConfirmDialog {
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"Yes", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" && onConfirm != nil {
				onConfirm()
			} else if onCancel != nil {
				onCancel()
			}
		})

	modal.SetBackgroundColor(tcell.ColorBlack)
	modal.SetTextColor(tcell.ColorWhite)
	modal.SetButtonBackgroundColor(tcell.NewRGBColor(50, 50, 50))
	modal.SetButtonTextColor(tcell.ColorWhite)

	return &ConfirmDialog{Modal: modal}
}


