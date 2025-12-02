package components

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
	"github.com/idongju/t9s/internal/config"
	"github.com/idongju/t9s/internal/view"
)

// CommandExecutor handles command execution
type CommandExecutor struct {
	app         *tview.Application
	contentView *view.ContentView
	config      *config.Config
}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor(app *tview.Application, contentView *view.ContentView, cfg *config.Config) *CommandExecutor {
	return &CommandExecutor{
		app:         app,
		contentView: contentView,
		config:      cfg,
	}
}

// ExecutePlan executes terraform plan
func (ce *CommandExecutor) ExecutePlan(path string) {
	ce.runTerraformCommand("Plan", path, ce.config.Commands.PlanTemplate)
}

// ExecuteApply executes terraform apply
func (ce *CommandExecutor) ExecuteApply(path string) {
	ce.runTerraformCommand("Apply", path, ce.config.Commands.ApplyTemplate)
}

// ShowHistory shows terraform history
func (ce *CommandExecutor) ShowHistory(path string) {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		path = filepath.Dir(path)
	}

	ce.contentView.Clear()
	ce.contentView.SetTitle(" ⏰ Terraform History ")
	fmt.Fprintf(ce.contentView, "[yellow]Terraform History[white]\n")
	fmt.Fprintf(ce.contentView, "[cyan]Directory:[white] %s\n", path)
	fmt.Fprintf(ce.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))

	statePath := filepath.Join(path, "terraform.tfstate")
	stateInfo, err := os.Stat(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(ce.contentView, "[yellow]No local terraform.tfstate found.[white]\n")
			fmt.Fprintf(ce.contentView, "This might be using a remote backend or hasn't been initialized.\n\n")
		} else {
			fmt.Fprintf(ce.contentView, "[red]Error checking state file:[white] %v\n", err)
		}
	} else {
		fmt.Fprintf(ce.contentView, "[green]Last Local Apply:[white] %s\n", stateInfo.ModTime().Format("2006-01-02 15:04:05"))
		fmt.Fprintf(ce.contentView, "[gray](Based on terraform.tfstate modification time)[white]\n\n")
	}

	fmt.Fprintf(ce.contentView, "[cyan]Current State Summary (terraform show):[white]\n")
	go func() {
		cmd := exec.Command("terraform", "show", "-no-color")
		cmd.Dir = path
		output, err := cmd.CombinedOutput()
		
		ce.app.QueueUpdateDraw(func() {
			if err != nil {
				fmt.Fprintf(ce.contentView, "[red]Error executing terraform show:[white] %v\n", err)
			}
			ce.contentView.AppendText(string(output))
		})
	}()
}

// ShowHelm shows helm list output
func (ce *CommandExecutor) ShowHelm() {
	ce.contentView.Clear()
	ce.contentView.SetTitle(" ⎈ Helm Releases ")
	fmt.Fprintf(ce.contentView, "[yellow]Helm List -A[white]\n\n")
	fmt.Fprintf(ce.contentView, "[cyan]Executing:[white] helm list -A\n")
	fmt.Fprintf(ce.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))

	cmd := exec.Command("helm", "list", "-A")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(ce.contentView, "[red]Error executing helm:[white] %v\n\n", err)
		fmt.Fprintf(ce.contentView, "Make sure 'helm' is installed and available in your PATH.\n")
		if len(output) > 0 {
			fmt.Fprintf(ce.contentView, "\n[red]Output:[white]\n%s", string(output))
		}
		return
	}

	ce.contentView.AppendText(string(output))
}

// EditFile opens a file in the default editor
func (ce *CommandExecutor) EditFile(filePath string) {
	if filePath == "" {
		ce.contentView.DisplayText("No File", "[yellow]No file selected[white]\n\nPlease select a file from the tree first.")
		return
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	ce.app.Suspend(func() {
		cmd := exec.Command(editor, filePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running editor: %v\nPress Enter to continue...", err)
			fmt.Scanln()
		}
	})
}

// runTerraformCommand executes a terraform command
func (ce *CommandExecutor) runTerraformCommand(action string, path string, template string) {
	var workDir string
	var varFile string

	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		workDir = filepath.Dir(path)
		if strings.HasSuffix(path, ".tfvars") {
			varFile = path
		}
	} else {
		workDir = path
		defaultVar := filepath.Join(workDir, ce.config.Commands.TfvarsFile)
		if _, err := os.Stat(defaultVar); err == nil {
			varFile = defaultVar
		}
	}

	ce.contentView.Clear()
	ce.contentView.SetTitle(fmt.Sprintf(" 🚀 Terraform %s ", action))
	fmt.Fprintf(ce.contentView, "[yellow]Executing Terraform %s[white]\n", action)
	fmt.Fprintf(ce.contentView, "[cyan]Directory:[white] %s\n", workDir)
	if varFile != "" {
		fmt.Fprintf(ce.contentView, "[cyan]Var File:[white] %s\n", varFile)
	} else {
		fmt.Fprintf(ce.contentView, "[yellow]Warning:[white] No .tfvars file found or selected. Running without -var-file.\n")
	}
	fmt.Fprintf(ce.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))

	cmdStr := template
	if varFile != "" {
		cmdStr = strings.ReplaceAll(cmdStr, "{varfile}", varFile)
	} else {
		cmdStr = strings.ReplaceAll(cmdStr, "-var-file={varfile}", "")
	}

	fmt.Fprintf(ce.contentView, "[gray]Command: %s[white]\n\n", cmdStr)

	go func() {
		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			return
		}
		
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Dir = workDir
		
		output, err := cmd.CombinedOutput()
		
		ce.app.QueueUpdateDraw(func() {
			if err != nil {
				fmt.Fprintf(ce.contentView, "[red]Error executing command:[white] %v\n", err)
			}
			ce.contentView.AppendText(string(output))
			ce.contentView.AppendText("\n\n[green]Done.[white]")
		})
	}()
}

