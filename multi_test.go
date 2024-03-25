package render

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var multiHandlerTestCases = []struct {
	name       string
	handlers   []Handler
	value      any
	want       string
	wantPretty string
	wantErr    string
	wantErrIs  []error
}{
	{
		name: "no handler can render",
		handlers: []Handler{
			&mockHandler{err: ErrCannotRender},
			&mockHandler{err: ErrCannotRender},
		},
		value:     "test",
		wantErr:   "render: cannot render: string",
		wantErrIs: []error{ErrCannotRender},
	},
	{
		name: "one handler can render",
		handlers: []Handler{
			&mockHandler{err: ErrCannotRender},
			&mockHandler{output: "success output"},
			&mockHandler{err: ErrCannotRender},
		},
		value:      struct{}{},
		want:       "success output",
		wantPretty: "success output",
	},
	{
		name: "one pretty handler can render",
		handlers: []Handler{
			&mockHandler{err: ErrCannotRender},
			&mockPrettyHandler{
				output:       "success output",
				prettyOutput: "pretty success output",
			},
			&mockHandler{err: ErrCannotRender},
		},
		value:      struct{}{},
		want:       "success output",
		wantPretty: "pretty success output",
	},
	{
		name: "multiple handlers can render",
		handlers: []Handler{
			&mockHandler{err: ErrCannotRender},
			&mockHandler{output: "first output"},
			&mockHandler{output: "second output"},
		},
		value:      struct{}{},
		want:       "first output",
		wantPretty: "first output",
	},
	{
		name: "multiple pretty handlers can render",
		handlers: []Handler{
			&mockHandler{err: ErrCannotRender},
			&mockPrettyHandler{
				output:       "first output",
				prettyOutput: "pretty first output",
			},
			&mockPrettyHandler{
				output:       "second output",
				prettyOutput: "pretty second output",
			},
		},
		value:      struct{}{},
		want:       "first output",
		wantPretty: "pretty first output",
	},
	{
		name: "first handler fails",
		handlers: []Handler{
			&mockHandler{err: errors.New("mock error")},
			&mockHandler{output: "success output"},
		},
		value:   struct{}{},
		wantErr: "mock error",
	},
	{
		name: "fails after cannot render",
		handlers: []Handler{
			&mockHandler{err: ErrCannotRender},
			&mockHandler{err: errors.New("mock error")},
			&mockHandler{output: "success output"},
		},
		value:   struct{}{},
		wantErr: "mock error",
	},
	{
		name: "fails after success render",
		handlers: []Handler{
			&mockHandler{output: "success output"},
			&mockHandler{err: errors.New("mock error")},
			&mockHandler{err: ErrCannotRender},
		},
		value:      struct{}{},
		want:       "success output",
		wantPretty: "success output",
	},
	{
		name: "fails after success render with prettier handlers",
		handlers: []Handler{
			&mockPrettyHandler{
				output:       "success output",
				prettyOutput: "pretty success output",
			},
			&mockHandler{err: errors.New("mock error")},
			&mockHandler{err: ErrCannotRender},
		},
		value:      struct{}{},
		want:       "success output",
		wantPretty: "pretty success output",
	},
}

func TestMulti_Render(t *testing.T) {
	for _, tt := range multiHandlerTestCases {
		t.Run(tt.name, func(t *testing.T) {
			mr := &Multi{
				Handlers: tt.handlers,
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

func TestMulti_RenderPretty(t *testing.T) {
	for _, tt := range multiHandlerTestCases {
		t.Run(tt.name, func(t *testing.T) {
			mr := &Multi{
				Handlers: tt.handlers,
			}
			var buf bytes.Buffer

			err := mr.RenderPretty(&buf, tt.value)
			got := buf.String()

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
			for _, e := range tt.wantErrIs {
				assert.ErrorIs(t, err, e)
			}

			if tt.wantErr == "" && len(tt.wantErrIs) == 0 {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPretty, got)
			}
		})
	}
}

func TestMulti_Formats(t *testing.T) {
	tests := []struct {
		name     string
		handlers []Handler
		want     []string
	}{
		{
			name: "single handler without a Formats method",
			handlers: []Handler{
				&mockHandler{},
			},
			want: []string{},
		},
		{
			name: "multiple handlers without a Formats method",
			handlers: []Handler{
				&mockHandler{},
			},
			want: []string{},
		},
		{
			name: "single handler with a Formats method",
			handlers: []Handler{
				&mockFormatsHandler{formats: []string{"yaml", "yml"}},
			},
			want: []string{"yaml", "yml"},
		},
		{
			name: "multiple handlers without a Formats method",
			handlers: []Handler{
				&mockFormatsHandler{formats: []string{"yaml", "yml"}},
				&mockFormatsHandler{formats: []string{"text", "txt"}},
			},
			want: []string{"yaml", "yml", "text", "txt"},
		},
		{
			name: "mixture of handlers with and without a Formats method",
			handlers: []Handler{
				&mockFormatsHandler{formats: []string{"yaml", "yml"}},
				&mockHandler{},
				&mockFormatsHandler{formats: []string{"binary", "bin"}},
				&mockHandler{},
			},
			want: []string{"yaml", "yml", "binary", "bin"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &Multi{
				Handlers: tt.handlers,
			}

			got := mr.Formats()

			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
