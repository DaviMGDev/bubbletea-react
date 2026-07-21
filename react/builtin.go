package react

import "fmt"

// --- Text Elements ---

func Text(text string) Element {
	return TextElement{Text: text}
}

func Textf(format string, args ...any) Element {
	return TextElement{Text: fmt.Sprintf(format, args...)}
}

func Bold(text string) Element {
	return BoldElement{Text: text}
}

// --- Interactive Elements ---

func Button(label string, onClick func()) Element {
	return ButtonElement{Label: label, OnClick: onClick}
}

// ButtonWithOptions returns a button with optional focus/blur callbacks.
func ButtonWithOptions(label string, onClick func(), opts ...ButtonOption) Element {
	b := ButtonElement{Label: label, OnClick: onClick}
	applyButtonOpts(&b, opts)
	return b
}

type ButtonOption func(*ButtonElement)

func WithOnFocus(fn func()) ButtonOption {
	return func(b *ButtonElement) { b.OnFocus = fn }
}

func WithOnBlur(fn func()) ButtonOption {
	return func(b *ButtonElement) { b.OnBlur = fn }
}

func applyButtonOpts(b *ButtonElement, opts []ButtonOption) {
	for _, o := range opts {
		o(b)
	}
}

func Input(value string, onChange func(string), placeholder string, width int) Element {
	return InputElement{
		Value:       value,
		OnChange:    onChange,
		Placeholder: placeholder,
		Width:       width,
	}
}

// InputWithOptions returns an input with optional focus/blur callbacks.
func InputWithOptions(value string, onChange func(string), placeholder string, width int, opts ...InputOption) Element {
	i := InputElement{
		Value:       value,
		OnChange:    onChange,
		Placeholder: placeholder,
		Width:       width,
	}
	applyInputOpts(&i, opts)
	return i
}

type InputOption func(*InputElement)

func InputOnFocus(fn func()) InputOption {
	return func(i *InputElement) { i.OnFocus = fn }
}

func InputOnBlur(fn func()) InputOption {
	return func(i *InputElement) { i.OnBlur = fn }
}

func applyInputOpts(i *InputElement, opts []InputOption) {
	for _, o := range opts {
		o(i)
	}
}

func Checkbox(label string, checked bool, onToggle func(bool)) Element {
	return CheckboxElement{Label: label, Checked: checked, OnChange: onToggle}
}

func Select(options []SelectOption, selected int, onChange func(string)) Element {
	opts := make([]SelectOption, len(options))
	copy(opts, options)
	return SelectElement{Options: opts, Selected: selected, OnChange: onChange}
}

// --- Layout Elements ---

func Column(children ...Element) Element {
	return ColumnElement{Children: children}
}

// Col returns a Column with layout options.
// Usage: Col(WithAlign(Center, Left)).Children(Text("a"), Text("b"))
func Col(opts ...LayoutOption) ColumnBuilder {
	return ColumnBuilder{opts: opts}
}

type ColumnBuilder struct {
	opts []LayoutOption
}

func (b ColumnBuilder) Children(children ...Element) Element {
	var align Align
	var hasAlign bool
	var wrap bool
	var collapse bool
	var collapseAt int
	var scrollable bool
	var height int
	for _, o := range b.opts {
		o(&align, &hasAlign, &wrap, &collapse, &collapseAt, &scrollable, &height)
	}
	var ap *Align
	if hasAlign {
		ap = &align
	}
	return ColumnElement{
		Children:   children,
		Align:      ap,
		Scrollable: scrollable,
		Height:     height,
	}
}

func Row(children ...Element) Element {
	return RowElement{Children: children}
}

// RowOpts returns a Row with layout options.
func RowOpts(opts ...LayoutOption) RowBuilder {
	return RowBuilder{opts: opts}
}

type RowBuilder struct {
	opts []LayoutOption
}

func (b RowBuilder) Children(children ...Element) Element {
	var align Align
	var hasAlign bool
	var wrap bool
	var collapse bool
	var collapseAt int
	var scrollable bool
	var height int
	for _, o := range b.opts {
		o(&align, &hasAlign, &wrap, &collapse, &collapseAt, &scrollable, &height)
	}
	var ap *Align
	if hasAlign {
		ap = &align
	}
	return RowElement{
		Children:   children,
		Align:      ap,
		Wrap:       wrap,
		Collapse:   collapse,
		CollapseAt: collapseAt,
	}
}

// LayoutOption configures a Column or Row.
type LayoutOption func(align *Align, hasAlign *bool, wrap *bool, collapse *bool, collapseAt *int, scrollable *bool, height *int)

func WithAlign(v, h AlignPosition) LayoutOption {
	return func(align *Align, hasAlign *bool, wrap *bool, collapse *bool, collapseAt *int, scrollable *bool, height *int) {
		*align = Align{V: v, H: h}
		*hasAlign = true
	}
}

func WithWrap(w bool) LayoutOption {
	return func(align *Align, hasAlign *bool, wrap *bool, collapse *bool, collapseAt *int, scrollable *bool, height *int) {
		*wrap = w
	}
}

func WithCollapse(c bool) LayoutOption {
	return func(align *Align, hasAlign *bool, wrap *bool, collapse *bool, collapseAt *int, scrollable *bool, height *int) {
		*collapse = c
	}
}

func WithCollapseAt(threshold int) LayoutOption {
	return func(align *Align, hasAlign *bool, wrap *bool, collapse *bool, collapseAt *int, scrollable *bool, height *int) {
		*collapseAt = threshold
		*collapse = true
	}
}

func WithScrollable(h int) LayoutOption {
	return func(align *Align, hasAlign *bool, wrap *bool, collapse *bool, collapseAt *int, scrollable *bool, height *int) {
		*scrollable = true
		*height = h
	}
}

// --- Container Elements ---

func Box(child Element, opts ...BoxOption) Element {
	b := BoxElement{Child: child}
	for _, opt := range opts {
		opt(&b)
	}
	return b
}

type BoxOption func(*BoxElement)

func WithBorder() BoxOption {
	return func(b *BoxElement) { b.Border = true }
}

func WithPadding(padding int) BoxOption {
	return func(b *BoxElement) { b.Padding = padding }
}

func WithWidth(width int) BoxOption {
	return func(b *BoxElement) { b.Width = width }
}

func WithHeight(height int) BoxOption {
	return func(b *BoxElement) { b.Height = height }
}

func WithTitle(title string) BoxOption {
	return func(b *BoxElement) { b.Title = title }
}

func Spacer(height int) Element {
	return SpacerElement{Height: height}
}

func View(children ...Element) Element {
	return Column(children...)
}

func Wrap(c Component, key string) Element {
	return ComponentElement{Component: c, Key: key}
}

// Scrollable returns a scrollable container with the given viewport height.
func Scrollable(child Element, height int) Element {
	return ScrollableElement{Child: child, Height: height}
}

// Form returns a form container.
func Form(onSubmit func(), children ...Element) Element {
	return FormElement{Children: children, OnSubmit: onSubmit}
}

// --- Decorative Elements ---

func Divider(opts ...DividerOption) Element {
	d := DividerElement{Style: "thin"}
	for _, opt := range opts {
		opt(&d)
	}
	return d
}

type DividerOption func(*DividerElement)

func DividerStyle(style string) DividerOption {
	return func(d *DividerElement) { d.Style = style }
}

func DividerWidth(width int) DividerOption {
	return func(d *DividerElement) { d.Width = width }
}

func DividerLabel(label string) DividerOption {
	return func(d *DividerElement) { d.Label = label }
}

func Progress(current, total int, opts ...ProgressOption) Element {
	p := ProgressElement{Current: current, Total: total, Width: 40}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

type ProgressOption func(*ProgressElement)

func ProgressWidth(width int) ProgressOption {
	return func(p *ProgressElement) { p.Width = width }
}

func ProgressLabel(label string) ProgressOption {
	return func(p *ProgressElement) { p.Label = label }
}

// Tabs returns a tab container.
func Tabs(tabs []Tab, active int, onChange func(int)) Element {
	return TabsElement{Tabs: tabs, Active: active, OnChange: onChange}
}

// Table returns a data table.
func Table(columns []TableColumn, rows [][]string, opts ...TableOption) Element {
	t := TableElement{Columns: columns, Rows: rows, Header: true}
	for _, opt := range opts {
		opt(&t)
	}
	return t
}

type TableOption func(*TableElement)

func TableNoHeader() TableOption {
	return func(t *TableElement) { t.Header = false }
}
