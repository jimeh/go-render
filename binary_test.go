package render_test

import (
	"errors"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockBinaryMarshaler struct {
	data []byte
	err  error
}

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
			wantErr:   "render: marshal error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "error writing to writer",
			writeErr:  errors.New("write error!!1"),
			value:     &mockBinaryMarshaler{data: []byte("test string")},
			wantErr:   "render: write error!!1",
			wantErrIs: []error{render.Err},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &render.Binary{}

			var err error
			var got string
			w := &mockWriter{WriteErr: tt.writeErr}

			err = b.Render(w, tt.value)
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
