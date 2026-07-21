package main

import (
	"fmt"
	"os"
	"time"

	react "github.com/DaviMGDev/bubbletea-react/react"
)

// DashboardDemo — compact 5-widget grid designed for 158×41 terminals.
// All widgets are height-balanced within each row to avoid ugly gaps.
type DashboardDemo struct{}

func (d *DashboardDemo) Render(ctx *react.Context) react.Element {
	now := time.Now()
	uptime := int(now.Sub(time.Date(2025, 7, 21, 0, 0, 0, 0, time.UTC)).Hours())

	// ── State ─────────────────────────────────────────────────────────
	cpu, setCPU := react.UseState(ctx, 67)
	activeTab, setActiveTab := react.UseState(ctx, 1)
	formInput, setFormInput := react.UseState(ctx, "")
	formSubmitted, setFormSubmitted := react.UseState(ctx, "(none)")
	formError, setFormError := react.UseState(ctx, false)
	feedback, setFeedback := react.UseState(ctx, "Ready")

	adjustCPU := func(d int) {
		v := cpu + d
		if v < 0 {
			v = 0
		}
		if v > 100 {
			v = 100
		}
		setCPU(v)
		setFeedback(fmt.Sprintf("CPU → %d%%", v))
	}

	submitForm := func() {
		if formInput == "" {
			setFormError(true)
			setFeedback("⚠ Empty input")
			return
		}
		setFormSubmitted(formInput)
		setFeedback(fmt.Sprintf("✓ Submitted: %s", formInput))
		setFormInput("")
		setFormError(false)
	}

	// ── Row 1 widgets (height-balanced at ~9 lines each) ──────────────
	sysInfo := react.Box(
		react.Column(
			react.Bold("System"),
			react.Textf("Host:    bubbletea-react"),
			react.Textf("Uptime:  %dh", uptime),
			react.Textf("Version: v1.0.0"),
			react.Textf("License: MIT"),
		),
		react.WithBorder(), react.WithWidth(36), react.WithPadding(1),
	)

	resPanel := react.Box(
		react.Column(
			react.Bold("Resources"),
			react.Progress(cpu, 100, react.ProgressWidth(28), react.ProgressLabel("CPU")),
			react.Progress(42, 100, react.ProgressWidth(28), react.ProgressLabel("MEM")),
			react.Progress(88, 100, react.ProgressWidth(28), react.ProgressLabel("DSK")),
			react.Row(
				react.Button("-5", func() { adjustCPU(-5) }),
				react.Text(" "),
				react.Button("+5", func() { adjustCPU(+5) }),
			),
		),
		react.WithBorder(), react.WithWidth(44), react.WithPadding(1),
	)

	tabsPanel := react.Box(
		react.Column(
			react.Bold("Tabs"),
			react.Tabs([]react.Tab{
				{
					Label: "Overview",
					Content: react.Column(
						react.Text("  All systems nominal."),
						react.Text("  0 alerts, 0 warnings."),
					),
				},
				{
					Label: "Logs",
					Content: react.Column(
						react.Text("  [INFO]  Server started"),
						react.Text("  [INFO]  Connection OK"),
						react.Text("  [WARN]  Disk at 88%"),
					),
				},
				{
					Label: "Config",
					Content: react.Column(
						react.Text("  Theme:  dark"),
						react.Text("  Port:   8080"),
						react.Text("  Debug:  off"),
					),
				},
			}, activeTab, func(i int) { setActiveTab(i) }),
		),
		react.WithBorder(), react.WithWidth(46), react.WithPadding(1),
	)

	// ── Row 2 widgets (height-balanced at ~9 lines each) ──────────────
	svcPanel := react.Box(
		react.Column(
			react.Bold("Services"),
			react.Table(
				[]react.TableColumn{
					{Label: "Service", Width: 14},
					{Label: "Status", Width: 14},
					{Label: "Uptime", Width: 8},
				},
				[][]string{
					{"web-server", "🟢 online", "14d"},
					{"database", "🟢 online", "14d"},
					{"cache", "🟡 degraded", "7d"},
					{"worker", "🔴 offline", "0d"},
				},
				react.TableNoHeader(),
			),
		),
		react.WithBorder(), react.WithWidth(64), react.WithPadding(1),
	)

	formField := func() react.Element {
		if formError {
			return react.Column(
				react.InputWithOptions(formInput, func(v string) { setFormInput(v) },
					"Type a message...", 30,
					react.InputOnSubmit(func() { submitForm() }),
				),
				react.Text("⚠ Required — cannot be empty"),
			)
		}
		return react.InputWithOptions(formInput, func(v string) { setFormInput(v) },
			"Type a message...", 30,
			react.InputOnSubmit(func() { submitForm() }),
		)
	}

	formPanel := react.Box(
		react.Column(
			react.Bold("Form"),
			react.Text("The Form element wraps children and"),
			react.Text("appends a Submit button."),
			react.Spacer(0),
			react.Form(func() { submitForm() }, formField()),
			react.Spacer(0),
			react.Textf("Submitted: %s", func() string {
				if formSubmitted != "" {
					return formSubmitted
				}
				return "(none)"
			}()),
		),
		react.WithBorder(), react.WithWidth(64), react.WithPadding(1),
	)

	// ── Assemble ──────────────────────────────────────────────────────
	return react.View(
		react.Box(
			react.Column(
				react.Bold("📊 Dashboard Demo"),
				react.Text("Box • Progress • Tabs • TableNoHeader • Form • DividerStyle"),
				react.Spacer(0),

				react.Row(
					sysInfo, react.Text(" "), resPanel, react.Text(" "), tabsPanel,
				),
				react.Row(
					svcPanel, react.Text(" "), formPanel,
				),

				react.Spacer(0),
				react.Divider(react.DividerStyle("double")),
				react.Row(
					react.Divider(react.DividerWidth(12)),
					react.Text(" "),
					react.Divider(react.DividerStyle("double"), react.DividerWidth(12)),
					react.Text(" "),
					react.Divider(react.DividerStyle("dashed"), react.DividerWidth(12)),
					react.Text(" "),
					react.Divider(react.DividerStyle("thin"), react.DividerWidth(12), react.DividerLabel("x")),
				),
				react.Spacer(0),
				react.Text(feedback),
				react.Text("Tab Navigate • Enter Click • -5/+5 CPU • Enter submits"),
			),
			react.WithBorder(), react.WithPadding(1),
			react.WithTitle("Dashboard"),
		),
	)
}

func main() {
	p := react.Root(&DashboardDemo{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
