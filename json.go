package render

import (
	"encoding/json"
	"fmt"
	"io"
)

// JSON is a renderer that marshals values to JSON.
type JSON struct {
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

var _ FormatRenderer = (*JSON)(nil)

// Render marshals the given value to JSON.
func (j *JSON) Render(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	if j.Pretty {
		prefix := j.Prefix
		indent := j.Indent
		if indent == "" {
			indent = "  "
		}

		enc.SetIndent(prefix, indent)
	}

	err := enc.Encode(v)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	return nil
}
