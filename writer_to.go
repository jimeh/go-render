package render

import (
	"fmt"
	"io"
)

// WriterTo is a renderer that renders a value to an io.Writer using the
// WriteTo method.
type WriterTo struct{}

var _ Renderer = (*WriterTo)(nil)

// Render renders a value to an io.Writer using the WriteTo method. If the value
// does not implement io.WriterTo, ErrCannotRender is returned.
func (wt *WriterTo) Render(w io.Writer, v any) error {
	x, ok := v.(io.WriterTo)
	if !ok {
		return ErrCannotRender
	}

	_, err := x.WriteTo(w)
	if err != nil {
		return fmt.Errorf("%w: %w", Err, err)
	}

	return nil
}
