package render

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		handlers map[string]Handler
		want     *Renderer
	}{
		{
			name: "nil handlers map",
			want: &Renderer{Handlers: map[string]Handler{}},
		},
		{
			name:     "empty handlers map",
			handlers: map[string]Handler{},
			want:     &Renderer{Handlers: map[string]Handler{}},
		},
		{
			name: "single handler",
			handlers: map[string]Handler{
				"mock": &mockHandler{},
			},
			want: &Renderer{Handlers: map[string]Handler{
				"mock": &mockHandler{},
			}},
		},
		{
			name: "multiple handlers",
			handlers: map[string]Handler{
				"mock":  &mockHandler{},
				"other": &mockHandler{output: "other output"},
			},
			want: &Renderer{Handlers: map[string]Handler{
				"mock":  &mockHandler{},
				"other": &mockHandler{output: "other output"},
			}},
		},
		{
			name: "multiple handlers with alias formats",
			handlers: map[string]Handler{
				"mock":  &mockFormatsHandler{formats: []string{"mock", "m"}},
				"other": &mockFormatsHandler{formats: []string{"other", "o"}},
			},
			want: &Renderer{Handlers: map[string]Handler{
				"mock":  &mockFormatsHandler{formats: []string{"mock", "m"}},
				"m":     &mockFormatsHandler{formats: []string{"mock", "m"}},
				"other": &mockFormatsHandler{formats: []string{"other", "o"}},
				"o":     &mockFormatsHandler{formats: []string{"other", "o"}},
			}},
		},
		{
			name: "multiple handlers with custom formats",
			handlers: map[string]Handler{
				"foo": &mockFormatsHandler{formats: []string{"mock", "m"}},
				"bar": &mockFormatsHandler{formats: []string{"other", "o"}},
			},
			want: &Renderer{Handlers: map[string]Handler{
				"foo":   &mockFormatsHandler{formats: []string{"mock", "m"}},
				"mock":  &mockFormatsHandler{formats: []string{"mock", "m"}},
				"m":     &mockFormatsHandler{formats: []string{"mock", "m"}},
				"bar":   &mockFormatsHandler{formats: []string{"other", "o"}},
				"other": &mockFormatsHandler{formats: []string{"other", "o"}},
				"o":     &mockFormatsHandler{formats: []string{"other", "o"}},
			}},
		},
		{
			name: "multiple handlers with capitalized formats",
			handlers: map[string]Handler{
				"Foo": &mockFormatsHandler{formats: []string{"MOCK", "m"}},
				"Bar": &mockFormatsHandler{formats: []string{"OTHER", "o"}},
			},
			want: &Renderer{Handlers: map[string]Handler{
				"foo":   &mockFormatsHandler{formats: []string{"MOCK", "m"}},
				"mock":  &mockFormatsHandler{formats: []string{"MOCK", "m"}},
				"m":     &mockFormatsHandler{formats: []string{"MOCK", "m"}},
				"bar":   &mockFormatsHandler{formats: []string{"OTHER", "o"}},
				"other": &mockFormatsHandler{formats: []string{"OTHER", "o"}},
				"o":     &mockFormatsHandler{formats: []string{"OTHER", "o"}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.handlers)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRenderer_Add(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		handler Handler
		want    []string
	}{
		{
			name:    "handler without Formats",
			format:  "tackle",
			handler: &mockHandler{},
			want:    []string{"tackle"},
		},
		{
			name:    "hander with Formats",
			format:  "hackle",
			handler: &mockFormatsHandler{formats: []string{"hackle"}},
			want:    []string{"hackle"},
		},
		{
			name:    "hander with alias formats",
			format:  "hackle",
			handler: &mockFormatsHandler{formats: []string{"hackle", "hack"}},
			want:    []string{"hackle", "hack"},
		},
		{
			name:    "given format differs from Formats",
			format:  "foobar",
			handler: &mockFormatsHandler{formats: []string{"hackle", "hack"}},
			want:    []string{"foobar", "hackle", "hack"},
		},
		{
			name:    "lowercases capitalized formats",
			format:  "FooBar",
			handler: &mockFormatsHandler{formats: []string{"HACKLE", "Hack"}},
			want:    []string{"foobar", "hackle", "hack"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{Handlers: map[string]Handler{}}

			r.Add(tt.format, tt.handler)

			for _, f := range tt.want {
				got, ok := r.Handlers[f]
				assert.Truef(t, ok, "not added as %q format", f)
				assert.Equal(t, tt.handler, got)
			}

			gotFormats := []string{}
			for f := range r.Handlers {
				gotFormats = append(gotFormats, f)
			}
			assert.ElementsMatch(t, tt.want, gotFormats)
		})
	}
}

func TestRenderer_Render(t *testing.T) {
	tests := []struct {
		name      string
		handlers  map[string]Handler
		format    string
		pretty    bool
		value     any
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "no pretty with handler that supports pretty",
			handlers: map[string]Handler{
				"mock": &mockPrettyHandler{
					output:       "plain output",
					prettyOutput: "pretty output",
				},
			},
			format: "mock",
			pretty: false,
			value:  struct{}{},
			want:   "plain output",
		},
		{
			name: "no pretty with handler that does not support pretty",
			handlers: map[string]Handler{
				"mock": &mockHandler{output: "plain output"},
			},
			format: "mock",
			pretty: false,
			value:  struct{}{},
			want:   "plain output",
		},
		{
			name: "pretty with handler that supports pretty",
			handlers: map[string]Handler{
				"mock": &mockPrettyHandler{
					output:       "plain output",
					prettyOutput: "pretty output",
				},
			},
			format: "mock",
			pretty: true,
			value:  struct{}{},
			want:   "pretty output",
		},
		{
			name: "pretty with handler that does not support pretty",
			handlers: map[string]Handler{
				"mock": &mockHandler{
					output: "plain output",
				},
			},
			format: "mock",
			pretty: true,
			value:  struct{}{},
			want:   "plain output",
		},
		{
			name: "normalizes given format to lowercase",
			handlers: map[string]Handler{
				"mock": &mockHandler{
					output: "plain output",
				},
			},
			format: "MOCK",
			pretty: true,
			value:  struct{}{},
			want:   "plain output",
		},
		{
			name: "handler returns error",
			handlers: map[string]Handler{
				"other": &mockHandler{
					output: "mock output",
					err:    errors.New("mock error"),
				},
			},
			format:  "other",
			value:   struct{}{},
			wantErr: "render: failed: mock error",
		},
		{
			name: "handler returns ErrCannotRender",
			handlers: map[string]Handler{
				"other": &mockHandler{
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
			name:      "non-existing handler",
			handlers:  map[string]Handler{},
			format:    "unknown",
			value:     struct{}{},
			wantErr:   "render: unsupported format: unknown",
			wantErrIs: []error{Err, ErrUnsupportedFormat},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				Handlers: tt.handlers,
			}
			var buf bytes.Buffer

			err := r.Render(&buf, tt.format, tt.pretty, tt.value)
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

func TestRenderer_Compact(t *testing.T) {
	tests := []struct {
		name      string
		handlers  map[string]Handler
		format    string
		value     any
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "handler supports pretty",
			handlers: map[string]Handler{
				"mock": &mockPrettyHandler{
					output:       "plain output",
					prettyOutput: "pretty output",
				},
			},
			format: "mock",
			value:  struct{}{},
			want:   "plain output",
		},
		{
			name: "handler does not support pretty",
			handlers: map[string]Handler{
				"mock": &mockHandler{output: "plain output"},
			},
			format: "mock",
			value:  struct{}{},
			want:   "plain output",
		},
		{
			name: "handler returns error",
			handlers: map[string]Handler{
				"other": &mockHandler{
					output: "mock output",
					err:    errors.New("mock error"),
				},
			},
			format:  "other",
			value:   struct{}{},
			wantErr: "render: failed: mock error",
		},
		{
			name: "handler returns ErrCannotRender",
			handlers: map[string]Handler{
				"other": &mockHandler{
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
			name:      "non-existing handler",
			handlers:  map[string]Handler{},
			format:    "unknown",
			value:     struct{}{},
			wantErr:   "render: unsupported format: unknown",
			wantErrIs: []error{Err, ErrUnsupportedFormat},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				Handlers: tt.handlers,
			}
			var buf bytes.Buffer

			err := r.Compact(&buf, tt.format, tt.value)
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

func TestRenderer_Pretty(t *testing.T) {
	tests := []struct {
		name      string
		handlers  map[string]Handler
		format    string
		value     any
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "handler supports pretty",
			handlers: map[string]Handler{
				"mock": &mockPrettyHandler{
					output:       "plain output",
					prettyOutput: "pretty output",
				},
			},
			format: "mock",
			value:  struct{}{},
			want:   "pretty output",
		},
		{
			name: "handler does not support pretty",
			handlers: map[string]Handler{
				"mock": &mockHandler{
					output: "plain output",
				},
			},
			format: "mock",
			value:  struct{}{},
			want:   "plain output",
		},
		{
			name: "handler returns error",
			handlers: map[string]Handler{
				"other": &mockHandler{
					output: "mock output",
					err:    errors.New("mock error"),
				},
			},
			format:  "other",
			value:   struct{}{},
			wantErr: "render: failed: mock error",
		},
		{
			name: "handler returns ErrCannotRender",
			handlers: map[string]Handler{
				"other": &mockHandler{
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
			name:      "non-existing handler",
			handlers:  map[string]Handler{},
			format:    "unknown",
			value:     struct{}{},
			wantErr:   "render: unsupported format: unknown",
			wantErrIs: []error{Err, ErrUnsupportedFormat},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{
				Handlers: tt.handlers,
			}
			var buf bytes.Buffer

			err := r.Pretty(&buf, tt.format, tt.value)
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
		for _, pretty := range []bool{false, true} {
			for _, format := range tt.formats {
				name := format + " format " + tt.name
				if pretty {
					name = "pretty " + name
				}

				t.Run(name, func(t *testing.T) {
					w := &mockWriter{WriteErr: tt.writeErr}

					value := tt.value
					if tt.valueFunc != nil {
						value = tt.valueFunc()
					}

					var err error
					var panicRes any

					func() {
						defer func() {
							if r := recover(); r != nil {
								panicRes = r
							}
						}()
						err = Base.Render(w, format, pretty, value)
					}()

					got := w.String()
					want := tt.want
					if pretty && tt.wantPretty != "" {
						want = tt.wantPretty
					} else if tt.wantCompact != "" {
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
}

func TestRenderer_CompactAllFormats(t *testing.T) {
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

				func() {
					defer func() {
						if r := recover(); r != nil {
							panicRes = r
						}
					}()
					err = Base.Compact(w, format, value)
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

func TestRenderer_PrettyAllFormats(t *testing.T) {
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

				func() {
					defer func() {
						if r := recover(); r != nil {
							panicRes = r
						}
					}()
					err = Base.Pretty(w, format, value)
				}()

				got := w.String()
				var want string
				if tt.wantPretty == "" && tt.wantCompact == "" {
					want = tt.want
				} else {
					want = tt.wantPretty
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
