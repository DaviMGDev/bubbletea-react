package react

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// testComponent is a simple component used in tests.
type testComponent struct {
	initial int
}

func (c *testComponent) Render(ctx *Context) Element {
	count, setCount := UseState(ctx, c.initial)
	_ = setCount
	return View(
		Box(
			Column(
				Bold("Test"),
				Spacer(1),
				Textf("Count: %d", count),
			),
			WithBorder(),
			WithPadding(1),
			WithWidth(40),
			WithTitle("Test Box"),
		),
	)
}

func TestRender(t *testing.T) {
	cmp := &testComponent{initial: 42}
	ctx := &Context{
		hooks: make([]any, 0),
	}
	result := renderElement(cmp.Render(ctx), 80, -1, newCtxCache())
	if result.view == "" {
		t.Fatal("expected non-empty view string")
	}
	t.Logf("rendered view:\n%s", result.view)

	if len(result.interactives) != 0 {
		t.Errorf("expected 0 interactive elements, got %d", len(result.interactives))
	}
}

func TestUseState(t *testing.T) {
	cmp := &testComponent{initial: 10}
	ctx := &Context{
		hooks: make([]any, 0),
	}

	// First render
	el1 := cmp.Render(ctx)
	result1 := renderElement(el1, 80, -1, newCtxCache())
	if !contains(result1.view, "Count: 10") {
		t.Errorf("expected 'Count: 10' in view, got:\n%s", result1.view)
	}

	// Second render with same context (simulates re-render)
	ctx.ResetHooks()
	el2 := cmp.Render(ctx)
	result2 := renderElement(el2, 80, -1, newCtxCache())
	if !contains(result2.view, "Count: 10") {
		t.Errorf("expected 'Count: 10' in second render, got:\n%s", result2.view)
	}
}

func TestUseStateUpdate(t *testing.T) {
	ctx := &Context{
		hooks: make([]any, 0),
	}

	// Get state and setter
	val, setVal := UseState(ctx, 5)
	if val != 5 {
		t.Errorf("expected 5, got %d", val)
	}

	// Update state
	setVal(20)

	// Verify the hook value was updated (this is in-memory)
	hook := ctx.hooks[0].(*stateHook[int])
	if hook.value != 20 {
		t.Errorf("expected hook.value = 20, got %d", hook.value)
	}

	// Verify state on re-render
	ctx.ResetHooks()
	val2, _ := UseState(ctx, 5)
	if val2 != 20 {
		t.Errorf("expected 20 on re-render, got %d", val2)
	}
}

func TestHookOrderMismatch(t *testing.T) {
	ctx := &Context{
		hooks: make([]any, 0),
	}

	// First render: register hooks in order
	UseState(ctx, 1)
	UseState(ctx, "hello")
	_ = UseRef(ctx, 3.14)

	// Second render: correct order
	ctx.ResetHooks()
	UseState(ctx, 1)
	UseState(ctx, "hello")
	UseRef(ctx, 3.14)
	// Should not panic

	// Third render: wrong order (wrong type)
	ctx.ResetHooks()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on hook type mismatch")
		}
	}()
	UseState(ctx, 1)
	UseState(ctx, 99) // int instead of string - should panic
}

func TestUseRef(t *testing.T) {
	ctx := &Context{
		hooks: make([]any, 0),
	}

	ref := UseRef(ctx, "initial")
	if ref.Value != "initial" {
		t.Errorf("expected 'initial', got %q", ref.Value)
	}

	// Modify ref
	ref.Value = "changed"

	// Re-render: should see changed value
	ctx.ResetHooks()
	ref2 := UseRef(ctx, "ignored")
	if ref2.Value != "changed" {
		t.Errorf("expected 'changed' from cache, got %q", ref2.Value)
	}
}

func TestButtonRender(t *testing.T) {
	btn := ButtonElement{Label: "Click Me", OnClick: func() {}}
	result := renderButton(btn.Label, false)
	if !contains(result, "Click Me") {
		t.Errorf("expected 'Click Me' in button, got: %s", result)
	}
}

func TestInputRender(t *testing.T) {
	input := &InputElement{
		Value:       "hello",
		Placeholder: "type here",
		Width:       20,
	}
	result := renderInput(input, false)
	if !contains(result, "hello") {
		t.Errorf("expected 'hello' in input, got: %s", result)
	}
}

func TestInputRenderEmpty(t *testing.T) {
	input := &InputElement{
		Value:       "",
		Placeholder: "type here",
		Width:       20,
	}
	result := renderInput(input, false)
	if !contains(result, "type here") {
		t.Errorf("expected placeholder 'type here' in input, got: %s", result)
	}
}

func TestRowLayout(t *testing.T) {
	row := Row(
		Text("left"),
		Text("right"),
	)
	result := renderElement(row, 80, -1, newCtxCache())
	if result.view == "" {
		t.Fatal("expected non-empty view")
	}
	t.Logf("row view: %s", result.view)
}

func TestBoxWithTitle(t *testing.T) {
	box := Box(
		Text("content"),
		WithBorder(),
		WithPadding(1),
		WithWidth(30),
		WithTitle("My Box"),
	)
	result := renderElement(box, 80, -1, newCtxCache())
	if result.view == "" {
		t.Fatal("expected non-empty view")
	}
	t.Logf("box view:\n%s", result.view)
	if !contains(result.view, "My Box") {
		t.Errorf("expected 'My Box' in view, got:\n%s", result.view)
	}
}

func TestEmptyComponent(t *testing.T) {
	row := Column()
	result := renderElement(row, 80, -1, newCtxCache())
	// Empty column should produce empty string
	// This is valid - no panic expected
	_ = result
}

func TestNilElement(t *testing.T) {
	result := renderElement(nil, 80, -1, newCtxCache())
	if result.view != "" {
		t.Errorf("expected empty view for nil element, got: %s", result.view)
	}
}

func TestQuitKeys(t *testing.T) {
	m := &rootModel{
		root: &testComponent{initial: 0},
		ctx: &Context{
			hooks: make([]any, 0),
		},
		ctxCache: newCtxCache(),
	}

	// Ctrl+C should quit
	t.Run("ctrl+c", func(t *testing.T) {
		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		if cmd == nil {
			t.Fatal("expected tea.Quit cmd for Ctrl+C")
		}
		msg := cmd()
		if _, ok := msg.(tea.QuitMsg); !ok {
			t.Errorf("expected QuitMsg, got %T", msg)
		}
	})

	// Reset for next test
	m.quitting = false

	// Escape should quit
	t.Run("escape", func(t *testing.T) {
		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEscape})
		if cmd == nil {
			t.Error("expected tea.Quit cmd for Escape")
		}
	})

	m.quitting = false

	// 'q' should quit (no input focused)
	t.Run("q key", func(t *testing.T) {
		// No interactives, so focusIndex = -1
		_, cmd := m.Update(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'q'},
		})
		if cmd == nil {
			t.Error("expected tea.Quit cmd for 'q' key")
		}
	})

	m.quitting = false

	// 'a' should NOT quit
	t.Run("non-quit key", func(t *testing.T) {
		_, cmd := m.Update(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'a'},
		})
		if cmd != nil {
			t.Error("expected nil cmd for non-quit key 'a'")
		}
	})
}

// contains reports whether s contains substr.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
