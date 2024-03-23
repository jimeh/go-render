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

var (
	// DefaultBinary is the default binary marshaler renderer. It
	// renders values using the encoding.BinaryMarshaler interface.
	DefaultBinary = &Binary{}

	// DefaultJSON is the default JSON renderer. It renders values using the
	// encoding/json package, with pretty printing enabled.
	DefaultJSON = &JSON{Pretty: true}

	// DefaultText is the default text renderer, used by the package level
	// Render function. It renders values using the DefaultStringer and
	// DefaultWriterTo renderers. This means a value must implement either the
	// fmt.Stringer or io.WriterTo interfaces to be rendered.
	DefaultText = &Text{}

	// DefaultXML is the default XML renderer. It renders values using the
	// encoding/xml package, with pretty printing enabled.
	DefaultXML = &XML{Pretty: true}

	// DefaultYAML is the default YAML renderer. It renders values using the
	// gopkg.in/yaml.v3 package, with an indentation of 2 spaces.
	DefaultYAML = &YAML{Indent: 2}

	// DefaultRenderer is used by the package level Render function. It supports
	// the text", "json", and "yaml" formats. If you need to support another set
	// of formats, use the New function to create a custom FormatRenderer.
	DefaultRenderer = MustNew("json", "text", "yaml")
)

// Render renders the given value to the given writer using the given format.
// The format must be one of the formats supported by the default renderer.
//
// By default it supports the following formats:
//
//   - "text": Renders values using the fmt.Stringer and io.WriterTo interfaces.
//   - "json": Renders values using the encoding/json package, with pretty
//     printing enabled.
//   - "yaml": Renders values using the gopkg.in/yaml.v3 package, with an
//     indentation of 2 spaces.
//
// If the format is not supported, a ErrUnsupportedFormat error will be
// returned.
//
// If you need to support a custom set of formats, use the New function to
// create a new FormatRenderer with the formats you need. If you need new custom
// renderers, manually create a new FormatRenderer.
func Render(w io.Writer, format string, v any) error {
	return DefaultRenderer.Render(w, format, v)
}

// New creates a new *FormatRenderer with support for the given formats.
//
// Supported formats are:
//
//   - "binary": Renders values using DefaultBinary.
//   - "json": Renders values using DefaultJSON.
//   - "text": Renders values using DefaultText.
//   - "xml": Renders values using DefaultXML.
//   - "yaml": Renders values using DefaultYAML.
//
// If an unsupported format is given, an ErrUnsupportedFormat error will be
// returned.
func New(formats ...string) (*Renderer, error) {
	renderers := map[string]FormatRenderer{}

	if len(formats) == 0 {
		return nil, fmt.Errorf("%w: no formats specified", Err)
	}

	for _, format := range formats {
		switch format {
		case "binary":
			renderers[format] = DefaultBinary
		case "json":
			renderers[format] = DefaultJSON
		case "text":
			renderers[format] = DefaultText
		case "xml":
			renderers[format] = DefaultXML
		case "yaml":
			renderers[format] = DefaultYAML
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
		}
	}

	return NewFormatRenderer(renderers), nil
}

// MustNew is like New, but panics if an error occurs.
func MustNew(formats ...string) *Renderer {
	r, err := New(formats...)
	if err != nil {
		panic(err.Error())
	}

	return r
}
