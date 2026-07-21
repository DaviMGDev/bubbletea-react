package react

import (
	tea "github.com/charmbracelet/bubbletea"
)

// rootModel implements tea.Model by wrapping a root Component.
type rootModel struct {
	root         Component
	ctx          *Context
	width        int
	height       int
	viewString   string
	interactives []interactiveEntry
	focusIndex   int
	ctxCache     *ctxCache
	quitting     bool
	prevFocus    int // previous focus index for firing blur events
}

func (m *rootModel) Init() tea.Cmd {
	m.renderComponent()
	return nil
}

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
		m.renderComponent()

	default:
	}

	return m, nil
}

func (m *rootModel) View() string {
	if m.quitting {
		return ""
	}
	return m.viewString
}

func (m *rootModel) isQuitKey(msg tea.KeyMsg) bool {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEscape:
		return true
	case tea.KeyRunes:
		if m.focusIndex >= 0 && m.focusIndex < len(m.interactives) {
			if m.interactives[m.focusIndex].Type == "input" {
				return false
			}
		}
		return len(msg.Runes) == 1 && msg.Runes[0] == 'q'
	}
	return false
}

func (m *rootModel) renderComponent() {
	m.ctx.ResetHooks()
	result := renderElement(m.root.Render(m.ctx), m.width, m.focusIndex, m.ctxCache)
	m.viewString = result.view
	m.interactives = result.interactives

	if len(m.interactives) == 0 {
		m.focusIndex = -1
	} else if m.focusIndex >= len(m.interactives) {
		m.focusIndex = len(m.interactives) - 1
	} else if m.focusIndex < 0 && len(m.interactives) > 0 {
		m.focusIndex = 0
	}
}

func (m *rootModel) handleKeyMsg(msg tea.KeyMsg) {
	if len(m.interactives) == 0 {
		return
	}

	if m.focusIndex < 0 || m.focusIndex >= len(m.interactives) {
		// No interactives at all — navigate keys do nothing
		switch msg.Type {
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
			return
		}
		return
	}

	entry := m.interactives[m.focusIndex]
	currentRow := entry.RowID

	// Navigation keys
	switch msg.Type {
	case tea.KeyTab, tea.KeyDown:
		// Vertical: next element in a DIFFERENT row
		m.prevFocus = m.focusIndex
		m.focusIndex = nextDifferentRow(m.interactives, m.focusIndex, currentRow, +1)
		m.fireFocusBlur()
		return
	case tea.KeyShiftTab, tea.KeyUp:
		// Vertical: previous element in a DIFFERENT row
		m.prevFocus = m.focusIndex
		m.focusIndex = prevDifferentRow(m.interactives, m.focusIndex, currentRow)
		m.fireFocusBlur()
		return
	}

	// Route non-navigation keys to focused element
	leftRightConsumed := false

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

	case "checkbox":
		switch msg.Type {
		case tea.KeyEnter, tea.KeySpace:
			if entry.OnToggle != nil {
				entry.OnToggle(!entry.Checked)
			}
		}

	case "select":
		if entry.OptionCount > 0 && len(entry.OptionValues) > 0 {
			switch msg.Type {
			case tea.KeyRight:
				newIdx := (entry.SelectedIndex + 1) % entry.OptionCount
				if newIdx >= 0 && newIdx < len(entry.OptionValues) && entry.OnChange != nil {
					entry.OnChange(entry.OptionValues[newIdx])
					leftRightConsumed = true
				}
			case tea.KeyLeft:
				newIdx := (entry.SelectedIndex - 1 + entry.OptionCount) % entry.OptionCount
				if newIdx >= 0 && newIdx < len(entry.OptionValues) && entry.OnChange != nil {
					entry.OnChange(entry.OptionValues[newIdx])
					leftRightConsumed = true
				}
			}
		}

	case "tab":
		switch msg.Type {
		case tea.KeyEnter, tea.KeySpace:
			if entry.OnClick != nil {
				entry.OnClick()
			}
		}

	case "form":
		switch msg.Type {
		case tea.KeyEnter, tea.KeySpace:
			if entry.OnClick != nil {
				entry.OnClick()
			}
		}
	}

	// Left/Right: same-row navigation (unless consumed by Select)
	if !leftRightConsumed {
		switch msg.Type {
		case tea.KeyRight:
			if next := nextSameRow(m.interactives, m.focusIndex, currentRow); next >= 0 {
				m.prevFocus = m.focusIndex
				m.focusIndex = next
				m.fireFocusBlur()
			}
		case tea.KeyLeft:
			if prev := prevSameRow(m.interactives, m.focusIndex, currentRow); prev >= 0 {
				m.prevFocus = m.focusIndex
				m.focusIndex = prev
				m.fireFocusBlur()
			}
		}
	}
}

// nextDifferentRow finds the next element forward (wrapping) with a different row ID.
func nextDifferentRow(entries []interactiveEntry, from int, rowID int, dir int) int {
	n := len(entries)
	for i := 1; i < n; i++ {
		idx := (from + i*dir + n) % n
		if entries[idx].RowID != rowID {
			return idx
		}
	}
	return from
}

// prevDifferentRow finds the previous element with a different row ID.
func prevDifferentRow(entries []interactiveEntry, from int, rowID int) int {
	return nextDifferentRow(entries, from, rowID, -1)
}

// nextSameRow finds the next element forward with the same row ID.
// Returns -1 if none found.
func nextSameRow(entries []interactiveEntry, from int, rowID int) int {
	n := len(entries)
	for i := 1; i < n; i++ {
		idx := (from + i) % n
		if entries[idx].RowID == rowID {
			return idx
		}
	}
	return -1
}

// prevSameRow finds the previous element with the same row ID.
// Returns -1 if none found.
func prevSameRow(entries []interactiveEntry, from int, rowID int) int {
	n := len(entries)
	for i := 1; i < n; i++ {
		idx := (from - i + n) % n
		if entries[idx].RowID == rowID {
			return idx
		}
	}
	return -1
}

func (m *rootModel) fireFocusBlur() {
	// Fire blur on previously focused element
	if m.prevFocus >= 0 && m.prevFocus < len(m.interactives) {
		prev := m.interactives[m.prevFocus]
		if prev.OnBlur != nil {
			prev.OnBlur()
		}
	}
	// Fire focus on new element
	if m.focusIndex >= 0 && m.focusIndex < len(m.interactives) {
		cur := m.interactives[m.focusIndex]
		if cur.OnFocus != nil {
			cur.OnFocus()
		}
	}
}

func Root(root Component) *tea.Program {
	m := &rootModel{
		root: root,
		ctx: &Context{
			hooks: make([]any, 0),
		},
		width:      80,
		height:     24,
		focusIndex: 0,
		prevFocus:  -1,
		ctxCache:   newCtxCache(),
	}

	return tea.NewProgram(m)
}
