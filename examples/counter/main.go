package main

import (
	"fmt"
	"os"

	react "github.com/DaviMGDev/bubbletea-react/react"
)

// Counter is a simple counter component with increment/decrement/reset.
type Counter struct{}

func (c *Counter) Render(ctx *react.Context) react.Element {
	count, setCount := react.UseState(ctx, 0)

	return react.View(
		react.Box(
			react.Column(
				react.Bold("Counter Demo"),
				react.Spacer(1),
				react.Textf("Count: %d", count),
				react.Spacer(1),
				react.Row(
					react.Button("+1", func() { setCount(count + 1) }),
					react.Text("  "),
					react.Button("-1", func() { setCount(count - 1) }),
					react.Text("  "),
					react.Button("Reset", func() { setCount(0) }),
				),
				react.Spacer(1),
				react.Text("Tab/arrows to navigate, Enter to click."),
			),
			react.WithBorder(),
			react.WithPadding(1),
			react.WithWidth(40),
			react.WithTitle("Counter"),
		),
	)
}

func main() {
	p := react.Root(&Counter{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
