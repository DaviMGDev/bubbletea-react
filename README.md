# bubbletea-react

A **component-based, reactive** abstraction over [Bubble Tea](https://github.com/charmbracelet/bubbletea).  
Instead of the Elm Architecture (Model/Msg/Update), compose declarative components with hooks — familiar from React.

```go
type Counter struct{}
func (c *Counter) Render(ctx *react.Context) react.Element {
    count, setCount := react.UseState(ctx, 0)
    return react.Column(
        react.Textf("Count: %d", count),
        react.Button("+1", func() { setCount(count + 1) }),
        react.Button("-1", func() { setCount(count - 1) }),
    )
}

func main() { react.Root(&Counter{}).Run() }
```

## Features

- **Hooks**: `UseState`, `UseEffect`, `UseRef` — deterministic, order-based (like React)
- **Declarative elements**: `Text`, `Button`, `Column`, `Row`, `Box`, `Input`, `Spacer`
- **Composition**: Nest components via the `Wrap` function
- **Keyboard navigation**: Tab/arrows to focus, Enter to activate
- **Zero config**: `Root(component)` returns a ready-to-run `*tea.Program`

## Installation

```bash
go get github.com/DaviMGDev/bubbletea-react
```

Then import:

```go
import react "github.com/DaviMGDev/bubbletea-react/react"
```

## Quick Start

### Counter

```go
package main

import (
    react "github.com/DaviMGDev/bubbletea-react/react"
)

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
            ),
            react.WithBorder(),
            react.WithPadding(1),
            react.WithWidth(40),
            react.WithTitle("Counter"),
        ),
    )
}

func main() {
    react.Root(&Counter{}).Run()
}
```

## API

### Hooks

| Function | Description |
|---|---|
| `UseState[T](ctx, initial) (T, func(T))` | Stateful value with setter |
| `UseEffect(ctx, func() func(), []any)` | Side effects with dependency tracking |
| `UseRef[T](ctx, initial) *Ref[T]` | Mutable reference persisting across renders |

### Elements

| Function | Description |
|---|---|
| `Text(s)` | Plain text |
| `Textf(format, args...)` | Formatted text |
| `Bold(s)` | Bold text |
| `Button(label, onClick)` | Clickable button |
| `Column(children...)` | Vertical layout |
| `Row(children...)` | Horizontal layout |
| `Box(child, opts...)` | Styled container (border, padding, width, height, title) |
| `Spacer(height)` | Vertical spacing |
| `Input(value, onChange, placeholder, width)` | Text input |
| `Wrap(component, key)` | Nest a child component |
| `View(children...)` | Top-level wrapper (alias for Column) |

### Box Options

| Function | Description |
|---|---|
| `WithBorder()` | Rounded border |
| `WithPadding(n)` | Internal padding |
| `WithWidth(n)` | Width in characters |
| `WithHeight(n)` | Height in characters |
| `WithTitle(s)` | Header title |

## How It Works

1. Your component's `Render` method is called to produce an **element tree** (virtual nodes)
2. The **reconciler** serializes the tree to a string (Bubble Tea's `View`)
3. **Hooks** store state in a per-component context, keyed by call order
4. When a setter is called, the hook value is updated in-place and the bridge re-renders the tree
5. **Keyboard events** (Tab, arrows, Enter, Space) are routed to the focused interactive element

## Limitations (v1)

- Full tree re-render on every state change (no per-component bailout)
- Nested component state restored via position-based keys (explicit keys recommended for dynamic lists)
- Mouse click not yet supported (keyboard-only navigation)
- Window resize handling is basic (uses stored terminal width)

## License

MIT
