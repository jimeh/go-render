package render

import (
	"io"
)

// FormatRenderer is a renderer that delegates rendering to another renderer
// based on a format value.
type FormatRenderer struct {
	Renderers map[string]Renderer
}

// Render renders a value to an io.Writer using the specified format. If the
// format is not supported, ErrCannotRender is returned.
//
// If the format is supported, but the value cannot be rendered to the format,
// the error returned by the renderer is returned. In most cases this will be
// ErrCannotRender, but it could be a different error if the renderer returns
// one.
func (r *FormatRenderer) Render(w io.Writer, format string, v any) error {
	renderer, ok := r.Renderers[format]
	if ok {
		return renderer.Render(w, v)
	}

	return ErrCannotRender
}
