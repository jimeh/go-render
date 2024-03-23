package render

import (
	"fmt"
	"io"
)

// Text is a renderer that writes the given value to the writer as is. It
// supports rendering the following types as plain text:
//
//   - []byte
//   - []rune
//   - string
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64
//   - float32, float64
//   - bool
//   - io.Reader
//   - io.WriterTo
//   - fmt.Stringer
//   - error
//
// If the value is of any other type, a ErrCannotRender error will be returned.
type Text struct{}

var _ FormatRenderer = (*Text)(nil)

// Render writes the given value to the writer as text.
func (t *Text) Render(w io.Writer, v any) error {
	var err error
	switch x := v.(type) {
	case []byte:
		_, err = w.Write(x)
	case []rune:
		_, err = w.Write([]byte(string(x)))
	case string:
		_, err = w.Write([]byte(x))
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, bool:
		_, err = fmt.Fprintf(w, "%v", x)
	case io.Reader:
		_, err = io.Copy(w, x)
	case io.WriterTo:
		_, err = x.WriteTo(w)
	case fmt.Stringer:
		_, err = w.Write([]byte(x.String()))
	case error:
		_, err = w.Write([]byte(x.Error()))
	default:
		return fmt.Errorf("%w: %T", ErrCannotRender, v)
	}

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailed, err)
	}

	return nil
}
