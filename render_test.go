package render_test

import (
	"bytes"
	"io"
)

type mockWriter struct {
	WriteErr error
	buf      bytes.Buffer
}

func (mw *mockWriter) Write(p []byte) (n int, err error) {
	if mw.WriteErr != nil {
		return 0, mw.WriteErr
	}

	return mw.buf.Write(p)
}

func (mw *mockWriter) String() string {
	return mw.buf.String()
}

type mockRenderer struct {
	output string
	err    error
}

func (m *mockRenderer) Render(w io.Writer, _ any) error {
	_, err := w.Write([]byte(m.output))

	if m.err != nil {
		return m.err
	}

	return err
}
