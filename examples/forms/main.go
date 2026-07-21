package main

import (
	"fmt"
	"os"

	react "github.com/DaviMGDev/bubbletea-react/react"
)

type App struct{}

func (a *App) Render(ctx *react.Context) react.Element {
	name, setName := react.UseState(ctx, "")
	agree, setAgree := react.UseState(ctx, false)
	choice, setChoice := react.UseState(ctx, 0)
	progress, setProgress := react.UseState(ctx, 0)
	tab, setTab := react.UseState(ctx, 0)

	options := []react.SelectOption{
		{Label: "Option A", Value: "a"},
		{Label: "Option B", Value: "b"},
		{Label: "Option C", Value: "c"},
	}

	return react.View(
		react.Box(
			react.Column(
				react.Bold("📋 Forms & Controls Demo"),
				react.Text("Tab/arrows to navigate, Enter/Space to activate."),

				react.Bold("Text Input:"),
				react.Input(name, func(v string) { setName(v) }, "Type here...", 30),

				react.Bold("Checkbox:"),
				react.Checkbox("I agree to the terms", agree, func(v bool) { setAgree(v) }),

				react.Bold("Select:"),
				react.Select(options, choice, func(v string) {
					for i, o := range options {
						if o.Value == v {
							setChoice(i)
							return
						}
					}
				}),

				react.Divider(react.DividerLabel(" Progress ")),
				react.Progress(progress, 10, react.ProgressWidth(30)),
				react.Row(
					react.Button("+1", func() {
						if progress < 10 {
							setProgress(progress + 1)
						}
					}),
					react.Text("  "),
					react.Button("-1", func() {
						if progress > 0 {
							setProgress(progress - 1)
						}
					}),
					react.Text("  "),
					react.Button("Reset", func() { setProgress(0) }),
				),

				react.Divider(react.DividerLabel(" Tabs ")),
				react.Tabs([]react.Tab{
					{Label: "Tab A", Content: react.Text("Content of Tab A")},
					{Label: "Tab B", Content: react.Column(
						react.Text("Content of Tab B"),
						react.Text("- Item 1"),
						react.Text("- Item 2"),
					)},
					{Label: "Tab C", Content: react.Text("Content of Tab C")},
				}, tab, func(i int) { setTab(i) }),
			),
			react.WithBorder(),
			react.WithPadding(1),
			react.WithWidth(70),
			react.WithTitle("Forms Demo"),
		),

		react.Spacer(1),

		// Status line
		react.Textf("Name: %s | Agreed: %v | Choice: %s | Progress: %d/10",
			map[bool]string{true: name, false: "-"}[name != ""],
			map[bool]string{true: "yes", false: "no"}[agree],
			map[bool]string{true: options[choice].Label, false: "-"}[choice >= 0 && choice < len(options)],
			progress,
		),

		react.Spacer(1),
		react.Text("q to quit"),
	)
}

func main() {
	p := react.Root(&App{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
