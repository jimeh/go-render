package render

import (
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
)

var (
	// YAMLDefaultIndent is the default number of spaces to use for indenting
	// nested blocks in the output YAML.
	YAMLDefaultIndent = 2

	// YAMLDefaultIndentSequence is the default value for the IndentSequence
	// option passed to yaml.NewEncoder(). When true, sequences will be
	// indented by default.
	YAMLDefaultIndentSequence = true
)

// YAML is a Handler that marshals the given value to YAML.
type YAML struct {
	// Indent is the number of spaces to use for indenting nested blocks in the
	// output YAML. When empty/zero, YAMLDefaultIndent will be used.
	Indent int

	// EncodeOptions is a list of options to pass to yaml.NewEncoder().
	//
	// These options will be appended to the default indent and sequence indent
	// options. With duplicate options, the last one will take precedence.
	EncodeOptions []yaml.EncodeOption
}

var (
	_ Handler        = (*YAML)(nil)
	_ FormatsHandler = (*YAML)(nil)
)

func (y *YAML) buildEncodeOptions() []yaml.EncodeOption {
	indent := YAMLDefaultIndent
	if y.Indent > 0 {
		indent = y.Indent
	}

	opts := []yaml.EncodeOption{
		yaml.Indent(indent),
		yaml.IndentSequence(YAMLDefaultIndentSequence),
	}

	if len(y.EncodeOptions) > 0 {
		opts = append(opts, y.EncodeOptions...)
	}

	return opts
}

// Render marshals the given value to YAML.
func (y *YAML) Render(w io.Writer, v any) error {
	enc := yaml.NewEncoder(w, y.buildEncodeOptions()...)

	err := enc.Encode(v)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	err = enc.Close()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	return nil
}

// Formats returns a list of format strings that this Handler supports.
func (y *YAML) Formats() []string {
	return []string{"yaml", "yml"}
}
