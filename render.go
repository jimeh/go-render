package render

import (
	"fmt"
	"io"
)

var (
	Err = fmt.Errorf("render")

	// ErrCannotRender is returned when a value cannot be rendered. This may be
	// due to the value not supporting the format, or the value itself not being
	// renderable.
	ErrCannotRender = fmt.Errorf("%w: cannot render", Err)
)

// Renderer is the interface that that individual renderers must implement.
type Renderer interface {
	Render(w io.Writer, v any) error
}

var (
	// DefaultJSON is the default JSON renderer. It renders values using the
	// encoding/json package, with pretty printing enabled.
	DefaultJSON = &JSON{Pretty: true}

	// DefaultXML is the default XML renderer. It renders values using the
	// encoding/xml package, with pretty printing enabled.
	DefaultXML = &XML{Pretty: true}

	// DefaultYAML is the default YAML renderer. It renders values using the
	// gopkg.in/yaml.v3 package, with an indentation of 2 spaces.
	DefaultYAML = &YAML{Indent: 2}

	// DefaultWriterTo is the default writer to renderer. It renders values
	// using the io.WriterTo interface.
	DefaultWriterTo = &WriterTo{}

	// DefaultStringer is the default stringer renderer. It renders values
	// using the fmt.Stringer interface.
	DefaultStringer = &Stringer{}

	// DefaultText is the default text renderer, used by the package level
	// Render function. It renders values using the DefaultStringer and
	// DefaultWriterTo renderers. This means a value must implement either the
	// fmt.Stringer or io.WriterTo interfaces to be rendered.
	DefaultText = &MultiRenderer{
		Renderers: []Renderer{DefaultStringer, DefaultWriterTo},
	}

	// DefaultBinaryMarshaler is the default binary marshaler renderer. It
	// renders values using the encoding.BinaryMarshaler interface.
	DefaultBinaryMarshaler = &Binary{}

	// DefaultRenderer is the default renderer, used by the package level Render
	// function. It supports the "json", "xml", "yaml", "text", "binary"
	// formats.
	DefaultRenderer = &FormatRenderer{map[string]Renderer{
		"json":   DefaultJSON,
		"xml":    DefaultXML,
		"yaml":   DefaultYAML,
		"yml":    DefaultYAML,
		"text":   DefaultText,
		"binary": DefaultBinaryMarshaler,
		"bin":    DefaultBinaryMarshaler,
	}}
)

// Render renders the given value to the given writer using the given format.
// The format must be one of the formats supported by the default renderer.
//
// By default it supports the following formats:
//
//   - "json": Renders values using the encoding/json package, with pretty
//     printing enabled.
//   - "yaml": Renders values using the gopkg.in/yaml.v3 package, with an
//     indentation of 2 spaces.
//   - "yml": Alias for "yaml".
//   - "xml": Renders values using the encoding/xml package, with pretty
//     printing enabled.
//   - "text": Renders values using the fmt.Stringer and io.WriterTo interfaces.
//     This means a value must implement either the fmt.Stringer or io.WriterTo
//     interfaces to be rendered.
//   - "binary": Renders values using the encoding.BinaryMarshaler interface.
//   - "bin": Alias for "binary".
//
// If the format is not supported, a ErrCannotRender error will be returned.
func Render(w io.Writer, format string, v any) error {
	return DefaultRenderer.Render(w, format, v)
}
