package render

import (
	"errors"
	"fmt"
	"io"
)

// Multi is a renderer that tries multiple renderers until one succeeds.
type Multi struct {
	Renderers []FormatRenderer
}

var _ FormatRenderer = (*Multi)(nil)

// Render tries each renderer in order until one succeeds. If none succeed,
// ErrCannotRender is returned. If a renderer returns an error that is not
// ErrCannotRender, that error is returned.
func (mr *Multi) Render(w io.Writer, v any) error {
	for _, r := range mr.Renderers {
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

func (mr *Multi) Formats() []string {
	formats := make(map[string]struct{})

	for _, r := range mr.Renderers {
		if x, ok := r.(Formats); ok {
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
