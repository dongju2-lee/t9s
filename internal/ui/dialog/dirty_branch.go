package dialog

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// DirtyBranchDialog handles branch switching when there are local changes
type DirtyBranchDialog struct {
	*tview.Flex
	form         *tview.Form
	onStash      func()
	onCommit     func()
	onForce      func()
	onCancel     func()
	targetBranch string
	modifiedFiles []string
}

// NewDirtyBranchDialog creates a dialog for handling dirty working tree
func NewDirtyBranchDialog(
	currentBranch, targetBranch string,
	modifiedFiles []string,
	onStash func(),
	onCommit func(),
	onForce func(),
	onCancel func(),
) *DirtyBranchDialog {
	dbd := &DirtyBranchDialog{
		Flex:          tview.NewFlex(),
		onStash:       onStash,
		onCommit:      onCommit,
		onForce:       onForce,
		onCancel:      onCancel,
		targetBranch:  targetBranch,
		modifiedFiles: modifiedFiles,
	}

	// Create info text
	infoText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	
	fmt.Fprintf(infoText, "[yellow]⚠️  Working Tree has Uncommitted Changes[white]\n\n")
	fmt.Fprintf(infoText, "Current Branch: [cyan]%s[white]\n", currentBranch)
	fmt.Fprintf(infoText, "Target Branch:  [green]%s[white]\n\n", targetBranch)
	fmt.Fprintf(infoText, "[yellow]Modified Files:[white]\n")
	
	displayFiles := modifiedFiles
	if len(displayFiles) > 10 {
		displayFiles = modifiedFiles[:10]
		fmt.Fprintf(infoText, "[gray]%s[white]\n", strings.Join(displayFiles, "\n"))
		fmt.Fprintf(infoText, "[gray]... and %d more files[white]\n\n", len(modifiedFiles)-10)
	} else {
		fmt.Fprintf(infoText, "[gray]%s[white]\n\n", strings.Join(displayFiles, "\n"))
	}
	
	fmt.Fprintf(infoText, "How would you like to proceed?\n")

	infoText.SetBorder(false)

	// Create form with options
	dbd.form = tview.NewForm().
		AddButton("💾 Stash & Switch", func() {
			if dbd.onStash != nil {
				dbd.onStash()
			}
		}).
		AddButton("📝 Commit & Switch", func() {
			if dbd.onCommit != nil {
				dbd.onCommit()
			}
		}).
		AddButton("⚠️  Force Switch (Discard)", func() {
			if dbd.onForce != nil {
				dbd.onForce()
			}
		}).
		AddButton("Cancel", func() {
			if dbd.onCancel != nil {
				dbd.onCancel()
			}
		})

	dbd.form.SetBorder(false)
	
	dbd.form.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tcell.NewRGBColor(0, 100, 100)).
		SetButtonTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.NewRGBColor(30, 30, 30))

	// Set up key bindings
	dbd.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if dbd.onCancel != nil {
				dbd.onCancel()
			}
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				if dbd.onCancel != nil {
					dbd.onCancel()
				}
				return nil
			case '1':
				if dbd.onStash != nil {
					dbd.onStash()
				}
				return nil
			case '2':
				if dbd.onCommit != nil {
					dbd.onCommit()
				}
				return nil
			case '3':
				if dbd.onForce != nil {
					dbd.onForce()
				}
				return nil
			}
		}
		return event
	})

	// Create container
	container := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(infoText, 0, 1, false).
		AddItem(dbd.form, 3, 0, true)

	container.SetBorder(true).
		SetTitle(" ⚠️  Branch Switch Warning ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(tcell.NewRGBColor(255, 165, 0))

	// Create layout
	dbd.SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false).
			AddItem(container, 80, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	dbd.SetBackgroundColor(tcell.ColorDefault)

	return dbd
}

// GetForm returns the form component
func (dbd *DirtyBranchDialog) GetForm() *tview.Form {
	return dbd.form
}

