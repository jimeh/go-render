package render

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrUnsupportedFormat is returned when a format is not supported by any
// Handler.
var ErrUnsupportedFormat = fmt.Errorf("%w: unsupported format", Err)

// Renderer exposes methods for rendering values to different formats. The
// Renderer delegates rendering to format specific handlers based on the format
// string given.
type Renderer struct {
	// Handlers is a map of format names to Handler. When Render is called,
	// the format is used to look up the Handler to use.
	Handlers map[string]Handler
}

// New returns a new Renderer that delegates rendering to the specified
// Handlers.
//
// Any Handlers which implement the FormatsHandler interface, will also be set
// as the handler for all format strings returned by Formats() on the handler.
func New(handlers map[string]Handler) *Renderer {
	r := &Renderer{Handlers: make(map[string]Handler, len(handlers))}

	for format, handler := range handlers {
		r.Add(format, handler)
	}

	return r
}

// Add adds a Handler to the Renderer. If the handler implements the
// FormatsHandler interface, the handler will be added for all formats returned
// by Formats().
func (r *Renderer) Add(format string, handler Handler) {
	if format != "" {
		r.Handlers[strings.ToLower(format)] = handler
	}

	if x, ok := handler.(FormatsHandler); ok {
		for _, f := range x.Formats() {
			if f != "" && f != format {
				r.Handlers[strings.ToLower(f)] = handler
			}
		}
	}
}

// Render renders a value to the given io.Writer using the specified format.
//
// If pretty is true, it will attempt to render the value with pretty
// formatting if the underlying Handler supports pretty formatting.
//
// If the format is not supported or the value cannot be rendered to the format,
// a ErrUnsupportedFormat error is returned.
func (r *Renderer) Render(
	w io.Writer,
	format string,
	pretty bool,
	v any,
) error {
	handler, ok := r.Handlers[strings.ToLower(format)]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
	}

	prettyHandler, ok := handler.(PrettyHandler)
	var err error
	if pretty && ok {
		err = prettyHandler.RenderPretty(w, v)
	} else {
		err = handler.Render(w, v)
	}

	if err != nil {
		if errors.Is(err, ErrCannotRender) {
			return fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
		}

		// Ensure that the error is wrapped with ErrFailed if it is not already.
		if !errors.Is(err, ErrFailed) {
			return fmt.Errorf("%w: %w", ErrFailed, err)
		}

		return err
	}

	return nil
}

// Compact is a convenience method that calls Render with pretty set to false.
func (r *Renderer) Compact(w io.Writer, format string, v any) error {
	return r.Render(w, format, false, v)
}

// Pretty is a convenience method that calls Render with pretty set to true.
func (r *Renderer) Pretty(w io.Writer, format string, v any) error {
	return r.Render(w, format, true, v)
}

// NewWith creates a new Renderer with the formats given, if they have handlers
// in the current Renderer. It essentially allows to restrict a Renderer to a
// only a sub-set of supported formats.
func (r *Renderer) NewWith(formats ...string) *Renderer {
	handlers := make(map[string]Handler, len(formats))

	for _, format := range formats {
		if r, ok := r.Handlers[strings.ToLower(format)]; ok {
			handlers[format] = r
		}
	}

	return New(handlers)
}
