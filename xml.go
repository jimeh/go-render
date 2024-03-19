package render

import (
	"encoding/xml"
	"fmt"
	"io"
)

// XML is a Renderer that marshals a value to XML.
type XML struct {
	// Pretty specifies whether the output should be pretty-printed. If true,
	// the output will be indented and newlines will be added.
	Pretty bool

	// Prefix is the prefix added to each level of indentation when Pretty is
	// true.
	Prefix string

	// Indent is the string added to each level of indentation when Pretty is
	// true. If empty, two spaces will be used instead.
	Indent string
}

var _ FormatRenderer = (*XML)(nil)

// Render marshals the given value to XML.
func (x *XML) Render(w io.Writer, v any) error {
	enc := xml.NewEncoder(w)
	if x.Pretty {
		prefix := x.Prefix
		indent := x.Indent
		if indent == "" {
			indent = "  "
		}

		enc.Indent(prefix, indent)
	}

	err := enc.Encode(v)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	return nil
}
