package react

import (
	tea "github.com/charmbracelet/bubbletea"
)

// rootModel implements tea.Model by wrapping a root Component.
// It manages the component rendering lifecycle and keyboard routing.
type rootModel struct {
	root         Component     // the root component
	ctx          *Context      // hooks context for the root component
	width        int           // terminal width
	height       int           // terminal height
	viewString   string        // cached view output
	interactives []interactiveEntry // collected interactive elements
	focusIndex   int           // currently focused interactive element index (-1 = none)
	ctxCache     *ctxCache     // persistent context cache for nested components
	quitting     bool          // set when a quit key is pressed
}

// Init implements tea.Model.
func (m *rootModel) Init() tea.Cmd {
	m.renderComponent()
	return nil
}

// Update implements tea.Model.
func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.renderComponent()

	case tea.KeyMsg:
		if m.isQuitKey(msg) {
			m.quitting = true
			return m, tea.Quit
		}
		m.handleKeyMsg(msg)
		// Always re-render after a key event — hook setters may have been
		// called inside handleKeyMsg and updated state in-place.
		// NOTE: we do NOT call program.Send from inside setters because
		// Bubble Tea's message channel is unbuffered and sending from
		// within Update would deadlock.
		m.renderComponent()

	default:
		// Unknown messages are ignored.
	}

	return m, nil
}

// View implements tea.Model.
func (m *rootModel) View() string {
	if m.quitting {
		return ""
	}
	return m.viewString
}

// isQuitKey returns true if the key signals the program should quit.
func (m *rootModel) isQuitKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEscape:
		return true
	case tea.KeyRunes:
		// 'q' to quit when no input field is focused.
		if m.focusIndex >= 0 && m.focusIndex < len(m.interactives) {
			if m.interactives[m.focusIndex].Type == "input" {
				return false
			}
		}
		return len(msg.Runes) == 1 && msg.Runes[0] == 'q'
	}
	return false
}

// renderComponent renders the root component and caches the view string.
func (m *rootModel) renderComponent() {
	m.ctx.ResetHooks()
	result := renderElement(m.root.Render(m.ctx), m.width, m.focusIndex, m.ctxCache)
	m.viewString = result.view
	m.interactives = result.interactives

	// Clamp focus index to valid range.
	if len(m.interactives) == 0 {
		m.focusIndex = -1
	} else if m.focusIndex >= len(m.interactives) {
		m.focusIndex = len(m.interactives) - 1
	} else if m.focusIndex < 0 && len(m.interactives) > 0 {
		m.focusIndex = 0
	}
}

// handleKeyMsg routes keyboard input to the currently focused interactive element.
func (m *rootModel) handleKeyMsg(msg tea.KeyMsg) {
	if len(m.interactives) == 0 {
		return
	}

	// Handle navigation keys regardless of focus type.
	switch msg.Type {
	case tea.KeyTab, tea.KeyDown:
		m.focusIndex = (m.focusIndex + 1) % len(m.interactives)
		return

	case tea.KeyShiftTab, tea.KeyUp:
		m.focusIndex = (m.focusIndex - 1 + len(m.interactives)) % len(m.interactives)
		return
	}

	// Route the key to the focused interactive element.
	if m.focusIndex < 0 || m.focusIndex >= len(m.interactives) {
		return
	}

	entry := m.interactives[m.focusIndex]

	switch entry.Type {
	case "button":
		switch msg.Type {
		case tea.KeyEnter, tea.KeySpace:
			if entry.OnClick != nil {
				entry.OnClick()
			}
		}

	case "input":
		switch msg.Type {
		case tea.KeyRunes:
			if entry.OnChange != nil {
				newValue := entry.Value + string(msg.Runes)
				entry.OnChange(newValue)
			}

		case tea.KeySpace:
			if entry.OnChange != nil {
				newValue := entry.Value + " "
				entry.OnChange(newValue)
			}

		case tea.KeyBackspace:
			if entry.OnChange != nil && len(entry.Value) > 0 {
				newValue := entry.Value[:len(entry.Value)-1]
				entry.OnChange(newValue)
			}
		}
	}
}

// Root creates a new Bubble Tea program from the given root component.
// The component's Render method will be called to produce the initial view,
// and state changes will trigger automatic re-renders.
//
// Press Ctrl+C, Escape, or q (when no input is focused) to quit.
//
// Example:
//
//	program := react.Root(&Counter{})
//	if _, err := program.Run(); err != nil {
//	    panic(err)
//	}
func Root(root Component) *tea.Program {
	m := &rootModel{
		root: root,
		ctx: &Context{
			hooks: make([]any, 0),
		},
		width:      80,
		height:     24,
		focusIndex: 0,
		ctxCache:   newCtxCache(),
	}

	return tea.NewProgram(m)
}
