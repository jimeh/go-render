package render

import (
	"encoding/xml"
	"fmt"
	"io"
)

// XMLDefaultIndent is the default indentation string used by XML instances when
// pretty rendering if no Indent value is set.
var XMLDefaultIndent = "  "

// XML is a Renderer that marshals a value to XML.
type XML struct {
	// Prefix is the prefix added to each level of indentation when pretty
	// rendering.
	Prefix string

	// Indent is the string added to each level of indentation when pretty
	// rendering. If empty, XMLDefaultIndent will be used.
	Indent string
}

var (
	_ Handler        = (*XML)(nil)
	_ PrettyHandler  = (*XML)(nil)
	_ FormatsHandler = (*XML)(nil)
)

// Render marshals the given value to XML.
func (x *XML) Render(w io.Writer, v any) error {
	err := xml.NewEncoder(w).Encode(v)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	return nil
}

// RenderPretty marshals the given value to XML with line breaks and
// indentation.
func (x *XML) RenderPretty(w io.Writer, v any) error {
	prefix := x.Prefix
	indent := x.Indent
	if indent == "" {
		indent = XMLDefaultIndent
	}

	enc := xml.NewEncoder(w)
	enc.Indent(prefix, indent)

	err := enc.Encode(v)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	return nil
}

// Formats returns a list of format strings that this Handler supports.
func (x *XML) Formats() []string {
	return []string{"xml"}
}
