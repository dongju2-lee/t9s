package view

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ContentView represents the content display area
type ContentView struct {
	*tview.TextView
}

// NewContentView creates a new content view
func NewContentView() *ContentView {
	cv := &ContentView{
		TextView: tview.NewTextView().
			SetDynamicColors(true).
			SetScrollable(true).
			SetWordWrap(true),
	}

	cv.SetBackgroundColor(tcell.ColorBlack)
	cv.SetBorderColor(tcell.NewRGBColor(0, 255, 255))
	cv.SetTitle(" 📄 Content ")
	cv.SetBorder(true)

	cv.ShowWelcome()

	return cv
}

// ShowWelcome displays the welcome message
func (cv *ContentView) ShowWelcome() {
	cv.Clear()
	cv.SetTitle(" 📄 Content ")
	fmt.Fprintf(cv, "[yellow]Welcome to T9s![white]\n\n")
	fmt.Fprintf(cv, "Select a file from the tree to view its content.\n\n")
	fmt.Fprintf(cv, "[cyan]Available Commands:[white]\n")
	fmt.Fprintf(cv, "  • [green]h[white] - View terraform history\n")
	fmt.Fprintf(cv, "  • [green]p[white] - Terraform plan\n")
	fmt.Fprintf(cv, "  • [green]a[white] - Terraform apply\n")
	fmt.Fprintf(cv, "  • [green]e[white] - Edit current file\n")
	fmt.Fprintf(cv, "  • [green]s[white] - Settings\n")
	fmt.Fprintf(cv, "  • [green]?[white] or [green]Shift+H[white] - Help\n")
	fmt.Fprintf(cv, "  • [green]/[white] - Command mode\n")
	fmt.Fprintf(cv, "  • [green]Shift+C[white] - Show this screen\n")
	fmt.Fprintf(cv, "  • [green]q[white] - Quit\n")
}

// DisplayFile displays the content of a file
func (cv *ContentView) DisplayFile(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		cv.Clear()
		fmt.Fprintf(cv, "[red]Error reading file: %v[white]", err)
		return err
	}

	cv.Clear()
	cv.SetTitle(fmt.Sprintf(" 📄 %s ", filepath.Base(path)))
	
	fmt.Fprintf(cv, "[yellow]File:[white] %s\n", path)
	fmt.Fprintf(cv, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))
	fmt.Fprintf(cv, "%s", string(content))

	if strings.HasSuffix(path, ".tfvars") || strings.HasSuffix(path, ".tf") {
		fmt.Fprintf(cv, "\n\n[gray]Press 'e' to edit this file[white]")
	}

	return nil
}

// DisplayText displays arbitrary text with a title
func (cv *ContentView) DisplayText(title, content string) {
	cv.Clear()
	cv.SetTitle(fmt.Sprintf(" %s ", title))
	fmt.Fprintf(cv, "%s", content)
}

// AppendText appends text to the current content
func (cv *ContentView) AppendText(text string) {
	fmt.Fprintf(cv, "%s", text)
}

