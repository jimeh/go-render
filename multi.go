package render

import (
	"errors"
	"fmt"
	"io"
)

// Multi is a Handler that tries multiple handlers until one succeeds.
type Multi struct {
	Handlers []Handler
}

var (
	_ Handler        = (*Multi)(nil)
	_ PrettyHandler  = (*Multi)(nil)
	_ FormatsHandler = (*Multi)(nil)
)

// Render tries each handler in order until one succeeds. If none succeed,
// ErrCannotRender is returned. If a handler returns an error that is not
// ErrCannotRender, that error is returned.
func (mr *Multi) Render(w io.Writer, v any) error {
	for _, r := range mr.Handlers {
		err := r.Render(w, v)
		if err == nil {
			return nil
		}
		if !errors.Is(err, ErrCannotRender) {
			return err
		}
	}

	return fmt.Errorf("%w: %T", ErrCannotRender, v)
}

// RenderPretty tries each handler in order until one succeeds. If none
// succeed, ErrCannotRender is returned. If a handler returns an error that is
// not ErrCannotRender, that error is returned.
//
// If a handler implements PrettyHandler, then the RenderPretty method is used
// instead of Render. Otherwise, the Render method is used.
func (mr *Multi) RenderPretty(w io.Writer, v any) error {
	for _, r := range mr.Handlers {
		var err error
		if x, ok := r.(PrettyHandler); ok {
			err = x.RenderPretty(w, v)
		} else {
			err = r.Render(w, v)
		}
		if err == nil {
			return nil
		}
		if !errors.Is(err, ErrCannotRender) {
			return err
		}
	}

	return fmt.Errorf("%w: %T", ErrCannotRender, v)
}

// Formats returns a list of format strings that this Handler supports.
func (mr *Multi) Formats() []string {
	formats := make(map[string]struct{})

	for _, r := range mr.Handlers {
		if x, ok := r.(FormatsHandler); ok {
			for _, f := range x.Formats() {
				formats[f] = struct{}{}
			}
		}
	}

	result := make([]string, 0, len(formats))
	for f := range formats {
		result = append(result, f)
	}

	return result
}
