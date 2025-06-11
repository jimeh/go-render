package render

import (
	"errors"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

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
		indent         int
		encoderOptions []yaml.EncodeOption
		value          any
		want           string
		writeErr       error
		wantErr        string
		wantErrIs      []error
		wantPanic      string
	}{
		{
			name:  "simple object default indent",
			value: map[string]int{"age": 30},
			want:  "age: 30\n",
		},
		{
			name: "nested structure",
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
				"books": []string{
					"The Great Gatsby",
					"1984",
				},
			},
			want: "books:\n  - The Great Gatsby\n  - \"1984\"\n",
		},
		{
			name: "sequences without SequenceIndent",
			encoderOptions: []yaml.EncodeOption{
				yaml.IndentSequence(false),
			},
			value: map[string]any{
				"books": []string{
					"The Great Gatsby",
					"1984",
				},
			},
			want: "books:\n- The Great Gatsby\n- \"1984\"\n",
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
			indent:    0,
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

			var err error
			var panicRes any
			func() {
				defer func() {
					if r := recover(); r != nil {
						panicRes = r
					}
				}()
				err = j.Render(w, tt.value)
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

func TestYAML_Formats(t *testing.T) {
	h := &YAML{}

	assert.Equal(t, []string{"yaml", "yml"}, h.Formats())
}
