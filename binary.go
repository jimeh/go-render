package render

import (
	"encoding"
	"fmt"
	"io"
)

// Binary can render values which implment the encoding.BinaryMarshaler
// interface.
type Binary struct{}

var _ FormatRenderer = (*Binary)(nil)

// Render writes result of calling MarshalBinary() on v. If v does not implment
// encoding.BinaryMarshaler the ErrCannotRander error will be returned.
func (bm *Binary) Render(w io.Writer, v any) error {
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
