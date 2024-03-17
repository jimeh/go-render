package render_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type mockYAMLMarshaler struct {
	val any
	err error
}

var _ yaml.Marshaler = (*mockYAMLMarshaler)(nil)

func (m *mockYAMLMarshaler) MarshalYAML() (any, error) {
	return m.val, m.err
}

func TestYAML_Render(t *testing.T) {
	tests := []struct {
		name      string
		indent    int
		value     interface{}
		want      string
		wantErr   string
		wantErrIs []error
		wantPanic string
	}{
		{
			name:  "simple object default indent",
			value: map[string]int{"age": 30},
			want:  "age: 30\n",
		},
		{
			name:   "nested structure",
			indent: 0, // This will use the default indent of 2 spaces
			value: map[string]any{
				"user": map[string]any{
					"age":  30,
					"name": "John Doe",
				},
			},
			want: "user:\n  age: 30\n  name: John Doe\n",
		},
		{
			name:   "simple object custom indent",
			indent: 4,
			value: map[string]any{
				"user": map[string]any{
					"age":  30,
					"name": "John Doe",
				},
			},
			want: "user:\n    age: 30\n    name: John Doe\n",
		},
		{
			name:  "implements yaml.Marshaler",
			value: &mockYAMLMarshaler{val: map[string]int{"age": 30}},
			want:  "age: 30\n",
		},
		{
			name:      "error from yaml.Marshaler",
			value:     &mockYAMLMarshaler{err: errors.New("mock error")},
			wantErr:   "render: mock error",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "invalid value",
			indent:    0,
			value:     make(chan int),
			wantPanic: "cannot marshal type: chan int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &render.YAML{
				Indent: tt.indent,
			}

			var buf bytes.Buffer
			var err error
			var panicRes any
			func() {
				defer func() {
					if r := recover(); r != nil {
						panicRes = r
					}
				}()
				err = j.Render(&buf, tt.value)
			}()

			got := buf.String()

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
