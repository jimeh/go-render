package render_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatRenderer_Render(t *testing.T) {
	tests := []struct {
		name      string
		renderers map[string]render.Renderer
		format    string
		value     interface{}
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "existing renderer",
			renderers: map[string]render.Renderer{
				"mock": &mockRenderer{output: "mock output"},
			},
			format: "mock",
			value:  struct{}{},
			want:   "mock output",
		},
		{
			name: "existing renderer returns error",
			renderers: map[string]render.Renderer{
				"other": &mockRenderer{
					output: "mock output",
					err:    errors.New("mock error"),
				},
			},
			format:  "other",
			value:   struct{}{},
			wantErr: "mock error",
		},
		{
			name:      "non-existing renderer",
			renderers: map[string]render.Renderer{},
			format:    "unknown",
			value:     struct{}{},
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &render.FormatRenderer{
				Renderers: tt.renderers,
			}

			var buf bytes.Buffer
			err := fr.Render(&buf, tt.format, tt.value)
			got := buf.String()

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
