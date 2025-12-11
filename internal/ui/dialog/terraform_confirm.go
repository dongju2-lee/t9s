package dialog

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TerraformConfirmDialog creates a detailed confirmation dialog for terraform commands
type TerraformConfirmDialog struct {
	*tview.Flex
}

// NewTerraformConfirmDialog creates a new terraform confirmation dialog
func NewTerraformConfirmDialog(command, workDir, configFile, fileContent string, onExecute, onAutoApprove, onCancel func()) *TerraformConfirmDialog {
	td := &TerraformConfirmDialog{
		Flex: tview.NewFlex(),
	}

	// Header
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	header.SetBackgroundColor(tcell.ColorBlack)
	fmt.Fprintf(header, "[::b][yellow]╔═══════════════════════════════════════════════════════════════════╗\n")
	fmt.Fprintf(header, "[::b][yellow]║[white]  Terraform Command Confirmation  [yellow]║\n")
	fmt.Fprintf(header, "[::b][yellow]╚═══════════════════════════════════════════════════════════════════╝")

	// Info section
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	info.SetBackgroundColor(tcell.ColorBlack)
	fmt.Fprintf(info, "\n[cyan]Command:[white] %s\n", command)
	fmt.Fprintf(info, "[cyan]Directory:[white] %s\n", workDir)
	fmt.Fprintf(info, "[cyan]Config File:[white] %s\n\n", configFile)
	fmt.Fprintf(info, "[yellow]Config Content:[white]\n")
	fmt.Fprintf(info, "[green]%s[white]\n", strings.Repeat("─", 65))
	if fileContent != "" {
		fmt.Fprintf(info, "%s\n", fileContent)
	} else {
		fmt.Fprintf(info, "[gray](empty or file not found)[white]\n")
	}
	fmt.Fprintf(info, "[green]%s[white]\n\n", strings.Repeat("─", 65))

	// Question
	question := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	question.SetBackgroundColor(tcell.ColorBlack)
	fmt.Fprintf(question, "[yellow]Do you want to proceed with this command?[white]\n")
	fmt.Fprintf(question, "[gray](Execute: manual 'yes' required | Auto Approve: automatic execution)[white]\n")

	// Buttons
	form := tview.NewForm().
		AddButton("Execute", func() {
			if onExecute != nil {
				onExecute()
			}
		}).
		AddButton("Auto Approve", func() {
			if onAutoApprove != nil {
				onAutoApprove()
			}
		}).
		AddButton("Cancel", func() {
			if onCancel != nil {
				onCancel()
			}
		})
	form.SetBackgroundColor(tcell.ColorBlack)
	form.SetButtonBackgroundColor(tcell.NewRGBColor(50, 50, 50))
	form.SetButtonTextColor(tcell.ColorWhite)
	form.SetButtonsAlign(tview.AlignCenter)

	td.SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(info, 0, 1, false).
		AddItem(question, 2, 0, false).
		AddItem(form, 3, 0, true)
	td.SetBackgroundColor(tcell.ColorBlack)
	td.SetBorder(true)
	td.SetBorderColor(tcell.NewRGBColor(255, 165, 0)) // Orange for warning
	td.SetTitle(" ⚠️  Confirmation Required ")

	return td
}

// GetForm returns the form for focus management
func (td *TerraformConfirmDialog) GetForm() *tview.Form {
	// Get the form from the flex layout
	item := td.GetItem(3)
	if form, ok := item.(*tview.Form); ok {
		return form
	}
	return nil
}

