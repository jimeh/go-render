package render_test

import (
	"errors"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStringer struct {
	value string
}

func (ms *mockStringer) String() string {
	return ms.value
}

func TestStringer_Render(t *testing.T) {
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
			name:      "does not implement fmt.Stringer",
			value:     struct{}{},
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		{
			name:      "error writing to writer",
			writeErr:  errors.New("write error!!1"),
			value:     &mockStringer{value: "test string"},
			wantErr:   "render: write error!!1",
			wantErrIs: []error{render.Err},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &render.Stringer{}

			var err error
			var got string
			w := &mockWriter{WriteErr: tt.writeErr}

			err = s.Render(w, tt.value)
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
