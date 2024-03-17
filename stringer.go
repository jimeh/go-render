package render

import (
	"fmt"
	"io"
)

// Stringer is a renderer that renders a value to an io.Writer using the
// String method.
type Stringer struct{}

var _ Renderer = (*Stringer)(nil)

// Render renders a value to an io.Writer using the String method. If the value
// does not implement fmt.Stringer, ErrCannotRender is returned.
func (s *Stringer) Render(w io.Writer, v any) error {
	x, ok := v.(fmt.Stringer)
	if !ok {
		return ErrCannotRender
	}

	_, err := fmt.Fprint(w, x.String())
	if err != nil {
		return fmt.Errorf("%w: %w", Err, err)
	}

	return nil
}
