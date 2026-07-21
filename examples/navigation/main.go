package main

import (
	"fmt"
	"os"

	react "github.com/DaviMGDev/bubbletea-react/react"
)

// NavigationDemo demonstrates ButtonWithOptions, InputWithOptions,
// and the spatial navigation focus/blur callback system.
type NavigationDemo struct{}

func (d *NavigationDemo) Render(ctx *react.Context) react.Element {
	// Track focus status for display.
	focusStatus, setFocusStatus := react.UseState(ctx, "(no element focused)")
	inputValue, setInputValue := react.UseState(ctx, "hello")
	clickLog, setClickLog := react.UseState(ctx, "(none)")

	// Helper to update the status line.
	onFocus := func(name string) func() {
		return func() { setFocusStatus("FOCUS: " + name) }
	}
	onBlur := func(name string) func() {
		return func() { setFocusStatus("BLUR: " + name) }
	}

	return react.View(
		react.Box(
			react.Column(
				react.Bold("🎮 Focus Navigation Demo"),
				react.Text("Features: ButtonWithOptions, InputWithOptions, OnFocus/OnBlur callbacks."),
				react.Spacer(1),
				react.Text("Tab/arrows move focus. Observe the status line below."),
				react.Divider(react.DividerLabel(" Interactive Elements ")),

				// Buttons with focus/blur callbacks.
				// ButtonWithOptions extends Button with optional OnFocus/OnBlur hooks.
				react.Bold("Buttons:"),
				react.Row(
					react.ButtonWithOptions("Button A", func() {
						setClickLog("Clicked: Button A")
					},
						react.WithOnFocus(onFocus("Button A")),
						react.WithOnBlur(onBlur("Button A")),
					),
					react.Text("  "),
					react.ButtonWithOptions("Button B", func() {
						setClickLog("Clicked: Button B")
					},
						react.WithOnFocus(onFocus("Button B")),
						react.WithOnBlur(onBlur("Button B")),
					),
					react.Text("  "),
					react.ButtonWithOptions("Button C", func() {
						setClickLog("Clicked: Button C")
					},
						react.WithOnFocus(onFocus("Button C")),
						react.WithOnBlur(onBlur("Button C")),
					),
				),

				react.Spacer(1),

				// Inputs with focus/blur callbacks.
				// InputWithOptions extends Input with optional focus/blur hooks.
				react.Bold("Inputs:"),
				react.Column(
					react.InputWithOptions(inputValue, func(v string) { setInputValue(v) },
						"Input 1 placeholder", 20,
						react.InputOnFocus(onFocus("Input 1")),
						react.InputOnBlur(onBlur("Input 1")),
					),
					react.Spacer(0),
					react.InputWithOptions("", func(v string) {},
						"Input 2 placeholder", 20,
						react.InputOnFocus(onFocus("Input 2")),
						react.InputOnBlur(onBlur("Input 2")),
					),
				),

				react.Spacer(1),

				// Checkbox (plain — demonstrates mixed interactive types).
				react.Bold("Checkbox (plain):"),
				react.Checkbox("Subscribe to newsletter", false, func(bool) {}),

				react.Spacer(1),
				react.Divider(react.DividerLabel(" Status ")),

				// Status display.
				react.Textf("Focus: %s", focusStatus),
				react.Textf("Click: %s", clickLog),

				react.Spacer(1),
				react.Text("Tab/arrows • Enter/Space to click • Observe focus/blur on status line"),
			),
			react.WithBorder(),
			react.WithPadding(1),
			react.WithWidth(60),
			react.WithTitle("Navigation Demo"),
		),
	)
}

func main() {
	p := react.Root(&NavigationDemo{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
