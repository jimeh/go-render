package render

import (
	"errors"
	"fmt"
	"io"
)

// ErrUnsupportedFormat is returned when a format is not supported by a
// renderer. Any method that accepts a format string may return this error.
var ErrUnsupportedFormat = fmt.Errorf("%w: unsupported format", Err)

// Renderer is a renderer that delegates rendering to another renderer
// based on a format value.
type Renderer struct {
	// Renderers is a map of format names to renderers. When Render is called,
	// the format is used to look up the renderer to use.
	Renderers map[string]FormatRenderer
}

// New returns a new Renderer that delegates rendering to the specified
// renderers.
//
// Any renderers which implement the Formats interface, will also be set as the
// renderer for all format strings returned by Format() on the renderer.
func New(renderers map[string]FormatRenderer) *Renderer {
	newRenderers := make(map[string]FormatRenderer, len(renderers))

	for format, r := range renderers {
		newRenderers[format] = r

		if x, ok := r.(Formats); ok {
			for _, f := range x.Formats() {
				if f != format {
					newRenderers[f] = r
				}
			}
		}
	}

	return &Renderer{Renderers: newRenderers}
}

// Render renders a value to an io.Writer using the specified format. If the
// format is not supported, ErrUnsupportedFormat is returned.
//
// If the format is supported, but the value cannot be rendered to the format,
// the error returned by the renderer is returned. In most cases this will be
// ErrCannotRender, but it could be a different error if the renderer returns
// one.
func (r *Renderer) Render(w io.Writer, format string, v any) error {
	renderer, ok := r.Renderers[format]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
	}

	err := renderer.Render(w, v)
	if err != nil {
		if errors.Is(err, ErrCannotRender) {
			return fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
		}
		if !errors.Is(err, ErrFailed) {
			return fmt.Errorf("%w: %w", ErrFailed, err)
		}

		return err
	}

	return nil
}

func (r *Renderer) OnlyWith(formats ...string) *Renderer {
	renderers := make(map[string]FormatRenderer, len(formats))

	for _, format := range formats {
		if r, ok := r.Renderers[format]; ok {
			renderers[format] = r
		}
	}

	return New(renderers)
}
