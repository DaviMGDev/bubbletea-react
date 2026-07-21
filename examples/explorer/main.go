package main

import (
	"fmt"
	"os"

	react "github.com/DaviMGDev/bubbletea-react/react"
)

// ============================================================================
// Data model
// ============================================================================

type TreeNode struct {
	Name     string
	Type     string // "file" or "dir"
	Children []TreeNode
	Size     string
}

var sampleTree = []TreeNode{
	{
		Name: "src", Type: "dir", Children: []TreeNode{
			{Name: "main.go", Type: "file", Size: "2.1 KB"},
			{Name: "utils.go", Type: "file", Size: "4.3 KB"},
			{Name: "components", Type: "dir", Children: []TreeNode{
				{Name: "header.go", Type: "file", Size: "1.8 KB"},
				{Name: "footer.go", Type: "file", Size: "0.9 KB"},
				{Name: "sidebar.go", Type: "file", Size: "3.2 KB"},
			}},
			{Name: "api", Type: "dir", Children: []TreeNode{
				{Name: "handler.go", Type: "file", Size: "5.1 KB"},
				{Name: "router.go", Type: "file", Size: "2.7 KB"},
			}},
		},
	},
	{
		Name: "docs", Type: "dir", Children: []TreeNode{
			{Name: "README.md", Type: "file", Size: "0.5 KB"},
			{Name: "CONTRIBUTING.md", Type: "file", Size: "1.2 KB"},
		},
	},
	{Name: "Makefile", Type: "file", Size: "0.3 KB"},
	{Name: "package.json", Type: "file", Size: "0.2 KB"},
}

func flattenTree(nodes []TreeNode, prefix string) [][]string {
	var rows [][]string
	for _, n := range nodes {
		path := prefix + n.Name
		if n.Type == "file" {
			rows = append(rows, []string{path, n.Size})
		}
		if len(n.Children) > 0 {
			rows = append(rows, flattenTree(n.Children, path+"/")...)
		}
	}
	return rows
}

// ============================================================================
// Nested component: TreeNodeComponent
// ============================================================================

type TreeNodeComponent struct {
	Node  TreeNode
	Depth int
}

func (c *TreeNodeComponent) Render(ctx *react.Context) react.Element {
	expanded, setExpanded := react.UseState(ctx, false)

	isDir := len(c.Node.Children) > 0
	indent := ""
	for i := 0; i < c.Depth; i++ {
		indent += "  "
	}

	toggle := func() {
		if isDir {
			setExpanded(!expanded)
		}
	}

	icon := "📄"
	if isDir {
		if expanded {
			icon = "📂"
		} else {
			icon = "📁"
		}
	}

	label := fmt.Sprintf("%s%s %s", indent, icon, c.Node.Name)

	children := []react.Element{
		react.Row(
			react.Text(label),
			react.Text(" "),
			react.Button("⏺", toggle),
		),
	}

	if expanded && isDir {
		for _, child := range c.Node.Children {
			key := fmt.Sprintf("node-%d-%s", c.Depth, child.Name)
			children = append(children,
				react.Wrap(&TreeNodeComponent{Node: child, Depth: c.Depth + 1}, key),
			)
		}
	}

	return react.Column(children...)
}

// ============================================================================
// Main app — compact layout for 41-row terminals
// ============================================================================

type ExplorerDemo struct{}

func (d *ExplorerDemo) Render(ctx *react.Context) react.Element {
	allFiles := flattenTree(sampleTree, "")

	tableColumns := []react.TableColumn{
		{Label: "File Path", Width: 28},
		{Label: "Size", Width: 10},
	}

	tableEl := react.Element(react.Text("  (no files)"))
	if len(allFiles) > 0 {
		tableEl = react.Table(tableColumns, allFiles)
	}

	var treeEls []react.Element
	for _, node := range sampleTree {
		treeEls = append(treeEls,
			react.Wrap(&TreeNodeComponent{Node: node, Depth: 0}, "root-"+node.Name),
		)
	}

	// No outer Box — just direct content. Saves 4 lines.
	return react.View(
		react.Column(
			react.Bold("🗂 File Tree Explorer"),
			react.Text("Deep Wrap nesting • Scrollable • Table • Per-component state isolation"),
			react.Row(
				react.Box(
					react.Column(
						react.Bold("Folders"),
						react.Column(treeEls...),
					),
					react.WithBorder(), react.WithWidth(34), react.WithTitle("Tree"),
				),
				react.Text(" "),
				react.Box(
					react.Column(
						react.Bold("All Files"),
						tableEl,
					),
					react.WithBorder(), react.WithWidth(40), react.WithTitle("Details"),
				),
			),
			react.Text("Click ⏺ to expand/collapse folders • Each node has its own expand state"),
		),
	)
}

func main() {
	p := react.Root(&ExplorerDemo{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
