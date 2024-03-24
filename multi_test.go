package render

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiRenderer_Render(t *testing.T) {
	successRenderer := &mockRenderer{output: "success output"}
	cannotRenderer := &mockRenderer{err: ErrCannotRender}
	failRenderer := &mockRenderer{err: errors.New("mock error")}

	tests := []struct {
		name      string
		renderers []FormatRenderer
		value     interface{}
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "no renderer can render",
			renderers: []FormatRenderer{
				cannotRenderer,
				cannotRenderer,
			},
			value:     "test",
			wantErr:   "render: cannot render: string",
			wantErrIs: []error{ErrCannotRender},
		},
		{
			name: "one renderer can render",
			renderers: []FormatRenderer{
				cannotRenderer,
				successRenderer,
				cannotRenderer,
			},
			value: struct{}{},
			want:  "success output",
		},
		{
			name: "multiple renderers can render",
			renderers: []FormatRenderer{
				&mockRenderer{err: ErrCannotRender},
				&mockRenderer{output: "first output"},
				&mockRenderer{output: "second output"},
			},
			value: struct{}{},
			want:  "first output",
		},
		{
			name: "first renderer fails",
			renderers: []FormatRenderer{
				failRenderer,
				successRenderer,
			},
			value:   struct{}{},
			wantErr: "mock error",
		},
		{
			name: "fails after cannot render",
			renderers: []FormatRenderer{
				cannotRenderer,
				failRenderer,
				successRenderer,
			},
			value:   struct{}{},
			wantErr: "mock error",
		},
		{
			name: "fails after success render",
			renderers: []FormatRenderer{
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
			mr := &Multi{
				Renderers: tt.renderers,
			}
			var buf bytes.Buffer

			err := mr.Render(&buf, tt.value)
			got := buf.String()

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
