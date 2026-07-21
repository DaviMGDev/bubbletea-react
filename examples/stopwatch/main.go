package main

import (
	"fmt"
	"os"
	"time"

	react "github.com/DaviMGDev/bubbletea-react/react"
)

// StopwatchDemo demonstrates UseEffect (logging state transitions),
// UseRef (storing lap times), Progress bar, and button groups.
//
// Note: The timer advances via button presses (1s per "Tick").
// Real auto-timing would require Bubble Tea's tea.Tick command,
// which this library doesn't expose from components. This design
// still cleanly demonstrates UseEffect + UseRef + Progress together.
type StopwatchDemo struct{}

func (d *StopwatchDemo) Render(ctx *react.Context) react.Element {
	elapsed, setElapsed := react.UseState(ctx, 0)
	running, setRunning := react.UseState(ctx, false)

	// UseRef stores lap times across renders without triggering re-renders.
	lapsRef := react.UseRef(ctx, []string{})

	// UseEffect: log transitions when running state changes.
	// The effect closure captures `running` and `elapsed` at render time.
	// A cleanup function could cancel external resources (not needed here).
	react.UseEffect(ctx, func() func() {
		if running {
			entry := fmt.Sprintf("[START] %s — timer running", time.Now().Format("15:04:05"))
			lapsRef.Value = append(lapsRef.Value, entry)
		} else if elapsed > 0 {
			entry := fmt.Sprintf("[STOP]  %s — elapsed %ds", time.Now().Format("15:04:05"), elapsed)
			lapsRef.Value = append(lapsRef.Value, entry)
		}
		return nil // no cleanup
	}, []any{running})

	// UseEffect with empty deps: runs once on mount.
	react.UseEffect(ctx, func() func() {
		lapsRef.Value = append(lapsRef.Value, "[MOUNT] Stopwatch ready — press Start")
		return nil
	}, []any{})

	// Build lap entries for display (newest first).
	lapEntries := make([]react.Element, 0, len(lapsRef.Value))
	for i := len(lapsRef.Value) - 1; i >= 0; i-- {
		lapEntries = append(lapEntries, react.Text("  "+lapsRef.Value[i]))
	}
	if len(lapEntries) == 0 {
		lapEntries = append(lapEntries, react.Text("  (no laps yet)"))
	}

	// Clamp elapsed for display.
	displayElapsed := elapsed
	if displayElapsed < 0 {
		displayElapsed = 0
	}

	return react.View(
		react.Box(
			react.Column(
				react.Bold("⏱ Stopwatch Demo"),
				react.Text("Features: UseEffect, UseRef (laps), Progress bar, button controls."),
				react.Spacer(1),
				react.Divider(react.DividerLabel(" Timer ")),

				// Main timer display.
				react.Col(react.WithAlign(react.Center, react.Center)).Children(
					react.Textf("Elapsed: %ds", displayElapsed),
				),

				// Progress bar: current/max with custom width.
				react.Progress(displayElapsed, 30,
					react.ProgressWidth(30),
					react.ProgressLabel(fmt.Sprintf("%ds/30s", displayElapsed)),
				),

				react.Spacer(1),

				// Control buttons.
				react.Row(
					react.Button("Start", func() {
						if !running {
							setRunning(true)
						}
					}),
					react.Text(" "),
					react.Button("Stop", func() {
						if running {
							setRunning(false)
						}
					}),
					react.Text(" "),
					react.Button("Tick ⟳", func() {
						setElapsed(elapsed + 1)
					}),
					react.Text(" "),
					react.Button("Lap", func() {
						if running || elapsed > 0 {
							entry := fmt.Sprintf("[LAP]   %s — %ds", time.Now().Format("15:04:05"), elapsed)
							lapsRef.Value = append(lapsRef.Value, entry)
						}
					}),
					react.Text(" "),
					react.Button("Reset", func() {
						setElapsed(0)
						setRunning(false)
						lapsRef.Value = []string{}
					}),
				),

				react.Spacer(1),
				react.Divider(react.DividerLabel(" Laps & Events ")),
				react.Column(lapEntries...),

				react.Spacer(1),
				react.Text("Start/Stop toggles • Tick advances 1s • Lap records current time • Reset clears"),
			),
			react.WithBorder(),
			react.WithPadding(1),
			react.WithWidth(58),
			react.WithTitle("Stopwatch Demo"),
		),
	)
}

func main() {
	p := react.Root(&StopwatchDemo{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
