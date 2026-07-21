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

// --- Alignment ---

// AlignPosition describes a position along one axis.
type AlignPosition int

const (
	Top     AlignPosition = iota // top / left
	Bottom                       // bottom / right
	Center                       // centered
	Stretch                      // fill available space
)

// Align is a 2D alignment pair (vertical, horizontal).
type Align struct {
	V, H AlignPosition
}

// Weight describes how a child fills remaining space in a row/column.
type Weight int

const (
	WeightNone       Weight = iota // natural size
	WeightHorizontal               // fill remaining horizontal space
	WeightVertical                 // fill remaining vertical space
	WeightBoth                     // fill both axes
)

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
	OnFocus func()
	OnBlur  func()
}

func (ButtonElement) element() {}

// ColumnElement lays out children vertically.
type ColumnElement struct {
	Children []Element
	Align    *Align   // nil = default (top, left)
	Wrap     bool     // unused on column
	Collapse bool     // unused on column
	Scrollable bool   // constrain height and scroll
	Height   int      // viewport height when scrollable
}

func (ColumnElement) element() {}

// RowElement lays out children horizontally.
type RowElement struct {
	Children []Element
	Align    *Align   // nil = default (top, left)
	Wrap     bool     // auto-wrap to next line on overflow
	Collapse bool     // render as column when width < CollapseAt
	CollapseAt int    // width threshold for collapse (default 80)
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
	Value       string
	OnChange    func(string)
	Placeholder string
	Width       int
	OnFocus     func()
	OnBlur      func()
}

func (InputElement) element() {}

// ScrollableElement constrains its child to a fixed height with scrolling.
type ScrollableElement struct {
	Child  Element
	Height int // viewport height in lines
}

func (ScrollableElement) element() {}

// FormElement groups children and handles submission.
type FormElement struct {
	Children []Element
	OnSubmit func()
}

func (FormElement) element() {}

// CheckboxElement renders a toggleable checkbox.
type CheckboxElement struct {
	Label    string
	Checked  bool
	OnChange func(bool)
}

func (CheckboxElement) element() {}

// SelectOption is a single option in a Select.
type SelectOption struct {
	Label string
	Value string
}

// SelectElement renders a selection list.
type SelectElement struct {
	Options  []SelectOption
	Selected int // index of selected option (-1 = none)
	OnChange func(string)
}

func (SelectElement) element() {}

// DividerElement renders a horizontal separator.
type DividerElement struct {
	Style  string // "thin", "double", "dashed" (default "thin")
	Width  int    // 0 = full terminal width
	Label  string // optional centered label
}

func (DividerElement) element() {}

// ProgressElement renders a progress bar.
type ProgressElement struct {
	Current int
	Total   int
	Width   int // bar width in chars (0 = full terminal width)
	Label   string // optional label suffix
}

func (ProgressElement) element() {}

// Tab defines a single tab pane.
type Tab struct {
	Label string
	Content Element
}

// TabsElement renders tab navigation with content panes.
type TabsElement struct {
	Tabs     []Tab
	Active   int // index of active tab
	OnChange func(int)
}

func (TabsElement) element() {}

// TableColumn defines a column in a table.
type TableColumn struct {
	Label string
	Width int // 0 = auto
}

// TableElement renders tabular data.
type TableElement struct {
	Columns []TableColumn
	Rows    [][]string
	Header  bool // whether to render a header row
}

func (TableElement) element() {}
