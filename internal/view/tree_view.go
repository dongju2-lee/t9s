package view

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TreeView represents the file tree view
type TreeView struct {
	*tview.TreeView
	currentDir  string
	onFileSelect func(path string)
}

// NewTreeView creates a new tree view
func NewTreeView(rootDir string) *TreeView {
	tv := &TreeView{
		TreeView:   tview.NewTreeView(),
		currentDir: rootDir,
	}

	root := tview.NewTreeNode(filepath.Base(rootDir)).
		SetColor(tcell.NewRGBColor(255, 215, 0)).
		SetReference(rootDir).
		SetSelectable(true)

	tv.SetRoot(root)
	tv.SetCurrentNode(root)
	tv.SetBackgroundColor(tcell.ColorBlack)
	tv.SetBorderColor(tcell.NewRGBColor(0, 255, 255))
	tv.SetTitle(" 📂 File Tree ")
	tv.SetBorder(true)
	tv.SetGraphicsColor(tcell.NewRGBColor(0, 255, 255))

	tv.addTreeChildren(root, rootDir)
	tv.setupHandlers()

	return tv
}

// SetFileSelectHandler sets the handler for file selection
func (tv *TreeView) SetFileSelectHandler(handler func(path string)) {
	tv.onFileSelect = handler
}

// setupHandlers sets up event handlers
func (tv *TreeView) setupHandlers() {
	tv.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return
		}

		path := reference.(string)
		children := node.GetChildren()
		
		if len(children) == 0 {
			info, err := os.Stat(path)
			if err == nil && !info.IsDir() {
				if tv.onFileSelect != nil {
					tv.onFileSelect(path)
				}
			} else if err == nil && info.IsDir() {
				tv.addTreeChildren(node, path)
			}
		} else {
			node.SetChildren([]*tview.TreeNode{})
		}
	})
}

// addTreeChildren adds children to a tree node
func (tv *TreeView) addTreeChildren(target *tview.TreeNode, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") && file.Name() != ".terraform" {
			continue
		}

		fullPath := filepath.Join(path, file.Name())
		node := tview.NewTreeNode(file.Name()).
			SetReference(fullPath).
			SetSelectable(true)

		if file.IsDir() {
			node.SetColor(tcell.NewRGBColor(100, 200, 255))
		} else if strings.HasSuffix(file.Name(), ".tfvars") {
			node.SetColor(tcell.NewRGBColor(255, 100, 255))
		} else if strings.HasSuffix(file.Name(), ".tf") {
			node.SetColor(tcell.NewRGBColor(100, 255, 100))
		} else {
			node.SetColor(tcell.NewRGBColor(200, 200, 200))
		}

		target.AddChild(node)
	}
}

// GetCurrentPath returns the currently selected path
func (tv *TreeView) GetCurrentPath() string {
	node := tv.GetCurrentNode()
	if node == nil {
		return ""
	}
	ref := node.GetReference()
	if ref == nil {
		return ""
	}
	return ref.(string)
}

