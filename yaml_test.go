package render

import (
	"errors"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T {
	return &v
}

type mockYAMLInterfaceMarshaler struct {
	val any
	err error
}

var _ yaml.InterfaceMarshaler = (*mockYAMLInterfaceMarshaler)(nil)

func (m *mockYAMLInterfaceMarshaler) MarshalYAML() (any, error) {
	return m.val, m.err
}

type mockYAMLBytesMarshaler struct {
	val []byte
	err error
}

var _ yaml.BytesMarshaler = (*mockYAMLBytesMarshaler)(nil)

func (m *mockYAMLBytesMarshaler) MarshalYAML() ([]byte, error) {
	return m.val, m.err
}

func TestYAML_Render(t *testing.T) {
	tests := []struct {
		name           string
		encoderOptions []yaml.EncodeOption
		indent         int
		value          any
		want           string
		writeErr       error
		wantErr        string
		wantErrIs      []error
	}{
		{
			name:  "simple object",
			value: map[string]int{"age": 30},
			want:  "age: 30\n",
		},
		{
			name: "nested object",
			value: map[string]any{
				"user": map[string]any{
					"age":  30,
					"name": "John Doe",
				},
			},
			want: "user:\n  age: 30\n  name: John Doe\n",
		},
		{
			name: "sequences",
			value: map[string]any{
				"items": map[string]any{
					"books": []string{
						"The Great Gatsby",
						"1984",
					},
				},
			},
			want: "items:\n  books:\n    - The Great Gatsby\n    - \"1984\"\n",
		},
		{
			name: "sequences without IndentSequence",
			encoderOptions: []yaml.EncodeOption{
				yaml.IndentSequence(false),
			},
			value: map[string]any{
				"items": map[string]any{
					"books": []string{
						"The Great Gatsby",
						"1984",
					},
				},
			},
			want: "items:\n  books:\n  - The Great Gatsby\n  - \"1984\"\n",
		},
		{
			name: "custom AutoInt encoder option",
			encoderOptions: []yaml.EncodeOption{
				yaml.AutoInt(),
			},
			value: map[string]any{
				"age": 1.0,
			},
			want: "age: 1\n",
		},
		{
			name:   "nested object with custom indent",
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
			name:   "sequences with custom indent",
			indent: 4,
			value: map[string]any{
				"items": map[string]any{
					"books": []string{
						"The Great Gatsby",
						"1984",
					},
				},
			},
			want: "items:\n    " +
				"books:\n        - The Great Gatsby\n        - \"1984\"\n",
		},
		{
			name:   "sequences with custom indent and without IndentSequence",
			indent: 4,
			encoderOptions: []yaml.EncodeOption{
				yaml.IndentSequence(false),
			},
			value: map[string]any{
				"items": map[string]any{
					"books": []string{
						"The Great Gatsby",
						"1984",
					},
				},
			},
			want: "items:\n    " +
				"books:\n    - The Great Gatsby\n    - \"1984\"\n",
		},
		{
			name:  "implements yaml.InterfaceMarshaler",
			value: &mockYAMLInterfaceMarshaler{val: map[string]int{"age": 30}},
			want:  "age: 30\n",
		},
		{
			name: "error from yaml.InterfaceMarshaler",
			value: &mockYAMLInterfaceMarshaler{
				err: errors.New("mock error"),
			},
			wantErr:   "render: failed: mock error",
			wantErrIs: []error{Err, ErrFailed},
		},
		{
			name:  "implements yaml.BytesMarshaler",
			value: &mockYAMLBytesMarshaler{val: []byte("age: 30\n")},
			want:  "age: 30\n",
		},
		{
			name:      "error from yaml.BytesMarshaler",
			value:     &mockYAMLBytesMarshaler{err: errors.New("mock error")},
			wantErr:   "render: failed: mock error",
			wantErrIs: []error{Err, ErrFailed},
		},
		{
			name:      "yaml format with error writing to writer",
			writeErr:  errors.New("write error!!1"),
			value:     map[string]int{"age": 30},
			wantErr:   "render: failed: yaml: write error: write error!!1",
			wantErrIs: []error{Err, ErrFailed},
		},
		{
			name:      "invalid value",
			value:     make(chan int),
			wantErr:   "render: failed: unknown value type chan int",
			wantErrIs: []error{Err, ErrFailed},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &YAML{
				Indent:        tt.indent,
				EncodeOptions: tt.encoderOptions,
			}

			w := &mockWriter{WriteErr: tt.writeErr}

			err := j.Render(w, tt.value)
			got := w.String()

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

func TestYAML_Formats(t *testing.T) {
	h := &YAML{}

	assert.Equal(t, []string{"yaml", "yml"}, h.Formats())
}

func TestYAML_DefaultIndent(t *testing.T) {
	tests := []struct {
		name           string
		defaultVal     *int
		indent         int
		encoderOptions []yaml.EncodeOption
		want           string
	}{
		{
			name: "do not set default, indent, or encoder options",
			want: "items:\n  " +
				"books:\n    - The Great Gatsby\n    - \"1984\"\n",
		},
		{
			name:       "custom default takes precedence over default",
			defaultVal: ptr(4),
			want: "items:\n    " +
				"books:\n        - The Great Gatsby\n        - \"1984\"\n",
		},
		{
			name:   "indent takes precedence over default",
			indent: 4,
			want: "items:\n    " +
				"books:\n        - The Great Gatsby\n        - \"1984\"\n",
		},
		{
			name: "encoder option takes precedence over default",
			encoderOptions: []yaml.EncodeOption{
				yaml.Indent(4),
			},
			want: "items:\n    " +
				"books:\n        - The Great Gatsby\n        - \"1984\"\n",
		},
		{
			name:       "indent takes precedence over custom default",
			defaultVal: ptr(4),
			indent:     3,
			want: "items:\n   " +
				"books:\n      - The Great Gatsby\n      - \"1984\"\n",
		},
		{
			name:   "encoder option takes precedence over indent",
			indent: 4,
			encoderOptions: []yaml.EncodeOption{
				yaml.Indent(3),
			},
			want: "items:\n   " +
				"books:\n      - The Great Gatsby\n      - \"1984\"\n",
		},
		{
			name: "encoder option takes precedence over indent and " +
				"custom default",
			defaultVal: ptr(5),
			indent:     4,
			encoderOptions: []yaml.EncodeOption{
				yaml.Indent(3),
			},
			want: "items:\n   " +
				"books:\n      - The Great Gatsby\n      - \"1984\"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture and restore the original default value.
			originalDefault := YAMLDefaultIndent
			t.Cleanup(func() { YAMLDefaultIndent = originalDefault })

			if tt.defaultVal != nil {
				YAMLDefaultIndent = *tt.defaultVal
			}

			y := &YAML{
				Indent:        tt.indent,
				EncodeOptions: tt.encoderOptions,
			}
			w := &mockWriter{}

			value := map[string]any{
				"items": map[string]any{
					"books": []string{
						"The Great Gatsby",
						"1984",
					},
				},
			}

			err := y.Render(w, value)
			require.NoError(t, err)

			assert.Equal(t, tt.want, w.String())
		})
	}
}

func TestYAML_DefaultIndentSequence(t *testing.T) {
	tests := []struct {
		name           string
		defaultVal     *bool
		encoderOptions []yaml.EncodeOption
		want           string
	}{
		{
			name: "do not set default or encoder options",
			want: "items:\n  " +
				"books:\n    - The Great Gatsby\n    - \"1984\"\n",
		},
		{
			name:       "set default to true",
			defaultVal: ptr(true),
			want: "items:\n  " +
				"books:\n    - The Great Gatsby\n    - \"1984\"\n",
		},
		{
			name:       "set default to true and encoder option to false",
			defaultVal: ptr(true),
			encoderOptions: []yaml.EncodeOption{
				yaml.IndentSequence(false),
			},
			want: "items:\n  " +
				"books:\n  - The Great Gatsby\n  - \"1984\"\n",
		},
		{
			name:       "set default to false",
			defaultVal: ptr(false),
			want: "items:\n  " +
				"books:\n  - The Great Gatsby\n  - \"1984\"\n",
		},
		{
			name:       "set default to false and encoder option to true",
			defaultVal: ptr(false),
			encoderOptions: []yaml.EncodeOption{
				yaml.IndentSequence(true),
			},
			want: "items:\n  " +
				"books:\n    - The Great Gatsby\n    - \"1984\"\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture and restore the original default value.
			originalDefault := YAMLDefaultIndentSequence
			t.Cleanup(func() { YAMLDefaultIndentSequence = originalDefault })

			if tt.defaultVal != nil {
				YAMLDefaultIndentSequence = *tt.defaultVal
			}

			y := &YAML{
				EncodeOptions: tt.encoderOptions,
			}
			w := &mockWriter{}

			value := map[string]any{
				"items": map[string]any{
					"books": []string{
						"The Great Gatsby",
						"1984",
					},
				},
			}

			err := y.Render(w, value)
			require.NoError(t, err)

			assert.Equal(t, tt.want, w.String())
		})
	}
}
