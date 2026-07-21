package react

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// interactiveEntry describes a focusable interactive element in the rendered tree.
type interactiveEntry struct {
	Type     string // "button" or "input"
	Label    string // displayed label (for buttons)
	OnClick  func() // set for buttons
	OnChange func(string) // set for inputs
	Value    string // current value (for inputs)
}

// ctxCache stores contexts for component instances, keyed by instance key.
// This ensures hooks persist across renders for nested components.
type ctxCache struct {
	entries map[string]*Context
}

func newCtxCache() *ctxCache {
	return &ctxCache{entries: make(map[string]*Context)}
}

// getOrCreate returns a context for the given component instance key.
func (c *ctxCache) getOrCreate(key string) *Context {
	if ctx, ok := c.entries[key]; ok {
		ctx.ResetHooks()
		return ctx
	}
	ctx := &Context{
		hooks: make([]any, 0),
	}
	c.entries[key] = ctx
	return ctx
}

// renderResult holds both the rendered string and the collected interactive entries.
type renderResult struct {
	view         string
	interactives []interactiveEntry
}

// renderElement renders an Element tree and collects interactive entries.
// focusIndex is the index of the currently focused interactive element, or -1.
func renderElement(el Element, width int, focusIndex int, cache *ctxCache) renderResult {
	componentCounter = 0
	var result renderResult
	var interactiveIdx int // tracks position in the interactive list during serialization
	result.view = serializeElement(el, width, &result.interactives, cache, focusIndex, &interactiveIdx)
	return result
}

// componentCounter is incremented during serialization to give each
// ComponentElement a unique position-based key for context caching.
var componentCounter uint64

// componentKey generates a stable key for a component instance.
func componentKey(c Component, key string) string {
	componentCounter++
	if key != "" {
		return key
	}
	typeName := fmt.Sprintf("%T", c)
	return fmt.Sprintf("%s-%d", typeName, componentCounter)
}

// serializeElement converts an Element tree into a rendered string.
// It appends interactive entries to the provided slice as it encounters them.
// focusIndex indicates which interactive element is focused (-1 = none).
// interactiveIdx tracks the current position in the interactive list.
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
		parts := make([]string, 0, len(e.Children))
		for _, child := range e.Children {
			s := serializeElement(child, width, interactives, cache, focusIndex, interactiveIdx)
			if s != "" {
				parts = append(parts, s)
			}
		}
		return lipgloss.JoinVertical(lipgloss.Top, parts...)

	case ColumnElement:
		return serializeElement(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *RowElement:
		if e == nil {
			return ""
		}
		parts := make([]string, 0, len(e.Children))
		for _, child := range e.Children {
			s := serializeElement(child, width, interactives, cache, focusIndex, interactiveIdx)
			if s != "" {
				parts = append(parts, s)
			}
		}
		return lipgloss.JoinHorizontal(lipgloss.Top, parts...)

	case RowElement:
		return serializeElement(&e, width, interactives, cache, focusIndex, interactiveIdx)

	case *BoxElement:
		if e == nil {
			return ""
		}
		childStr := serializeElement(e.Child, e.Width-4, interactives, cache, focusIndex, interactiveIdx)
		style := lipgloss.NewStyle()
		if e.Width > 0 {
			style = style.Width(e.Width)
		}
		if e.Height > 0 {
			style = style.Height(e.Height)
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
				Value:    e.Value,
			})
		}
		return renderInput(e, isFocused)

	case InputElement:
		return serializeElement(&e, width, interactives, cache, focusIndex, interactiveIdx)

	default:
		return ""
	}
}

// renderButton renders a button as styled text.
func renderButton(label string, focused bool) string {
	bg := lipgloss.Color("236")
	fg := lipgloss.Color("255")
	if focused {
		bg = lipgloss.Color("63")  // bright purple/blue for focus
		fg = lipgloss.Color("255") // white text
	}
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(fg).
		Background(bg).
		Bold(true)
	return " " + style.Render(" "+label+" ") + " "
}

// renderInput renders an input field.
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
