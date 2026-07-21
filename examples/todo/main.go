package main

import (
	"fmt"
	"os"
	"strings"

	react "github.com/DaviMGDev/bubbletea-react/react"
)

// TodoItem represents a single todo item.
type TodoItem struct {
	Text string
	Done bool
}

// TodoApp is the main todo list component.
type TodoApp struct{}

func (t *TodoApp) Render(ctx *react.Context) react.Element {
	items, setItems := react.UseState(ctx, []TodoItem{})
	input, setInput := react.UseState(ctx, "")

	return react.View(
		react.Box(
			react.Column(
				react.Bold("📋 Todo List"),
				react.Spacer(1),

				// Input row
				react.Row(
					react.Input(input, func(val string) { setInput(val) }, "Add a todo...", 30),
					react.Text("  "),
					react.Button("Add", func() {
						trimmed := strings.TrimSpace(input)
						if trimmed != "" {
							setItems(append(items, TodoItem{Text: trimmed, Done: false}))
							setInput("")
						}
					}),
				),

				react.Spacer(1),

				// Todo list items
				renderTodoList(items, setItems),

				react.Spacer(1),

				// Footer
				func() react.Element {
					doneCount := 0
					for _, item := range items {
						if item.Done {
							doneCount++
						}
					}
					return react.Textf("%d / %d completed", doneCount, len(items))
				}(),
			),
			react.WithBorder(),
			react.WithPadding(1),
			react.WithWidth(60),
			react.WithTitle("Todo App"),
		),
	)
}

// renderTodoList renders the list of todo items with toggle and delete buttons.
func renderTodoList(items []TodoItem, setItems func([]TodoItem)) react.Element {
	if len(items) == 0 {
		return react.Text("  Nothing yet! Add a todo above.")
	}

	children := make([]react.Element, 0, len(items)*2)
	for i, item := range items {
		idx := i // capture for closure
		prefix := "[ ]"
		textStyle := item.Text
		if item.Done {
			prefix = "[✓]"
			textStyle = item.Text + " (done)"
		}

		children = append(children, react.Row(
			react.Text(fmt.Sprintf("  %s ", prefix)),
			react.Text(textStyle),
			react.Text("  "),
			react.Button("Toggle", func() {
				newItems := make([]TodoItem, len(items))
				copy(newItems, items)
				newItems[idx].Done = !newItems[idx].Done
				setItems(newItems)
			}),
			react.Text(" "),
			react.Button("Del", func() {
				newItems := make([]TodoItem, 0, len(items)-1)
				newItems = append(newItems, items[:idx]...)
				newItems = append(newItems, items[idx+1:]...)
				setItems(newItems)
			}),
		))

		// Add a spacer between items
		if i < len(items)-1 {
			children = append(children, react.Spacer(0))
		}
	}

	return react.Column(children...)
}

func main() {
	p := react.Root(&TodoApp{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
