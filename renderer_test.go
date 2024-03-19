package render_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderer_Render(t *testing.T) {
	tests := []struct {
		name      string
		renderers map[string]render.FormatRenderer
		format    string
		value     interface{}
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "existing renderer",
			renderers: map[string]render.FormatRenderer{
				"mock": &mockRenderer{output: "mock output"},
			},
			format: "mock",
			value:  struct{}{},
			want:   "mock output",
		},
		{
			name: "existing renderer returns error",
			renderers: map[string]render.FormatRenderer{
				"other": &mockRenderer{
					output: "mock output",
					err:    errors.New("mock error"),
				},
			},
			format:  "other",
			value:   struct{}{},
			wantErr: "render: failed: mock error",
		},
		{
			name: "existing renderer returns ErrCannotRender",
			renderers: map[string]render.FormatRenderer{
				"other": &mockRenderer{
					output: "mock output",
					err:    fmt.Errorf("%w: mock", render.ErrCannotRender),
				},
			},
			format:    "other",
			value:     struct{}{},
			wantErr:   "render: unsupported format: other",
			wantErrIs: []error{render.Err, render.ErrUnsupportedFormat},
		},
		{
			name:      "non-existing renderer",
			renderers: map[string]render.FormatRenderer{},
			format:    "unknown",
			value:     struct{}{},
			wantErr:   "render: unsupported format: unknown",
			wantErrIs: []error{render.Err, render.ErrUnsupportedFormat},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &render.Renderer{
				Formats: tt.renderers,
			}
			var buf bytes.Buffer

			err := fr.Render(&buf, tt.format, tt.value)
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

func TestRenderer_RenderAllFormats(t *testing.T) {
	tests := []renderFormatTestCase{}
	tests = append(tests, binaryFormattestCases...)
	tests = append(tests, jsonFormatTestCases...)
	tests = append(tests, textFormatTestCases...)
	tests = append(tests, xmlFormatTestCases...)
	tests = append(tests, yamlFormatTestCases...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &mockWriter{WriteErr: tt.writeErr}

			var err error
			var panicRes any
			renderer, err := render.New("binary", "json", "text", "xml", "yaml")
			require.NoError(t, err)

			func() {
				defer func() {
					if r := recover(); r != nil {
						panicRes = r
					}
				}()
				err = renderer.Render(w, tt.format, tt.value)
			}()

			got := w.String()

			if tt.wantPanic != "" {
				assert.Equal(t, tt.wantPanic, panicRes)
			}

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
			for _, e := range tt.wantErrIs {
				assert.ErrorIs(t, err, e)
			}

			if tt.wantPanic == "" &&
				tt.wantErr == "" && len(tt.wantErrIs) == 0 {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
