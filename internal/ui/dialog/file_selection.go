package dialog

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FileSelectionDialog shows a dialog to select a config file
type FileSelectionDialog struct {
	*tview.Flex
	list        *tview.List
	preview     *tview.TextView
	configDir   string
	filePattern string
	onSelect    func(string, string) // callback with (filepath, content)
	onCancel    func()
}

// NewFileSelectionDialog creates a new file selection dialog
func NewFileSelectionDialog(configDir, filePattern, title string, onSelect func(string, string), onCancel func()) *FileSelectionDialog {
	dialog := &FileSelectionDialog{
		configDir:   configDir,
		filePattern: filePattern,
		onSelect:    onSelect,
		onCancel:    onCancel,
	}

	// Create list for files
	dialog.list = tview.NewList()
	dialog.list.SetBorder(true).SetTitle(" 📁 Available Files ")
	dialog.list.SetBackgroundColor(tcell.ColorBlack)
	dialog.list.SetBorderColor(tcell.NewRGBColor(0, 255, 255))
	dialog.list.ShowSecondaryText(false)

	// Create preview pane
	dialog.preview = tview.NewTextView()
	dialog.preview.SetBorder(true).SetTitle(" 📄 File Preview ")
	dialog.preview.SetBackgroundColor(tcell.ColorBlack)
	dialog.preview.SetBorderColor(tcell.NewRGBColor(0, 255, 255))
	dialog.preview.SetDynamicColors(true)
	dialog.preview.SetScrollable(true)
	dialog.preview.SetWrap(true)

	// Load files
	dialog.loadFiles()

	// Handle selection change
	dialog.list.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		dialog.showPreview(mainText)
	})

	// Handle selection
	dialog.list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		content := dialog.getFileContent(mainText)
		if dialog.onSelect != nil {
			filePath := filepath.Join(dialog.configDir, mainText)
			dialog.onSelect(filePath, content)
		}
	})

	// Create title
	titleText := tview.NewTextView()
	titleText.SetTextAlign(tview.AlignCenter)
	titleText.SetDynamicColors(true)
	titleText.SetBackgroundColor(tcell.ColorBlack)
	fmt.Fprintf(titleText, "[::b][yellow]%s[white]\n\n", title)
	fmt.Fprintf(titleText, "[gray]Use ↑↓ to navigate, Enter to select, Esc to cancel[white]")

	// Create help text
	helpText := tview.NewTextView()
	helpText.SetTextAlign(tview.AlignCenter)
	helpText.SetDynamicColors(true)
	helpText.SetBackgroundColor(tcell.ColorBlack)
	fmt.Fprintf(helpText, "\n[yellow]↑↓[white] Navigate  [yellow]Enter[white] Select  [yellow]Esc[white] Cancel")

	// Create content layout (list + preview)
	contentFlex := tview.NewFlex().
		AddItem(dialog.list, 0, 1, true).
		AddItem(dialog.preview, 0, 2, false)
	contentFlex.SetBackgroundColor(tcell.ColorBlack)

	// Main layout
	dialog.Flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(titleText, 3, 0, false).
		AddItem(contentFlex, 0, 1, true).
		AddItem(helpText, 2, 0, false)
	dialog.Flex.SetBackgroundColor(tcell.ColorBlack)

	// Set up input capture
	dialog.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			if dialog.onCancel != nil {
				dialog.onCancel()
			}
			return nil
		}
		return event
	})

	// Show first file preview
	if dialog.list.GetItemCount() > 0 {
		mainText, _ := dialog.list.GetItemText(0)
		dialog.showPreview(mainText)
	}

	return dialog
}

// loadFiles loads files from config directory
func (d *FileSelectionDialog) loadFiles() {
	files, err := ioutil.ReadDir(d.configDir)
	if err != nil {
		d.list.AddItem("(Error reading directory)", "", 0, nil)
		return
	}

	count := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Check if file matches pattern
		if d.filePattern != "" {
			matched, _ := filepath.Match(d.filePattern, file.Name())
			if !matched {
				continue
			}
		}

		d.list.AddItem(file.Name(), "", 0, nil)
		count++
	}

	if count == 0 {
		d.list.AddItem("(No files found)", "", 0, nil)
		d.preview.SetText("[yellow]No configuration files found in this directory.[white]")
	}
}

// showPreview shows file content in preview pane
func (d *FileSelectionDialog) showPreview(filename string) {
	content := d.getFileContent(filename)
	if content == "" {
		d.preview.SetText("[yellow]No preview available[white]")
		return
	}

	d.preview.Clear()
	
	// Add file info header
	filePath := filepath.Join(d.configDir, filename)
	fmt.Fprintf(d.preview, "[cyan]File:[white] %s\n", filename)
	fmt.Fprintf(d.preview, "[cyan]Path:[white] %s\n", filePath)
	fmt.Fprintf(d.preview, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))
	
	// Add content with syntax highlighting
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			// Comment
			fmt.Fprintf(d.preview, "[gray]%s[white]\n", line)
		} else if strings.Contains(line, "=") {
			// Variable assignment
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				fmt.Fprintf(d.preview, "[yellow]%s[white]=[green]%s[white]\n", parts[0], parts[1])
			} else {
				fmt.Fprintf(d.preview, "%s\n", line)
			}
		} else {
			fmt.Fprintf(d.preview, "%s\n", line)
		}
	}
}

// getFileContent reads and returns file content
func (d *FileSelectionDialog) getFileContent(filename string) string {
	if filename == "(No files found)" || filename == "(Error reading directory)" {
		return ""
	}

	filePath := filepath.Join(d.configDir, filename)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("[red]Error reading file: %v[white]", err)
	}

	return string(content)
}

// GetList returns the list component for focus management
func (d *FileSelectionDialog) GetList() *tview.List {
	return d.list
}
