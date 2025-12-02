package view

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HeaderView represents the application header
type HeaderView struct {
	*tview.Flex
	currentDir string
	workspace  string
}

// NewHeaderView creates a new header view
func NewHeaderView(currentDir string) *HeaderView {
	hv := &HeaderView{
		Flex:       tview.NewFlex(),
		currentDir: currentDir,
		workspace:  "default",
	}

	hv.buildHeader()
	return hv
}

// buildHeader builds the header components
func (hv *HeaderView) buildHeader() {
	// Info section
	infoText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	infoText.SetBackgroundColor(tcell.ColorBlack)
	
	user := os.Getenv("USER")
	host, _ := os.Hostname()
	
	// Get terraform workspace
	workspace := "default"
	cmd := exec.Command("terraform", "workspace", "show")
	if out, err := cmd.Output(); err == nil {
		workspace = strings.TrimSpace(string(out))
	}
	hv.workspace = workspace

	fmt.Fprintf(infoText, "[cyan]Context:[white]  %s\n", workspace)
	fmt.Fprintf(infoText, "[cyan]Path:[white]     %s\n", hv.currentDir)
	fmt.Fprintf(infoText, "[cyan]User:[white]     %s@%s\n", user, host)
	fmt.Fprintf(infoText, "[cyan]Version:[white]  v0.2.0\n")

	// Shortcuts section
	shortcuts := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	shortcuts.SetBackgroundColor(tcell.ColorBlack)
	
	fmt.Fprintf(shortcuts, "[yellow]<s>[white] Settings    [yellow]<p>[white] Plan\n")
	fmt.Fprintf(shortcuts, "[yellow]<h>[white] History     [yellow]<a>[white] Apply\n")
	fmt.Fprintf(shortcuts, "[yellow]<H>[white] Helm List   [yellow]<e>[white] Edit\n")
	fmt.Fprintf(shortcuts, "[yellow]<Enter>[white] Select  [yellow]<q>[white] Quit")

	// Logo
	logo := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)
	logo.SetBackgroundColor(tcell.ColorBlack)

	fmt.Fprintf(logo, "[cyan] _______ ___      [orange] ______[white]\n")
	fmt.Fprintf(logo, "[cyan]|_   _/ _ \\___    [orange]/ ____|[white]\n")
	fmt.Fprintf(logo, "[cyan]  | | \\_, /(_-<   [orange]`--. \\ [white]\n")
	fmt.Fprintf(logo, "[cyan]  |_|  /_//__ /   [orange]/\\__/ /[white]\n")
	fmt.Fprintf(logo, "[cyan]                 [orange]\\____/ [white]")

	// Combine layouts
	leftFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(infoText, 40, 0, false).
		AddItem(shortcuts, 0, 1, false)
	leftFlex.SetBackgroundColor(tcell.ColorBlack)

	hv.SetDirection(tview.FlexColumn).
		AddItem(leftFlex, 0, 1, false).
		AddItem(logo, 26, 0, false)
	hv.SetBackgroundColor(tcell.ColorBlack)
}

// UpdateWorkspace updates the workspace information
func (hv *HeaderView) UpdateWorkspace(workspace string) {
	hv.workspace = workspace
	hv.buildHeader()
}

