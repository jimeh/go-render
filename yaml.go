package render

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

var YAMLDefaultIndent = 2

// YAML is a Handler that marshals the given value to YAML.
type YAML struct {
	// Indent controls how many spaces will be used for indenting nested blocks
	// in the output YAML. When Indent is zero, YAMLDefaultIndent will be used.
	Indent int
}

var (
	_ Handler        = (*YAML)(nil)
	_ FormatsHandler = (*YAML)(nil)
)

// Render marshals the given value to YAML.
func (y *YAML) Render(w io.Writer, v any) error {
	indent := y.Indent
	if indent == 0 {
		indent = YAMLDefaultIndent
	}

	enc := yaml.NewEncoder(w)
	enc.SetIndent(indent)

	err := enc.Encode(v)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	return nil
}

// Formats returns a list of format strings that this Handler supports.
func (y *YAML) Formats() []string {
	return []string{"yaml", "yml"}
}
