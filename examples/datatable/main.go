package main

import (
	"fmt"
	"os"

	react "github.com/DaviMGDev/bubbletea-react/react"
)

// --- Data types ---

type UserRecord struct {
	Name   string
	Role   string
	Status string
	Email  string
}

// sampleData is the dataset displayed in the table.
var sampleData = []UserRecord{
	{Name: "Alice Johnson", Role: "Admin", Status: "Active", Email: "alice@example.com"},
	{Name: "Bob Smith", Role: "Editor", Status: "Active", Email: "bob@example.com"},
	{Name: "Carol Davis", Role: "Viewer", Status: "Inactive", Email: "carol@example.com"},
	{Name: "Dave Wilson", Role: "Editor", Status: "Active", Email: "dave@example.com"},
	{Name: "Eve Martin", Role: "Admin", Status: "Suspended", Email: "eve@example.com"},
	{Name: "Frank Lee", Role: "Viewer", Status: "Active", Email: "frank@example.com"},
}

// --- Nested component: UserDetailPane ---

// UserDetailPane is a nested component rendered via Wrap().
// It has its OWN state (selectedIndex), independent of the parent.
type UserDetailPane struct {
	Users []UserRecord
}

func (p *UserDetailPane) Render(ctx *react.Context) react.Element {
	// This UseState is local to each UserDetailPane instance.
	// The hook index is scoped to this component's Context.
	selected, setSelected := react.UseState(ctx, 0)

	if len(p.Users) == 0 {
		return react.Text("  (no users)")
	}

	// Clamp selection.
	if selected < 0 {
		selected = 0
	}
	if selected >= len(p.Users) {
		selected = len(p.Users) - 1
	}

	user := p.Users[selected]

	return react.Column(
		react.Bold("Detail View"),
		react.Spacer(1),
		react.Textf("Name:   %s", user.Name),
		react.Textf("Role:   %s", user.Role),
		react.Textf("Status: %s", user.Status),
		react.Textf("Email:  %s", user.Email),
		react.Spacer(1),
		react.Text("Navigate with ↑/↓ to change selection."),
		react.Spacer(1),

		// Navigation buttons for this nested component's state.
		react.Row(
			react.Text("Select: "),
			react.Button("◀", func() {
				if selected > 0 {
					setSelected(selected - 1)
				}
			}),
			react.Textf(" %d/%d ", selected+1, len(p.Users)),
			react.Button("▶", func() {
				if selected < len(p.Users)-1 {
					setSelected(selected + 1)
				}
			}),
		),
	)
}

// --- Main app ---

// DataTableDemo demonstrates Table, Tabs, and Wrap (nested components).
type DataTableDemo struct{}

func (d *DataTableDemo) Render(ctx *react.Context) react.Element {
	activeTab, setActiveTab := react.UseState(ctx, 0)

	// Build table columns and rows from sample data.
	columns := []react.TableColumn{
		{Label: "Name", Width: 18},
		{Label: "Role", Width: 10},
		{Label: "Status", Width: 12},
		{Label: "Email", Width: 22},
	}

	rows := make([][]string, 0, len(sampleData))
	for _, u := range sampleData {
		rows = append(rows, []string{u.Name, u.Role, u.Status, u.Email})
	}

	// Table placeholder when no data.
	tableEl := react.Element(react.Text("(no data)"))
	if len(rows) > 0 {
		tableEl = react.Table(columns, rows)
	}

	return react.View(
		react.Box(
			react.Column(
				react.Bold("📊 Dataset Viewer"),
				react.Text("Features: Table, Tabs, Wrap (nested components), per-component state isolation."),
				react.Spacer(1),

				// Tabs with two panes.
				react.Tabs([]react.Tab{
					{
						Label: "Data Grid",
						Content: react.Column(
							react.Spacer(1),
							tableEl,
						),
					},
					{
						Label: "Detail",
						Content: react.Column(
							react.Spacer(1),
							// Wrap creates a nested component with its own Context.
							// The second argument ("detail") is a stable key for
							// reconciliation — it preserves state across re-renders
							// of the parent.
							react.Wrap(&UserDetailPane{Users: sampleData}, "detail"),
						),
					},
				}, activeTab, func(i int) { setActiveTab(i) }),

				react.Spacer(1),
				react.Text("←/→ switch tabs • Enter/Space activates tab • Tab to navigate"),
			),
			react.WithBorder(),
			react.WithPadding(1),
			react.WithWidth(76),
			react.WithTitle("Data Table Demo"),
		),
	)
}

func main() {
	p := react.Root(&DataTableDemo{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
