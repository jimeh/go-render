package render

import (
	"encoding"
	"fmt"
	"io"
)

// Binary can render values which implement the encoding.BinaryMarshaler
// interface.
type Binary struct{}

var (
	_ Handler        = (*Binary)(nil)
	_ FormatsHandler = (*Binary)(nil)
)

// Render writes result of calling MarshalBinary() on v. If v does not implement
// encoding.BinaryMarshaler the ErrCannotRender error will be returned.
func (br *Binary) Render(w io.Writer, v any) error {
	x, ok := v.(encoding.BinaryMarshaler)
	if !ok {
		return fmt.Errorf("%w: %T", ErrCannotRender, v)
	}

	b, err := x.MarshalBinary()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	return nil
}

// Formats returns a list of format strings that this Handler supports.
func (br *Binary) Formats() []string {
	return []string{"binary", "bin"}
}
