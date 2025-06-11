package render

import (
	"encoding/json"
	"fmt"
	"io"
)

// JSONDefaultIndent is the default indentation string used by JSON instances
// when pretty rendering if no Indent value is set on the JSON instance itself.
var JSONDefaultIndent = "  "

// JSON is a Handler that marshals values to JSON.
type JSON struct {
	// Prefix is the prefix added to each level of indentation when pretty
	// rendering.
	Prefix string

	// Indent is the string added to each level of indentation when pretty
	// rendering. If empty, JSONDefaultIndent will be used.
	Indent string
}

var (
	_ Handler        = (*JSON)(nil)
	_ PrettyHandler  = (*JSON)(nil)
	_ FormatsHandler = (*JSON)(nil)
)

// Render marshals the given value to JSON.
func (jr *JSON) Render(w io.Writer, v any) error {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	return nil
}

// RenderPretty marshals the given value to JSON with line breaks and
// indentation.
func (jr *JSON) RenderPretty(w io.Writer, v any) error {
	prefix := jr.Prefix
	indent := jr.Indent
	if indent == "" {
		indent = JSONDefaultIndent
	}

	enc := json.NewEncoder(w)
	enc.SetIndent(prefix, indent)

	err := enc.Encode(v)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	return nil
}

// Formats returns a list of format strings that this Handler supports.
func (jr *JSON) Formats() []string {
	return []string{"json"}
}
