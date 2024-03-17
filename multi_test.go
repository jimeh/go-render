package render_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiRenderer_Render(t *testing.T) {
	successRenderer := &mockRenderer{output: "success output"}
	cannotRenderer := &mockRenderer{err: render.ErrCannotRender}
	failRenderer := &mockRenderer{err: errors.New("mock error")}

	tests := []struct {
		name      string
		renderers []render.Renderer
		value     interface{}
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "no renderer can render",
			renderers: []render.Renderer{
				cannotRenderer,
				cannotRenderer,
			},
			value:     struct{}{},
			wantErrIs: []error{render.ErrCannotRender},
		},
		{
			name: "one renderer can render",
			renderers: []render.Renderer{
				cannotRenderer,
				successRenderer,
				cannotRenderer,
			},
			value: struct{}{},
			want:  "success output",
		},
		{
			name: "multiple renderers can render",
			renderers: []render.Renderer{
				&mockRenderer{err: render.ErrCannotRender},
				&mockRenderer{output: "first output"},
				&mockRenderer{output: "second output"},
			},
			value: struct{}{},
			want:  "first output",
		},
		{
			name: "first renderer fails",
			renderers: []render.Renderer{
				failRenderer,
				successRenderer,
			},
			value:   struct{}{},
			wantErr: "mock error",
		},
		{
			name: "fails after cannot render",
			renderers: []render.Renderer{
				cannotRenderer,
				failRenderer,
				successRenderer,
			},
			value:   struct{}{},
			wantErr: "mock error",
		},
		{
			name: "fails after success render",
			renderers: []render.Renderer{
				successRenderer,
				failRenderer,
				cannotRenderer,
			},
			value: struct{}{},
			want:  "success output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &render.MultiRenderer{
				Renderers: tt.renderers,
			}

			var buf bytes.Buffer
			err := mr.Render(&buf, tt.value)
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
