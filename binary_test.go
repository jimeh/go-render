package render_test

import (
	"encoding"
	"errors"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
)

type mockBinaryMarshaler struct {
	data []byte
	err  error
}

var _ encoding.BinaryMarshaler = (*mockBinaryMarshaler)(nil)

func (mbm *mockBinaryMarshaler) MarshalBinary() ([]byte, error) {
	return mbm.data, mbm.err
}

func TestBinary_Render(t *testing.T) {
	tests := []struct {
		name      string
		writeErr  error
		value     any
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name:  "implements encoding.BinaryMarshaler",
			value: &mockBinaryMarshaler{data: []byte("test string")},
			want:  "test string",
		},
		{
			name:      "does not implement encoding.BinaryMarshaler",
			value:     struct{}{},
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		{
			name: "error marshaling",
			value: &mockBinaryMarshaler{
				data: []byte("test string"),
				err:  errors.New("marshal error!!1"),
			},
			wantErr:   "render: failed: marshal error!!1",
			wantErrIs: []error{render.Err, render.ErrFailed},
		},
		{
			name:      "error writing to writer",
			writeErr:  errors.New("write error!!1"),
			value:     &mockBinaryMarshaler{data: []byte("test string")},
			wantErr:   "render: failed: write error!!1",
			wantErrIs: []error{render.Err, render.ErrFailed},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &render.Binary{}
			w := &mockWriter{WriteErr: tt.writeErr}

			err := b.Render(w, tt.value)
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
