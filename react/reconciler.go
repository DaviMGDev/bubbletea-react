package react

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// interactiveEntry describes a focusable interactive element in the rendered tree.
type interactiveEntry struct {
	Type     string // "button", "input", "checkbox", "select", "tab", "form"
	Label    string // displayed label
	OnClick  func() // activation handler (buttons, tabs, form submit)
	OnChange func(string) // value change handler (inputs, select)
	OnSubmit func() // called on Enter (input-specific)
	OnFocus  func()
	OnBlur   func()
	Value    string // current value (inputs) or serialized state
	Checked  bool   // for checkboxes
	OnToggle func(bool) // for checkboxes
	RowID    int    // which visual row this element belongs to
	// Select-specific
	SelectedIndex int
	OptionCount   int
	OptionValues  []string // option values for cycling
}

// ctxCache stores contexts for component instances, keyed by instance key.
type ctxCache struct {
	entries map[string]*Context
}

func newCtxCache() *ctxCache {
	return &ctxCache{entries: make(map[string]*Context)}
}

func (c *ctxCache) getOrCreate(key string) *Context {
	if ctx, ok := c.entries[key]; ok {
		ctx.ResetHooks()
		return ctx
	}
	ctx := &Context{hooks: make([]any, 0)}
	c.entries[key] = ctx
	return ctx
}

// renderResult holds both the rendered string and the collected interactive entries.
type renderResult struct {
	view         string
	interactives []interactiveEntry
}

// currentRowID tracks visual rows for spatial navigation.
// Reset each render, incremented at each new row boundary.
// Safe because Bubble Tea is single-threaded.
var currentRowID int

// componentCounter is incremented during serialization to give each
// ComponentElement a unique position-based key for context caching.
var componentCounter uint64

// renderElement renders an Element tree and collects interactive entries.
func renderElement(el Element, width int, focusIndex int, cache *ctxCache) renderResult {
	componentCounter = 0
	currentRowID = 0
	var result renderResult
	var interactiveIdx int
	result.view = serializeElement(el, width, &result.interactives, cache, focusIndex, &interactiveIdx)
	return result
}

func componentKey(c Component, key string) string {
	componentCounter++
	if key != "" {
		return key
	}
	typeName := fmt.Sprintf("%T", c)
	return fmt.Sprintf("%s-%d", typeName, componentCounter)
}

// nextRow bumps the row counter and returns the new value.
func nextRow() int {
	currentRowID++
	return currentRowID
}

// serializeElement converts an Element tree into a rendered string.
func serializeElement(el Element, width int, interactives *[]interactiveEntry, cache *ctxCache, focusIndex int, interactiveIdx *int) string {
	if el == nil {
		return ""
	}

	switch e := el.(type) {
	case *TextElement:
		if e == nil {
			return ""
		}
		return e.Text
	case TextElement:
		return e.Text

	case *BoldElement:
		if e == nil {
			return ""
		}
		return lipgloss.NewStyle().Bold(true).Render(e.Text)
	case BoldElement:
		return lipgloss.NewStyle().Bold(true).Render(e.Text)

	case *ComponentElement:
		if e == nil || e.Component == nil {
			return ""
		}
		key := componentKey(e.Component, e.Key)
		childCtx := cache.getOrCreate(key)
		childEl := e.Component.Render(childCtx)
		return serializeElement(childEl, width, interactives, cache, focusIndex, interactiveIdx)
	case ComponentElement:
		return serializeElement(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *ColumnElement:
		if e == nil {
			return ""
		}
		return serializeColumn(e, width, interactives, cache, focusIndex, interactiveIdx)
	case ColumnElement:
		return serializeColumn(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *RowElement:
		if e == nil {
			return ""
		}
		return serializeRow(e, width, interactives, cache, focusIndex, interactiveIdx)
	case RowElement:
		return serializeRow(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *BoxElement:
		if e == nil {
			return ""
		}
		// Calculate actual overhead: border (2 if present) + padding (2 * value).
		overhead := 2 * e.Padding
		if e.Border {
			overhead += 2
		}
		childWidth := e.Width - overhead
		if childWidth < 0 {
			childWidth = width - overhead
		}
		childStr := serializeElement(e.Child, childWidth, interactives, cache, focusIndex, interactiveIdx)
		style := lipgloss.NewStyle()
		if e.Width > 0 {
			style = style.Width(e.Width)
		}
		if e.Height > 0 {
			style = style.MaxHeight(e.Height)
		}
		if e.Padding > 0 {
			style = style.Padding(e.Padding)
		}
		if e.Border {
			style = style.Border(lipgloss.RoundedBorder())
		}
		if e.Title != "" {
			childStr = e.Title + "\n\n" + childStr
		}
		if childStr == "" {
			return style.Render("")
		}
		return style.Render(childStr)
	case BoxElement:
		return serializeElement(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *SpacerElement:
		if e == nil || e.Height <= 0 {
			return ""
		}
		return strings.Repeat("\n", e.Height)
	case SpacerElement:
		return serializeElement(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *ButtonElement:
		if e == nil {
			return ""
		}
		currentIdx := *interactiveIdx
		*interactiveIdx++
		isFocused := (focusIndex >= 0 && currentIdx == focusIndex)
		if interactives != nil {
			*interactives = append(*interactives, interactiveEntry{
				Type:    "button",
				Label:   e.Label,
				OnClick: e.OnClick,
				OnFocus: e.OnFocus,
				OnBlur:  e.OnBlur,
				RowID:   currentRowID,
			})
		}
		return renderButton(e.Label, isFocused)
	case ButtonElement:
		return serializeElement(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *InputElement:
		if e == nil {
			return ""
		}
		currentIdx := *interactiveIdx
		*interactiveIdx++
		isFocused := (focusIndex >= 0 && currentIdx == focusIndex)
		if interactives != nil {
			*interactives = append(*interactives, interactiveEntry{
				Type:     "input",
				OnChange: e.OnChange,
				OnSubmit: e.OnSubmit,
				OnFocus:  e.OnFocus,
				OnBlur:   e.OnBlur,
				Value:    e.Value,
				RowID:    currentRowID,
			})
		}
		return renderInput(e, isFocused)
	case InputElement:
		return serializeElement(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *ScrollableElement:
		if e == nil {
			return ""
		}
		return serializeScrollable(e, width, interactives, cache, focusIndex, interactiveIdx)
	case ScrollableElement:
		return serializeScrollable(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *FormElement:
		if e == nil {
			return ""
		}
		return serializeForm(e, width, interactives, cache, focusIndex, interactiveIdx)
	case FormElement:
		return serializeForm(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *CheckboxElement:
		if e == nil {
			return ""
		}
		currentIdx := *interactiveIdx
		*interactiveIdx++
		isFocused := (focusIndex >= 0 && currentIdx == focusIndex)
		if interactives != nil {
			*interactives = append(*interactives, interactiveEntry{
				Type:     "checkbox",
				Label:    e.Label,
				Checked:  e.Checked,
				OnToggle: e.OnChange,
				RowID:    currentRowID,
			})
		}
		return renderCheckbox(e, isFocused)
	case CheckboxElement:
		return serializeElement(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *SelectElement:
		if e == nil {
			return ""
		}
		currentIdx := *interactiveIdx
		*interactiveIdx++
		isFocused := (focusIndex >= 0 && currentIdx == focusIndex)
		optCount := len(e.Options)
		optValues := make([]string, optCount)
		for i, o := range e.Options {
			optValues[i] = o.Value
		}
		if interactives != nil {
			*interactives = append(*interactives, interactiveEntry{
				Type:          "select",
				OnChange: func(val string) {
					if e.OnChange != nil {
						e.OnChange(val)
					}
				},
				Value:         fmt.Sprintf("%d", e.Selected),
				SelectedIndex: e.Selected,
				OptionCount:   optCount,
				OptionValues:  optValues,
				RowID:         currentRowID,
			})
		}
		return renderSelect(e, isFocused)
	case SelectElement:
		return serializeElement(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *DividerElement:
		if e == nil {
			return ""
		}
		return renderDivider(e, width)
	case DividerElement:
		return renderDivider(&e, width)

	case *ProgressElement:
		if e == nil {
			return ""
		}
		return renderProgress(e, width)
	case ProgressElement:
		return renderProgress(&e, width)

	case *TabsElement:
		if e == nil {
			return ""
		}
		return serializeTabs(e, width, interactives, cache, focusIndex, interactiveIdx)
	case TabsElement:
		return serializeTabs(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *TableElement:
		if e == nil {
			return ""
		}
		return renderTable(e, width)
	case TableElement:
		return renderTable(&e, width)

	default:
		return ""
	}
}

// --- Column ---

func serializeColumn(e *ColumnElement, width int, interactives *[]interactiveEntry, cache *ctxCache, focusIndex int, interactiveIdx *int) string {
	parts := make([]string, 0, len(e.Children))
	for _, child := range e.Children {
		nextRow() // each Column child starts a new visual row
		s := serializeElement(child, width, interactives, cache, focusIndex, interactiveIdx)
		if s != "" {
			parts = append(parts, s)
		}
	}
	joined := lipgloss.JoinVertical(lipgloss.Top, parts...)
	if e.Align != nil {
		joined = lipgloss.JoinVertical(toLipglossPosition(e.Align.V), parts...)
	}
	if e.Scrollable && e.Height > 0 {
		lines := strings.Split(joined, "\n")
		if len(lines) > e.Height {
			lines = lines[:e.Height]
		}
		joined = strings.Join(lines, "\n")
		style := lipgloss.NewStyle().Height(e.Height)
		joined = style.Render(joined)
	}
	return joined
}

// --- Row ---

func serializeRow(e *RowElement, width int, interactives *[]interactiveEntry, cache *ctxCache, focusIndex int, interactiveIdx *int) string {
	if e.Collapse {
		threshold := e.CollapseAt
		if threshold <= 0 {
			threshold = 80
		}
		if width < threshold {
			parts := make([]string, 0, len(e.Children))
			for _, child := range e.Children {
				nextRow()
				s := serializeElement(child, width, interactives, cache, focusIndex, interactiveIdx)
				if s != "" {
					parts = append(parts, s)
				}
			}
			joined := lipgloss.JoinVertical(lipgloss.Top, parts...)
			if e.Align != nil {
				joined = lipgloss.JoinVertical(toLipglossPosition(e.Align.V), parts...)
			}
			return joined
		}
	}

	// All children in a Row share the same row ID
	nextRow()
	thisRow := currentRowID
	parts := make([]string, 0, len(e.Children))
	savedRow := currentRowID

	for _, child := range e.Children {
		currentRowID = thisRow // force same row for all Row children
		s := serializeElement(child, width, interactives, cache, focusIndex, interactiveIdx)
		if s != "" {
			parts = append(parts, s)
		}
	}
	currentRowID = savedRow // restore (Row doesn't consume a row for nested content)

	if e.Wrap {
		return joinWrapped(parts, width)
	}
	joined := lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	if e.Align != nil {
		joined = lipgloss.JoinHorizontal(toLipglossPosition(e.Align.H), parts...)
	}
	return joined
}

func joinWrapped(parts []string, width int) string {
	if len(parts) == 0 {
		return ""
	}

	// Split each part into lines and find the max number of lines.
	maxLines := 0
	partLines := make([][]string, len(parts))
	for i, part := range parts {
		partLines[i] = strings.Split(part, "\n")
		if len(partLines[i]) > maxLines {
			maxLines = len(partLines[i])
		}
	}

	// Pad shorter parts with empty lines so all have the same height.
	for i := range partLines {
		for len(partLines[i]) < maxLines {
			partLines[i] = append(partLines[i], "")
		}
	}

	// Build output: for each line index, join all parts horizontally.
	var lines []string
	for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
		var rowParts []string
		for _, pl := range partLines {
			rowParts = append(rowParts, pl[lineIdx])
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, rowParts...))
	}
	return strings.Join(lines, "\n")
}

// --- Scrollable ---

func serializeScrollable(e *ScrollableElement, width int, interactives *[]interactiveEntry, cache *ctxCache, focusIndex int, interactiveIdx *int) string {
	child := serializeElement(e.Child, width, interactives, cache, focusIndex, interactiveIdx)
	h := e.Height
	if h <= 0 {
		h = 10
	}
	lines := strings.Split(child, "\n")
	if len(lines) > h {
		lines = lines[:h]
	}
	style := lipgloss.NewStyle().Height(h).MaxHeight(h)
	return style.Render(strings.Join(lines, "\n"))
}

// --- Form ---

func serializeForm(e *FormElement, width int, interactives *[]interactiveEntry, cache *ctxCache, focusIndex int, interactiveIdx *int) string {
	parts := make([]string, 0, len(e.Children))
	for _, child := range e.Children {
		s := serializeElement(child, width, interactives, cache, focusIndex, interactiveIdx)
		if s != "" {
			parts = append(parts, s)
		}
	}
	if e.OnSubmit != nil {
		nextRow()
		currentIdx := *interactiveIdx
		*interactiveIdx++
		_ = currentIdx
		if interactives != nil {
			*interactives = append(*interactives, interactiveEntry{
				Type:    "form",
				Label:   "Submit",
				OnClick: e.OnSubmit,
				RowID:   currentRowID,
			})
		}
	}
	return strings.Join(parts, "\n")
}

// --- Checkbox ---

func renderCheckbox(e *CheckboxElement, focused bool) string {
	mark := " "
	if e.Checked {
		mark = "x"
	}
	box := fmt.Sprintf("[%s]", mark)
	style := lipgloss.NewStyle()
	if focused {
		style = style.Foreground(lipgloss.Color("63")).Bold(true)
	}
	return style.Render(box) + " " + e.Label
}

// --- Select ---

func renderSelect(e *SelectElement, focused bool) string {
	if len(e.Options) == 0 {
		return "(no options)"
	}
	sel := e.Selected
	if sel < 0 || sel >= len(e.Options) {
		sel = 0
	}
	opt := e.Options[sel]
	style := lipgloss.NewStyle().Padding(0, 1)
	if focused {
		style = style.Background(lipgloss.Color("236")).Foreground(lipgloss.Color("255"))
	}
	return style.Render(fmt.Sprintf("▸ %s", opt.Label))
}

// --- Divider ---

func renderDivider(e *DividerElement, width int) string {
	w := e.Width
	if w <= 0 {
		w = width
	}
	var char string
	switch e.Style {
	case "double":
		char = "═"
	case "dashed":
		char = " ─"
	default:
		char = "─"
	}
	line := strings.Repeat(char, w)
	if e.Label != "" {
		half := w / 2
		labelStart := half - len(e.Label)/2
		if labelStart < 0 {
			labelStart = 0
		}
		runes := []rune(line)
		for i, ch := range e.Label {
			pos := labelStart + i
			if pos < len(runes) {
				runes[pos] = ch
			}
		}
		line = string(runes)
	}
	return line
}

// --- Progress ---

func renderProgress(e *ProgressElement, width int) string {
	w := e.Width
	if w <= 0 {
		w = width
	}
	pct := 0.0
	if e.Total > 0 {
		pct = float64(e.Current) / float64(e.Total)
	}
	filled := int(pct * float64(w))
	if filled > w {
		filled = w
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", w-filled)
	label := fmt.Sprintf(" %d/%d", e.Current, e.Total)
	if e.Label != "" {
		label = " " + e.Label
	}
	return bar + label
}

// --- Tabs ---

func serializeTabs(e *TabsElement, width int, interactives *[]interactiveEntry, cache *ctxCache, focusIndex int, interactiveIdx *int) string {
	if len(e.Tabs) == 0 {
		return ""
	}

	// All tab headers share the same row (they're laid out horizontally)
	nextRow()
	thisRow := currentRowID
	savedRow := currentRowID

	headers := make([]string, 0, len(e.Tabs))
	for i, tab := range e.Tabs {
		currentRowID = thisRow // force same row for all tabs
		idx := i
		currentHeaderIdx := *interactiveIdx
		*interactiveIdx++
		headerFocused := (focusIndex >= 0 && currentHeaderIdx == focusIndex)
		if interactives != nil {
			*interactives = append(*interactives, interactiveEntry{
				Type:  "tab",
				Label: tab.Label,
				OnClick: func() {
					if e.OnChange != nil {
						e.OnChange(idx)
					}
				},
				Value: fmt.Sprintf("%d", idx),
				RowID: currentRowID,
			})
		}
		style := lipgloss.NewStyle().Padding(0, 2)
		if i == e.Active {
			style = style.Bold(true).Foreground(lipgloss.Color("63"))
		}
		if headerFocused {
			style = style.Background(lipgloss.Color("236"))
		}
		headers = append(headers, style.Render(tab.Label))
	}
	currentRowID = savedRow // restore (tabs don't consume extra rows)

	headerLine := lipgloss.JoinHorizontal(lipgloss.Top, headers...)
	divider := strings.Repeat("─", width)
	active := e.Active
	if active < 0 || active >= len(e.Tabs) {
		active = 0
	}
	content := serializeElement(e.Tabs[active].Content, width, interactives, cache, focusIndex, interactiveIdx)
	return headerLine + "\n" + divider + "\n" + content
}

// --- Table ---

func renderTable(e *TableElement, width int) string {
	if len(e.Columns) == 0 || len(e.Rows) == 0 {
		return ""
	}
	colWidths := make([]int, len(e.Columns))
	totalFixed := 0
	autoCols := 0
	for i, col := range e.Columns {
		if col.Width > 0 {
			colWidths[i] = col.Width
			totalFixed += col.Width
		} else {
			autoCols++
		}
	}
	remaining := width - totalFixed - (len(e.Columns) - 1)
	if remaining > 0 && autoCols > 0 {
		perCol := remaining / autoCols
		for i := range colWidths {
			if colWidths[i] == 0 {
				colWidths[i] = perCol
			}
		}
	}
	for i, col := range e.Columns {
		minW := runewidth.StringWidth(col.Label) + 2
		for _, row := range e.Rows {
			if i < len(row) {
				if w := runewidth.StringWidth(row[i]) + 2; w > minW {
					minW = w
				}
			}
		}
		if colWidths[i] < minW {
			colWidths[i] = minW
		}
	}
	var lines []string
	if e.Header {
		cells := make([]string, len(e.Columns))
		for i, col := range e.Columns {
			cells[i] = padOrTrunc(col.Label, colWidths[i])
		}
		headerStyle := lipgloss.NewStyle().Bold(true)
		lines = append(lines, headerStyle.Render(strings.Join(cells, " ")))
		sepParts := make([]string, len(e.Columns))
		for i, w := range colWidths {
			sepParts[i] = strings.Repeat("─", w)
		}
		lines = append(lines, strings.Join(sepParts, " "))
	}
	for _, row := range e.Rows {
		cells := make([]string, len(e.Columns))
		for i := range e.Columns {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			cells[i] = padOrTrunc(val, colWidths[i])
		}
		lines = append(lines, strings.Join(cells, " "))
	}
	return strings.Join(lines, "\n")
}

// --- Helpers ---

func toLipglossPosition(p AlignPosition) lipgloss.Position {
	switch p {
	case Top:
		return lipgloss.Top
	case Bottom:
		return lipgloss.Bottom
	case Center:
		return lipgloss.Center
	default:
		return lipgloss.Top
	}
}

func padOrTrunc(s string, w int) string {
	sw := runewidth.StringWidth(s)
	if sw >= w {
		return runewidth.Truncate(s, w, "…")
	}
	return s + strings.Repeat(" ", w-sw)
}

func renderButton(label string, focused bool) string {
	bg := lipgloss.Color("236")
	fg := lipgloss.Color("255")
	if focused {
		bg = lipgloss.Color("63")
		fg = lipgloss.Color("255")
	}
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(fg).
		Background(bg).
		Bold(true)
	return " " + style.Render(" "+label+" ") + " "
}

func renderInput(e *InputElement, focused bool) string {
	displayText := e.Value
	if displayText == "" {
		displayText = e.Placeholder
	}
	w := e.Width
	if w <= 0 {
		w = 30
	}
	if len(displayText) > w {
		displayText = displayText[len(displayText)-w:]
	}
	borderColor := lipgloss.Color("240")
	if focused {
		borderColor = lipgloss.Color("63")
	}
	style := lipgloss.NewStyle().
		Width(w).
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor)
	padding := w - len(displayText)
	if padding < 0 {
		padding = 0
	}
	return style.Render(displayText + strings.Repeat(" ", padding))
}
