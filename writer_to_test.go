package render_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockWriterTo struct {
	value string
	err   error
}

func (m *mockWriterTo) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(m.value))

	if m.err != nil {
		return int64(n), m.err
	}

	return int64(n), err
}

func TestWriterTo_Render(t *testing.T) {
	tests := []struct {
		name      string
		writeErr  error
		value     any
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name:  "implements io.WriterTo",
			value: &mockWriterTo{value: "test string"},
			want:  "test string",
		},
		{
			name:      "does not implement io.WriterTo",
			value:     struct{}{},
			wantErr:   "render: cannot render",
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		{
			name: "error writing to writer",
			value: &mockWriterTo{
				value: "test string",
				err:   errors.New("WriteTo error!!1"),
			},
			wantErr:   "render: WriteTo error!!1",
			wantErrIs: []error{render.Err},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wt := &render.WriterTo{}

			var err error
			var got string
			w := &bytes.Buffer{}

			err = wt.Render(w, tt.value)
			got = w.String()

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErr)
			}
			for _, e := range tt.wantErrIs {
				assert.ErrorIs(t, err, e)
			}

			if tt.wantErr == "" && len(tt.wantErrIs) == 0 {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
