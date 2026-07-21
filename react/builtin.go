package react

import "fmt"

// --- Builder Functions ---

// Text returns a plain text element.
func Text(text string) Element {
	return TextElement{Text: text}
}

// Textf returns a formatted text element (fmt.Sprintf style).
func Textf(format string, args ...any) Element {
	return TextElement{Text: fmt.Sprintf(format, args...)}
}

// Bold returns a bold text element.
func Bold(text string) Element {
	return BoldElement{Text: text}
}

// Button returns a clickable button element.
func Button(label string, onClick func()) Element {
	return ButtonElement{Label: label, OnClick: onClick}
}

// Column lays out children vertically (top to bottom).
func Column(children ...Element) Element {
	return ColumnElement{Children: children}
}

// Row lays out children horizontally (left to right).
func Row(children ...Element) Element {
	return RowElement{Children: children}
}

// Box wraps a child element in a styled container.
// Width, Height, Padding, and Border control the appearance.
func Box(child Element, opts ...BoxOption) Element {
	b := BoxElement{Child: child}
	for _, opt := range opts {
		opt(&b)
	}
	return b
}

// BoxOption is a functional option for Box configuration.
type BoxOption func(*BoxElement)

// WithBorder enables a rounded border on the box.
func WithBorder() BoxOption {
	return func(b *BoxElement) {
		b.Border = true
	}
}

// WithPadding sets the padding inside the box.
func WithPadding(padding int) BoxOption {
	return func(b *BoxElement) {
		b.Padding = padding
	}
}

// WithWidth sets the box width in characters.
func WithWidth(width int) BoxOption {
	return func(b *BoxElement) {
		b.Width = width
	}
}

// WithHeight sets the box height in characters.
func WithHeight(height int) BoxOption {
	return func(b *BoxElement) {
		b.Height = height
	}
}

// WithTitle sets a title displayed at the top of the box.
func WithTitle(title string) BoxOption {
	return func(b *BoxElement) {
		b.Title = title
	}
}

// Spacer adds vertical spacing of the given height.
func Spacer(height int) Element {
	return SpacerElement{Height: height}
}

// Input returns a text input element.
func Input(value string, onChange func(string), placeholder string, width int) Element {
	return InputElement{
		Value:       value,
		OnChange:    onChange,
		Placeholder: placeholder,
		Width:       width,
	}
}

// View is a convenience wrapper that groups children into a Column.
// It can be used as the top-level return of a component's Render method.
func View(children ...Element) Element {
	return Column(children...)
}

// Wrap wraps a child Component as an Element, enabling nesting.
// The optional key helps the reconciler identify this component across renders.
func Wrap(c Component, key string) Element {
	return ComponentElement{Component: c, Key: key}
}
