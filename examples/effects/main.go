package main

import (
	"fmt"
	"os"
	"time"

	react "github.com/DaviMGDev/bubbletea-react/react"
)

// EffectsDemo demonstrates UseEffect, UseRef, Scrollable, and layout options.
type EffectsDemo struct{}

func (d *EffectsDemo) Render(ctx *react.Context) react.Element {
	// --- Hooks ---
	count, setCount := react.UseState(ctx, 0)

	// UseRef: mutable buffer that persists across renders.
	// Unlike UseState, mutating a Ref does NOT trigger a re-render —
	// it's useful for data that changes frequently but doesn't need
	// to drive the UI layout directly.
	logRef := react.UseRef(ctx, []string{})

	// UseEffect with dependency [count]: runs AFTER every render where
	// `count` has changed. The closure captures the current `count`
	// value. The optional return func is cleanup (called before the
	// next effect run or on unmount).
	react.UseEffect(ctx, func() func() {
		ts := time.Now().Format("15:04:05.000")
		entry := fmt.Sprintf("[%s] Count → %d", ts, count)
		logRef.Value = append(logRef.Value, entry)
		// Keep at most 50 entries to avoid unbounded growth.
		if len(logRef.Value) > 50 {
			logRef.Value = logRef.Value[len(logRef.Value)-50:]
		}
		return nil // no cleanup needed for this effect
	}, []any{count})

	// UseEffect with empty deps []: runs exactly once on mount.
	// A non-nil cleanup would run only on unmount.
	react.UseEffect(ctx, func() func() {
		logRef.Value = append(logRef.Value, "[LIFECYCLE] Component mounted — press +/- to change count")
		return nil
	}, []any{})

	// UseEffect with nil deps: runs on EVERY render (no dep tracking).
	// Use sparingly — demonstrated here for completeness.
	react.UseEffect(ctx, func() func() {
		return nil
	}, nil)

	// Build log entry elements (newest entries first, so they appear
	// at the top of the scrollable viewport).
	logEntries := make([]react.Element, 0, len(logRef.Value))
	for i := len(logRef.Value) - 1; i >= 0; i-- {
		logEntries = append(logEntries, react.Text(logRef.Value[i]))
	}
	if len(logEntries) == 0 {
		logEntries = append(logEntries, react.Text("(no entries yet — press a button)"))
	}

	// --- Element Tree ---
	return react.View(
		react.Box(
			react.Column(
				react.Bold("🎯 Effects & Refs Demo"),
				react.Spacer(1),
				react.Text("Features: UseEffect, UseRef, Scrollable, layout alignment"),

				react.Divider(react.DividerLabel(" Counter ")),

				// Col with alignment: centers the count text within its cell.
				react.Col(react.WithAlign(react.Center, react.Center)).Children(
					react.Textf("Count: %d", count),
				),

				react.Spacer(1),
				react.Row(
					react.Button("+1", func() { setCount(count + 1) }),
					react.Text("  "),
					react.Button("-1", func() { setCount(count - 1) }),
					react.Text("  "),
					react.Button("Reset", func() { setCount(0) }),
				),

				react.Spacer(1),
				react.Row(
					react.Button("Clear Log", func() {
						logRef.Value = []string{}
					}),
				),

				react.Spacer(1),
				react.Divider(react.DividerLabel(" Log (scrollable) ")),

				// Scrollable constrains the child to a fixed viewport height.
				// Content exceeding Height lines is clipped.
				react.Scrollable(
					react.Column(logEntries...),
					8, // viewport height in lines
				),

				react.Spacer(1),
				react.Text("Tab/arrows to navigate • Enter to click • q to quit"),
			),
			react.WithBorder(),
			react.WithPadding(1),
			react.WithWidth(58),
			react.WithTitle("Effects Demo"),
		),
	)
}

func main() {
	p := react.Root(&EffectsDemo{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
