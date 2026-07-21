// Package react provides a component-based, reactive abstraction over
// Bubble Tea. Instead of the classic Elm Architecture (Model/Msg/Update),
// you compose declarative components with hooks (useState, useEffect).
//
// Basic usage:
//
//	type Counter struct{}
//	func (c *Counter) Render(ctx *react.Context) react.Element {
//	    count, setCount := ctx.UseState(0)
//	    return react.Column(
//	        react.Textf("Count: %d", count),
//	        react.Button("+", func() { setCount(count + 1) }),
//	        react.Button("-", func() { setCount(count - 1) }),
//	    )
//	}
//
//	func main() {
//	    p := react.Root(&Counter{})
//	    p.Run()
//	}
package react

// Element is a virtual node in the component tree.
// Every concrete element type must implement this marker interface.
type Element interface {
	element()
}

// Component is the interface all user-defined components must implement.
// Render returns an Element tree describing what to draw.
type Component interface {
	Render(ctx *Context) Element
}

// Context is injected into every Component.Render call.
// It provides hooks for state management and lifecycle.
// Do not create Context values yourself; they are provided by the framework.
//
// Use the package-level hook functions (UseState, UseEffect, UseRef)
// with the Context rather than calling methods on it.
type Context struct {
	componentID uint64
	hookIndex   int
	hooks       []any
}

// ResetHooks resets the hook index for a new render cycle.
// Called internally by the bridge before each render.
func (ctx *Context) ResetHooks() {
	ctx.hookIndex = 0
}

// Hooks returns the hook list for this context.
// Used internally by the reconciler.
func (ctx *Context) Hooks() []any {
	return ctx.hooks
}

// --- Concrete Element Types ---

// TextElement renders a plain string.
type TextElement struct {
	Text string
}

func (TextElement) element() {}

// BoldElement renders bold text.
type BoldElement struct {
	Text string
}

func (BoldElement) element() {}

// ComponentElement wraps a child component for nesting.
type ComponentElement struct {
	Component Component
	Key       string // for reconciliation identity
}

func (ComponentElement) element() {}

// ButtonElement renders a clickable button.
type ButtonElement struct {
	Label   string
	OnClick func()
}

func (ButtonElement) element() {}

// ColumnElement lays out children vertically.
type ColumnElement struct {
	Children []Element
}

func (ColumnElement) element() {}

// RowElement lays out children horizontally.
type RowElement struct {
	Children []Element
}

func (RowElement) element() {}

// BoxElement is a styled container with border and padding.
type BoxElement struct {
	Child     Element
	Border    bool
	Padding   int
	Width     int
	Height    int
	Title     string
}

func (BoxElement) element() {}

// SpacerElement adds vertical space.
type SpacerElement struct {
	Height int
}

func (SpacerElement) element() {}

// InputElement renders a text input field.
type InputElement struct {
	Value     string
	OnChange  func(string)
	Placeholder string
	Width     int
}

func (InputElement) element() {}
