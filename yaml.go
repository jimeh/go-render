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

	// YAMLDefaultEncoderOptions is the default list of options passed to
	// yaml.NewEncoder(). Any options specified in YAML.EncodeOptions will be
	// passed after these defaults.
	YAMLDefaultEncoderOptions = []yaml.EncodeOption{
		yaml.Indent(YAMLDefaultIndent),
		yaml.IndentSequence(true),
	}
)

// YAML is a Handler that marshals the given value to YAML.
type YAML struct {
	// EncodeOptions is a list of options to pass to yaml.NewEncoder().
	// If empty, YAMLDefaultEncoderOptions will be used.
	EncodeOptions []yaml.EncodeOption

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
	opts := YAMLDefaultEncoderOptions

	if y.Indent > 0 {
		opts = append(opts, yaml.Indent(y.Indent))
	}
	if len(y.EncodeOptions) > 0 {
		opts = append(opts, y.EncodeOptions...)
	}

	enc := yaml.NewEncoder(w, opts...)

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
