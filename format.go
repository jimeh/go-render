package render

import (
	"fmt"
	"io"
)

var ErrUnsupportedFormat = fmt.Errorf("%w: unsupported format", Err)

// FormatRenderer is a renderer that delegates rendering to another renderer
// based on a format value.
type FormatRenderer struct {
	// Renderers is a map of format names to renderers. When Render is called,
	// the format is used to look up the renderer to use.
	Renderers map[string]Renderer
}

// Render renders a value to an io.Writer using the specified format. If the
// format is not supported, ErrUnsupportedFormat is returned.
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

	return ErrUnsupportedFormat
}
