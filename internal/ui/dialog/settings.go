package dialog

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/idongju/t9s/internal/config"
)

// SettingsDialog represents the settings dialog
type SettingsDialog struct {
	*tview.Flex
	form   *tview.Form
	config *config.Config
}

// NewSettingsDialog creates a new settings dialog
func NewSettingsDialog(cfg *config.Config, onSave, onCancel func()) *SettingsDialog {
	sd := &SettingsDialog{
		Flex:   tview.NewFlex(),
		config: cfg,
	}

	form := tview.NewForm().
		AddInputField("Terraform Root Directory", cfg.TerraformRoot, 60, nil, func(text string) {
			cfg.TerraformRoot = text
		}).
		AddInputField("Terraform Plan Template", cfg.Commands.PlanTemplate, 60, nil, func(text string) {
			cfg.Commands.PlanTemplate = text
		}).
		AddInputField("Terraform Apply Template", cfg.Commands.ApplyTemplate, 60, nil, func(text string) {
			cfg.Commands.ApplyTemplate = text
		}).
		AddInputField("Default Var File", cfg.Commands.VarFile, 60, nil, func(text string) {
			cfg.Commands.VarFile = text
		}).
		AddButton("Save", func() {
			if err := cfg.Save(); err == nil && onSave != nil {
				onSave()
			}
		}).
		AddButton("Cancel", func() {
			if onCancel != nil {
				onCancel()
			}
		})

	form.SetBackgroundColor(tcell.ColorBlack)
	form.SetBorderColor(tcell.NewRGBColor(0, 255, 255))
	form.SetTitle(" ⚙️  Settings ").SetBorder(true)
	form.SetButtonBackgroundColor(tcell.NewRGBColor(50, 50, 50))
	form.SetButtonTextColor(tcell.NewRGBColor(255, 255, 255))

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	header.SetBackgroundColor(tcell.ColorBlack)
	fmt.Fprintf(header, "[::b][cyan]╔═══════════════════════════════════════════════════════════════════╗\n")
	fmt.Fprintf(header, "[::b][cyan]║[white]  T9s - Settings  [cyan]║\n")
	fmt.Fprintf(header, "[::b][cyan]╚═══════════════════════════════════════════════════════════════════╝")

	help := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	help.SetBackgroundColor(tcell.ColorBlack)
	fmt.Fprintf(help, "\n[yellow]Terraform Root Directory:[white]\n")
	fmt.Fprintf(help, "  Directory where your Terraform code is located (e.g., /home/user/terraform)\n")
	fmt.Fprintf(help, "\n[yellow]Template Variables:[white]\n")
	fmt.Fprintf(help, "  [cyan]{varfile}[white] - Will be replaced with the var file path\n")
	fmt.Fprintf(help, "\n[yellow]Examples:[white]\n")
	fmt.Fprintf(help, "  [green]terraform plan -var-file={varfile}[white]\n")
	fmt.Fprintf(help, "  [green]terraform apply -var-file={varfile} -auto-approve[white]\n")

	sd.SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(form, 0, 1, true).
		AddItem(help, 9, 0, false)
	sd.SetBackgroundColor(tcell.ColorBlack)

	sd.form = form
	return sd
}

// GetForm returns the form component
func (sd *SettingsDialog) GetForm() *tview.Form {
	return sd.form
}

