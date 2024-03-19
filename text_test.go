package render_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
)

type mockStringer struct {
	value string
}

var _ fmt.Stringer = (*mockStringer)(nil)

func (ms *mockStringer) String() string {
	return ms.value
}

type mockWriterTo struct {
	value string
	err   error
}

var _ io.WriterTo = (*mockWriterTo)(nil)

func (m *mockWriterTo) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(m.value))

	if m.err != nil {
		return int64(n), m.err
	}

	return int64(n), err
}

func TestText_Render(t *testing.T) {
	tests := []struct {
		name      string
		writeErr  error
		value     any
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name:  "implements fmt.Stringer",
			value: &mockStringer{value: "test string"},
			want:  "test string",
		},
		{
			name:      "error writing to writer with fmt.Stringer",
			writeErr:  errors.New("write error!!1"),
			value:     &mockStringer{value: "test string"},
			wantErr:   "render: failed: write error!!1",
			wantErrIs: []error{render.Err, render.ErrFailed},
		},
		{
			name:  "implements io.WriterTo",
			value: &mockWriterTo{value: "test string"},
			want:  "test string",
		},
		{
			name:      "does not implement fmt.Stringer or io.WriterTo",
			value:     struct{}{},
			wantErr:   "render: cannot render",
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		{
			name: "error writing to writer with io.WriterTo",
			value: &mockWriterTo{
				value: "test string",
				err:   errors.New("WriteTo error!!1"),
			},
			wantErr:   "render: failed: WriteTo error!!1",
			wantErrIs: []error{render.Err, render.ErrFailed},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &render.Text{}
			w := &mockWriter{WriteErr: tt.writeErr}

			err := s.Render(w, tt.value)
			got := w.String()

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
			for _, e := range tt.wantErrIs {
				assert.ErrorIs(t, err, e)
			}

			if tt.wantErr == "" && len(tt.wantErrIs) == 0 {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
