package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	react "github.com/DaviMGDev/bubbletea-react/react"
)

// ChatDemo demonstrates RowOpts WithCollapse, Input with InputOnSubmit,
// Button, and list rendering.
//
// The Input uses InputOnSubmit to capture Enter directly, so you can type
// a message and press Enter to send without tabbing to the Send button.
type ChatDemo struct{}

func (d *ChatDemo) Render(ctx *react.Context) react.Element {
	messages, setMessages := react.UseState(ctx, []string{})
	input, setInput := react.UseState(ctx, "")

	submitMessage := func() {
		trimmed := strings.TrimSpace(input)
		if trimmed == "" {
			return
		}
		ts := time.Now().Format("15:04")
		msg := fmt.Sprintf("[%s] %s", ts, trimmed)
		newMsgs := append(messages, msg)
		// Keep at most 50 messages to avoid unbounded growth.
		if len(newMsgs) > 50 {
			newMsgs = newMsgs[len(newMsgs)-50:]
		}
		setMessages(newMsgs)
		setInput("")
	}

	// Build message list elements (newest at bottom).
	msgEls := make([]react.Element, 0, len(messages))
	if len(messages) == 0 {
		msgEls = append(msgEls, react.Text("  No messages yet. Type below and press Send!"))
	} else {
		for _, m := range messages {
			msgEls = append(msgEls, react.Text("  │ "+m))
		}
	}

	return react.View(
		react.Box(
			react.Column(
				react.Bold("💬 Chat UI"),
				react.Text("Features: RowOpts WithCollapse, Input, Button, stateful list."),
				react.Spacer(1),

				// Message area.
				react.Divider(react.DividerLabel(" Messages ")),
				react.Column(msgEls...),

				react.Spacer(1),
				react.Divider(),

				// Input row: when terminal is narrow (<55 cols), the row
				// collapses into a column — Input stacks above Send button.
				react.RowOpts(
					react.WithCollapse(true),
					react.WithCollapseAt(55),
				).Children(
					react.InputWithOptions(input, func(v string) { setInput(v) }, "Type a message...", 30,
					react.InputOnSubmit(func() { submitMessage() }),
				),
					react.Text("  "),
					react.Button("Send", func() { submitMessage() }),
				),

				react.Spacer(1),
				react.Textf("%d message(s)", len(messages)),
				react.Text("Type • Enter to submit • Tab/arrows navigate • Resize terminal for collapse"),
			),
			react.WithBorder(),
			react.WithPadding(1),
			react.WithWidth(65),
			react.WithTitle("Chat Demo"),
		),
	)
}

func main() {
	p := react.Root(&ChatDemo{})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
