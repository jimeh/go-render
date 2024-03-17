package render

import (
	"errors"
	"io"
)

// MultiRenderer is a renderer that tries multiple renderers until one succeeds.
type MultiRenderer struct {
	Renderers []Renderer
}

var _ Renderer = (*MultiRenderer)(nil)

// Render tries each renderer in order until one succeeds. If none succeed,
// ErrCannotRender is returned. If a renderer returns an error that is not
// ErrCannotRender, that error is returned.
func (mr *MultiRenderer) Render(w io.Writer, v any) error {
	for _, r := range mr.Renderers {
		err := r.Render(w, v)
		if err == nil {
			return nil
		}
		if !errors.Is(err, ErrCannotRender) {
			return err
		}
	}

	return ErrCannotRender
}
