package render

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderer_Render(t *testing.T) {
	tests := []struct {
		name      string
		renderers map[string]FormatRenderer
		format    string
		value     interface{}
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "existing renderer",
			renderers: map[string]FormatRenderer{
				"mock": &mockRenderer{output: "mock output"},
			},
			format: "mock",
			value:  struct{}{},
			want:   "mock output",
		},
		{
			name: "existing renderer returns error",
			renderers: map[string]FormatRenderer{
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
			renderers: map[string]FormatRenderer{
				"other": &mockRenderer{
					output: "mock output",
					err:    fmt.Errorf("%w: mock", ErrCannotRender),
				},
			},
			format:    "other",
			value:     struct{}{},
			wantErr:   "render: unsupported format: other",
			wantErrIs: []error{Err, ErrUnsupportedFormat},
		},
		{
			name:      "non-existing renderer",
			renderers: map[string]FormatRenderer{},
			format:    "unknown",
			value:     struct{}{},
			wantErr:   "render: unsupported format: unknown",
			wantErrIs: []error{Err, ErrUnsupportedFormat},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &Renderer{
				Renderers: tt.renderers,
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
		for _, format := range tt.formats {
			t.Run(format+" format "+tt.name, func(t *testing.T) {
				w := &mockWriter{WriteErr: tt.writeErr}

				value := tt.value
				if tt.valueFunc != nil {
					value = tt.valueFunc()
				}

				var err error
				var panicRes any
				renderer := compactRenderer
				require.NoError(t, err)

				func() {
					defer func() {
						if r := recover(); r != nil {
							panicRes = r
						}
					}()
					err = renderer.Render(w, format, value)
				}()

				got := w.String()
				var want string
				if tt.wantPretty == "" && tt.wantCompact == "" {
					want = tt.want
				} else {
					want = tt.wantCompact
				}

				if tt.wantPanic != "" {
					assert.Equal(t, tt.wantPanic, panicRes)
				}

				if tt.wantErr != "" {
					wantErr := strings.ReplaceAll(
						tt.wantErr, "{{format}}", format,
					)
					assert.EqualError(t, err, wantErr)
				}
				for _, e := range tt.wantErrIs {
					assert.ErrorIs(t, err, e)
				}

				if tt.wantPanic == "" &&
					tt.wantErr == "" && len(tt.wantErrIs) == 0 {
					assert.NoError(t, err)
					assert.Equal(t, want, got)
				}
			})
		}
	}
}
