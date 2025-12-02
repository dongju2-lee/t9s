package view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CommandView represents the command input bar
type CommandView struct {
	*tview.Flex
	input      *tview.InputField
	pathLabel  *tview.TextView
	onExecute  func(cmd string)
	currentDir string
}

// NewCommandView creates a new command input view
func NewCommandView(currentDir string) *CommandView {
	cv := &CommandView{
		Flex:       tview.NewFlex(),
		currentDir: currentDir,
	}

	// Path label showing current directory
	cv.pathLabel = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	cv.pathLabel.SetBackgroundColor(tcell.ColorBlack)
	cv.UpdatePath(currentDir)

	// Command input field
	cv.input = tview.NewInputField().
		SetLabel("$ ").
		SetFieldBackgroundColor(tcell.NewRGBColor(30, 30, 30)).
		SetFieldTextColor(tcell.ColorWhite)
	cv.input.SetBackgroundColor(tcell.ColorBlack)
	cv.input.SetLabelColor(tcell.NewRGBColor(0, 255, 0))

	cv.SetDirection(tview.FlexRow).
		AddItem(cv.pathLabel, 1, 0, false).
		AddItem(cv.input, 1, 0, true)
	cv.SetBackgroundColor(tcell.ColorBlack)

	return cv
}

// UpdatePath updates the current directory path
func (cv *CommandView) UpdatePath(path string) {
	cv.currentDir = path
	cv.pathLabel.Clear()
	fmt.Fprintf(cv.pathLabel, "[cyan]Working Directory:[white] [yellow]%s[white]", path)
}

// SetExecuteHandler sets the handler for command execution
func (cv *CommandView) SetExecuteHandler(handler func(cmd string)) {
	cv.onExecute = handler
	cv.input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter && handler != nil {
			cmd := cv.input.GetText()
			if cmd != "" {
				handler(cmd)
			}
		}
	})
}

// GetInput returns the input field
func (cv *CommandView) GetInput() *tview.InputField {
	return cv.input
}

// Clear clears the input field
func (cv *CommandView) Clear() {
	cv.input.SetText("")
}

// GetCurrentDir returns the current directory
func (cv *CommandView) GetCurrentDir() string {
	return cv.currentDir
}

