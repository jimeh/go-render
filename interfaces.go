package render

import "io"

// Handler interface is for single format renderers, which can only render a
// single format. It is the basis of the multi-format support offerred by the
// render package.
type Handler interface {
	// Render writes v into w in the format that the Handler supports.
	//
	// If v does not implement a required interface, or otherwise cannot be
	// rendered to the format in question, then a ErrCannotRender error must be
	// returned. Any other errors should be returned as is.
	Render(w io.Writer, v any) error
}

// PrettyHandler interface is a optional interface that can be implemented by
// Handler implementations to render a value in a pretty way. This is
// useful for formats that support pretty printing, like in the case of JSON and
// XML.
type PrettyHandler interface {
	// RenderPretty writes v into w in the format that the Handler supports,
	// using a pretty variant of the format. The exact definition of "pretty" is
	// up to the handler. Typically this would be mean adding line breaks and
	// indentation, like in the case of JSON and XML.
	//
	// If v does not implement a required interface, or otherwise cannot be
	// rendered to the format in question, then a ErrCannotRender error must be
	// returned. Any other errors should be returned as is.
	RenderPretty(w io.Writer, v any) error
}

// FormatsHandler is an optional interface that can be implemented by Handler
// implementations to return a list of formats that the handler supports. This
// is used by the New function to allow format aliases like "yml" for "yaml".
type FormatsHandler interface {
	// Formats returns a list of strings which all target the same format. In
	// most cases this would just be a single value, but multiple values are
	// supported for the sake of aliases, like "yaml" and "yml".
	Formats() []string
}
