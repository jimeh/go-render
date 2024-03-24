// Package render provides a simple and flexible way to render a value to a
// io.Writer using different formats based on a format string argument.
//
// It allows rendering a custom type which can be marshaled to JSON, YAML, XML,
// while also supporting plain text by implementing fmt.Stringer or io.WriterTo.
// Binary output is also supported by implementing the encoding.BinaryMarshaler
// interface.
//
// Originally intended to easily implement CLI tools which can output their data
// as plain text, as well as JSON/YAML with a simple switch of a format string.
// But it can just as easily render to any io.Writer.
package render

import (
	"fmt"
	"io"
)

var (
	// Err is the base error for the package. All errors returned by this
	// package are wrapped with this error.
	Err       = fmt.Errorf("render")
	ErrFailed = fmt.Errorf("%w: failed", Err)

	// ErrCannotRender is returned when a value cannot be rendered. This may be
	// due to the value not supporting the format, or the value itself not being
	// renderable. Only Renderer implementations should return this error.
	ErrCannotRender = fmt.Errorf("%w: cannot render", Err)
)

// FormatRenderer interface is for single format renderers, which can only
// render a single format.
type FormatRenderer interface {
	// Render writes v into w in the format that the FormatRenderer supports.
	//
	// If v does not implement a required interface, or otherwise cannot be
	// rendered to the format in question, then a ErrCannotRender error must be
	// returned. Any other errors should be returned as is.
	Render(w io.Writer, v any) error
}

// Formats is an optional interface that can be implemented by FormatRenderer
// implementations to return a list of formats that the renderer supports. This
// is used by the NewRenderer function to allowing format aliases like "yml" for
// "yaml".
type Formats interface {
	Formats() []string
}

var (
	prettyRenderer = New(map[string]FormatRenderer{
		"binary": &Binary{},
		"json":   &JSON{Pretty: true},
		"text":   &Text{},
		"xml":    &XML{Pretty: true},
		"yaml":   &YAML{Indent: 2},
	})
	compactRenderer = New(map[string]FormatRenderer{
		"binary": &Binary{},
		"json":   &JSON{},
		"text":   &Text{},
		"xml":    &XML{},
		"yaml":   &YAML{},
	})

	DefaultPretty  = prettyRenderer.OnlyWith("json", "text", "xml", "yaml")
	DefaultCompact = compactRenderer.OnlyWith("json", "text", "xml", "yaml")
)

// Render renders the given value to the given writer using the given format.
// If pretty is true, the value will be rendered in a pretty way, otherwise it
// will be rendered in a compact way.
//
// By default it supports the following formats:
//
//   - "text": Renders values via a myriad of ways.
//   - "json": Renders values using the encoding/json package.
//   - "yaml": Renders values using the gopkg.in/yaml.v3 package.
//   - "xml": Renders values using the encoding/xml package.
//
// If the format is not supported, a ErrUnsupportedFormat error will be
// returned.
func Render(w io.Writer, format string, pretty bool, v any) error {
	if pretty {
		return DefaultPretty.Render(w, format, v)
	}

	return DefaultCompact.Render(w, format, v)
}

// Pretty renders the given value to the given writer using the given format.
// The format must be one of the formats supported by the default renderer.
//
// By default it supports the following formats:
//
//   - "text": Renders values via a myriad of ways.
//   - "json": Renders values using the encoding/json package, with pretty
//     printing enabled.
//   - "yaml": Renders values using the gopkg.in/yaml.v3 package, with an
//     indentation of 2 spaces.
//   - "xml": Renders values using the encoding/xml package, with pretty
//     printing enabled.
//
// If the format is not supported, a ErrUnsupportedFormat error will be
// returned.
//
// If you need to support a custom set of formats, use the New function to
// create a new Renderer with the formats you need. If you need new custom
// renderers, manually create a new Renderer.
func Pretty(w io.Writer, format string, v any) error {
	return DefaultPretty.Render(w, format, v)
}

// Compact renders the given value to the given writer using the given format.
// The format must be one of the formats supported by the default renderer.
//
// By default it supports the following formats:
//
//   - "text": Renders values via a myriad of ways..
//   - "json": Renders values using the encoding/json package.
//   - "yaml": Renders values using the gopkg.in/yaml.v3 package.
//   - "xml": Renders values using the encoding/xml package.
//
// If the format is not supported, a ErrUnsupportedFormat error will be
// returned.
//
// If you need to support a custom set of formats, use the New function to
// create a new Renderer with the formats you need. If you need new custom
// renderers, manually create a new Renderer.
func Compact(w io.Writer, format string, v any) error {
	return DefaultCompact.Render(w, format, v)
}

// NewCompact returns a new renderer which only supports the specified formats
// and renders structured formats compactly. If no formats are specified, a
// error is returned.
//
// If any of the formats are not supported by, a ErrUnsupported error is
// returned.
func NewCompact(formats ...string) (*Renderer, error) {
	if len(formats) == 0 {
		return nil, fmt.Errorf("%w: no formats specified", Err)
	}

	for _, format := range formats {
		if _, ok := compactRenderer.Renderers[format]; !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
		}
	}

	return compactRenderer.OnlyWith(formats...), nil
}

// NewPretty returns a new renderer which only supports the specified formats
// and renders structured formats in a pretty way. If no formats are specified,
// a error is returned.
//
// If any of the formats are not supported by, a ErrUnsupported error is
// returned.
func NewPretty(formats ...string) (*Renderer, error) {
	if len(formats) == 0 {
		return nil, fmt.Errorf("%w: no formats specified", Err)
	}

	for _, format := range formats {
		if _, ok := prettyRenderer.Renderers[format]; !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
		}
	}

	return prettyRenderer.OnlyWith(formats...), nil
}
