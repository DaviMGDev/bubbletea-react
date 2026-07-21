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
- **Declarative elements**: `Text`, `Button`, `Column`, `Row`, `Box`, `Input`, `Spacer`, `Checkbox`, `Select`, `Divider`, `Progress`, `Tabs`, `Table`, `Scrollable`, `Form`
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

## Examples

All examples are in the [`examples/`](examples/) directory. Run any of them with:

```bash
cd examples/<name>
go run .
```

| Example | Features Demonstrated |
|---|---|
| [`counter`](examples/counter/) | `UseState`, `Button`, `Row`, `Box`, `Bold`, `Spacer`, `Textf` |
| [`forms`](examples/forms/) | `Input`, `Checkbox`, `Select`, `Progress`, `Tabs`, `Divider`, multiple `UseState` |
| [`todo`](examples/todo/) | `UseState` with dynamic lists, `Input`, `Button`, list rendering |
| [`effects`](examples/effects/) | `UseEffect` (dependency tracking & cleanup), `UseRef`, `Scrollable`, `Col` with `WithAlign` |
| [`datatable`](examples/datatable/) | `Table`, `Tabs`, `Wrap` (nested components with isolated state) |
| [`chat`](examples/chat/) | `RowOpts` `WithCollapse` (responsive layout), `Input`, `Button`, list state |
| [`navigation`](examples/navigation/) | `ButtonWithOptions`, `InputWithOptions` (`OnFocus`/`OnBlur` callbacks) |
| [`stopwatch`](examples/stopwatch/) | `UseEffect`, `UseRef` (lap storage), `Progress`, `Button` control groups |
| [`explorer`](examples/explorer/) | Deep `Wrap` nesting tree, `Scrollable`, `Table` detail view, per-component state |
| [`dashboard`](examples/dashboard/) | `Col`/`RowOpts` layout (`WithWrap`, `WithAlign`), `Form`, `Box` variants, `TableNoHeader`, inline `DividerStyle` |

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
| `ButtonWithOptions(label, onClick, opts...)` | Button with focus/blur callbacks |
| `Input(value, onChange, placeholder, width)` | Text input |
| `InputWithOptions(value, onChange, placeholder, width, opts...)` | Input with focus/blur callbacks |
| `Checkbox(label, checked, onToggle)` | Toggleable checkbox |
| `Select(options, selected, onChange)` | Selection from options list |
| `Column(children...)` | Vertical layout |
| `Row(children...)` | Horizontal layout |
| `Col(opts...).Children(children...)` | Column with layout options (align, scrollable) |
| `RowOpts(opts...).Children(children...)` | Row with layout options (align, wrap, collapse) |
| `Box(child, opts...)` | Styled container (border, padding, width, height, title) |
| `Scrollable(child, height)` | Scrollable viewport |
| `Form(onSubmit, children...)` | Form container |
| `Spacer(height)` | Vertical spacing |
| `Divider(opts...)` | Horizontal separator |
| `Progress(current, total, opts...)` | Progress bar |
| `Tabs(tabs, active, onChange)` | Tab container with content panes |
| `Table(columns, rows, opts...)` | Data table |
| `Wrap(component, key)` | Nest a child component |
| `View(children...)` | Top-level wrapper (alias for Column) |

### Layout Options (Col / RowOpts)

| Function | Description |
|---|---|
| `WithAlign(v, h)` | Set vertical + horizontal alignment |
| `WithWrap(bool)` | Enable row wrapping on overflow |
| `WithCollapse(bool)` | Auto-collapse row to column on narrow terminals |
| `WithCollapseAt(threshold)` | Set collapse width threshold (default 80) |
| `WithScrollable(height)` | Make column scrollable with fixed viewport |

### Box Options

| Function | Description |
|---|---|
| `WithBorder()` | Rounded border |
| `WithPadding(n)` | Internal padding |
| `WithWidth(n)` | Width in characters |
| `WithHeight(n)` | Height in characters |
| `WithTitle(s)` | Header title |

### Divider Options

| Function | Description |
|---|---|
| `DividerStyle("thin" / "double" / "dashed")` | Line style |
| `DividerWidth(n)` | Line width |
| `DividerLabel(s)` | Centered label |

### Progress Options

| Function | Description |
|---|---|
| `ProgressWidth(n)` | Bar width in characters |
| `ProgressLabel(s)` | Additional label text |

### Button / Input Options

| Function | Applies to | Description |
|---|---|---|
| `WithOnFocus(fn)` | Button | Fires when element receives Tab focus |
| `WithOnBlur(fn)` | Button | Fires when element loses Tab focus |
| `InputOnFocus(fn)` | Input | Fires when input receives focus |
| `InputOnBlur(fn)` | Input | Fires when input loses focus |
| `InputOnSubmit(fn)` | Input | Fires when Enter is pressed while input is focused (form submission) |

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
