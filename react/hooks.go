package react

import (
	"fmt"
	"reflect"
)

// --- useState ---

// stateHook holds one piece of state for a component instance.
type stateHook[T any] struct {
	value T
}

// UseState returns a stateful value and a function to update it.
// The hook is identified by its call order within the component
// (like React's rules of hooks). State updates are in-place; the
// bridge triggers a re-render after each key event.
//
//	count, setCount := react.UseState(ctx, 0)
func UseState[T any](ctx *Context, initial T) (T, func(T)) {
	if ctx.hookIndex >= len(ctx.hooks) {
		// First render: create a new hook.
		h := &stateHook[T]{value: initial}
		ctx.hooks = append(ctx.hooks, h)
		idx := ctx.hookIndex
		ctx.hookIndex++
		hooks := ctx.hooks
		return h.value, func(val T) {
			hooks[idx].(*stateHook[T]).value = val
		}
	}

	// Re-render: retrieve existing hook.
	h, ok := ctx.hooks[ctx.hookIndex].(*stateHook[T])
	if !ok {
		expected := reflect.TypeOf((*T)(nil)).Elem()
		got := reflect.TypeOf(ctx.hooks[ctx.hookIndex])
		panic(fmt.Sprintf(
			"bubbletea-react: hook type mismatch at index %d: expected %v, got %v. "+
				"Hooks must be called in the same order and with the same types on every render.",
			ctx.hookIndex, expected, got,
		))
	}
	ctx.hookIndex++
	return h.value, func(val T) {
		h.value = val
	}
}

// --- useEffect ---

// effectHook stores an effect and its dependency array.
type effectHook struct {
	cleanup  func()
	lastDeps []any
}

// UseEffect runs a side effect function after every render where any of the
// specified dependencies have changed. If deps is nil, the effect runs after
// every render. If deps is an empty slice, the effect runs only once (on mount).
//
// The effect function can return an optional cleanup function:
//
//	react.UseEffect(ctx, func() func() {
//	    // setup
//	    return func() {
//	        // cleanup
//	    }
//	}, []any{count})
func UseEffect(ctx *Context, fn func() func(), deps []any) {
	if ctx.hookIndex >= len(ctx.hooks) {
		// First render: register effect.
		h := &effectHook{}
		ctx.hooks = append(ctx.hooks, h)
		// Run effect immediately.
		if fn != nil {
			if h.cleanup != nil {
				h.cleanup()
			}
			h.cleanup = fn()
		}
		h.lastDeps = cloneDeps(deps)
		ctx.hookIndex++
		return
	}

	h := ctx.hooks[ctx.hookIndex].(*effectHook)
	ctx.hookIndex++

	// Check if deps changed.
	if deps == nil || depsChanged(h.lastDeps, deps) {
		if h.cleanup != nil {
			h.cleanup()
		}
		if fn != nil {
			h.cleanup = fn()
		} else {
			h.cleanup = nil
		}
		h.lastDeps = cloneDeps(deps)
	}
}

// --- useRef ---

// Ref holds a mutable value that persists across renders.
type Ref[T any] struct {
	Value T
}

// refHook stores a Ref pointer that persists across renders.
type refHook[T any] struct {
	ref *Ref[T]
}

// UseRef returns a mutable ref object whose .Value field persists across renders.
// The returned pointer is stable across renders (same address on each call).
// The initial value is set on the first render and never resets.
//
//	inputRef := react.UseRef(ctx, "")
//	inputRef.Value = "new text"
func UseRef[T any](ctx *Context, initial T) *Ref[T] {
	if ctx.hookIndex >= len(ctx.hooks) {
		h := &refHook[T]{ref: &Ref[T]{Value: initial}}
		ctx.hooks = append(ctx.hooks, h)
		ctx.hookIndex++
		return h.ref
	}

	h, ok := ctx.hooks[ctx.hookIndex].(*refHook[T])
	if !ok {
		expected := reflect.TypeOf((*T)(nil)).Elem()
		got := reflect.TypeOf(ctx.hooks[ctx.hookIndex])
		panic(fmt.Sprintf(
			"bubbletea-react: hook type mismatch at index %d: expected %v, got %v. "+
				"Hooks must be called in the same order and with the same types on every render.",
			ctx.hookIndex, expected, got,
		))
	}
	ctx.hookIndex++
	return h.ref
}

// --- helpers ---

func depsChanged(old, new []any) bool {
	if len(old) != len(new) {
		return true
	}
	for i := range old {
		if !reflect.DeepEqual(old[i], new[i]) {
			return true
		}
	}
	return false
}

func cloneDeps(deps []any) []any {
	if deps == nil {
		return nil
	}
	c := make([]any, len(deps))
	copy(c, deps)
	return c
}
