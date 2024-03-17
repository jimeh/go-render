package render

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// YAML is a renderer that marshals the given value to YAML.
type YAML struct {
	// Indent controls how many spaces will be used for indenting nested blocks
	// in the output YAML. When Indent is zero, 2 will be used by default.
	Indent int
}

var _ Renderer = (*YAML)(nil)

// Render marshals the given value to YAML.
func (j *YAML) Render(w io.Writer, v any) error {
	enc := yaml.NewEncoder(w)

	indent := j.Indent
	if indent == 0 {
		indent = 2
	}

	enc.SetIndent(indent)

	err := enc.Encode(v)
	if err != nil {
		return fmt.Errorf("%w: %w", Err, err)
	}

	return nil
}
