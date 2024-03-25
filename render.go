// Package render provides a simple and flexible way to render a value to a
// io.Writer using different formats based on a format string argument.
//
// It is designed around using a custom type/struct to render your output.
// Thanks to Go's marshaling interfaces, you get JSON, YAML, and XML support
// almost for free. While plain text output is supported by the type
// implementing io.Reader, io.WriterTo, fmt.Stringer, or error interfaces, or by
// simply being a type which can easily be type cast to a byte slice.
//
// Originally intended to easily implement CLI tools which can output their data
// as plain text, as well as JSON/YAML with a simple switch of a format string.
// But it can just as easily render to any io.Writer.
//
// The package is designed to be flexible and extensible with a sensible set of
// defaults accessible via package level functions. You can create your own
// Renderer for custom formats, or create new handlers that support custom
// formats.
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

	// Base is a renderer that supports all formats. It is used by the package
	// level NewWith function to create new renderers with a sub-set of
	// formats.
	Base = New(map[string]Handler{
		"binary": &Binary{},
		"json":   &JSON{},
		"text":   &Text{},
		"xml":    &XML{},
		"yaml":   &YAML{},
	})

	// Default is the default renderer that is used by package level Render,
	// Compact, Pretty functions. It supports JSON, Text, and YAML formats.
	Default = Base.NewWith("json", "text", "yaml")
)

// Render renders the given value to the given writer using the given format. If
// pretty is true, the value will be rendered "pretty" if the target format
// supports it, otherwise it will be rendered in a compact way.
//
// It uses the default renderer to render the value, which supports JSON, Text,
// and YAML formats out of the box.
//
// If you need to support a custom set of formats, use the New function to
// create a new Renderer with the formats you need. If you need new custom
// renderers, manually create a new Renderer.
func Render(w io.Writer, format string, pretty bool, v any) error {
	return Default.Render(w, format, pretty, v)
}

// Compact is a convenience function that calls the Default renderer's Compact
// method. It is the same as calling Render with pretty set to false.
func Compact(w io.Writer, format string, v any) error {
	return Default.Compact(w, format, v)
}

// Pretty is a convenience function that calls the Default renderer's Pretty
// method. It is the same as calling Render with pretty set to true.
func Pretty(w io.Writer, format string, v any) error {
	return Default.Pretty(w, format, v)
}

// NewWith creates a new Renderer with the given formats. Only formats on the
// BaseRender will be supported.
func NewWith(formats ...string) *Renderer {
	return Base.NewWith(formats...)
}
