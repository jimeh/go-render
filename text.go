package render

import (
	"fmt"
	"io"
)

type Text struct{}

var _ FormatRenderer = (*Text)(nil)

func (t *Text) Render(w io.Writer, v any) error {
	var err error
	switch x := v.(type) {
	case fmt.Stringer:
		_, err = w.Write([]byte(x.String()))
	case io.WriterTo:
		_, err = x.WriteTo(w)
	default:
		return ErrCannotRender
	}

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	return nil
}
