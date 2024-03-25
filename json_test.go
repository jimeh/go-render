package render

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockJSONMarshaler struct {
	data []byte
	err  error
}

var _ json.Marshaler = (*mockJSONMarshaler)(nil)

func (mjm *mockJSONMarshaler) MarshalJSON() ([]byte, error) {
	return mjm.data, mjm.err
}

func TestJSON_Render(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		indent     string
		value      any
		want       string
		wantPretty string
		wantErr    string
		wantErrIs  []error
	}{
		{
			name:  "simple object",
			value: map[string]int{"age": 30},
			want:  "{\"age\":30}\n",
		},
		{
			name:   "ignores prefix and indent",
			prefix: "// ",
			indent: "\t",
			value:  map[string]int{"age": 30},
			want:   "{\"age\":30}\n",
		},
		{
			name:  "implements json.Marshaler",
			value: &mockJSONMarshaler{data: []byte(`{"age":30}`)},
			want:  "{\"age\":30}\n",
		},
		{
			name:      "error from json.Marshaler",
			value:     &mockJSONMarshaler{err: errors.New("marshal error!!1")},
			wantErrIs: []error{Err, ErrFailed},
		},
		{
			name:      "invalid value",
			value:     make(chan int),
			wantErr:   "render: failed: json: unsupported type: chan int",
			wantErrIs: []error{Err, ErrFailed},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSON{
				Prefix: tt.prefix,
				Indent: tt.indent,
			}
			var buf bytes.Buffer

			err := j.Render(&buf, tt.value)
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

func TestJSON_RenderPretty(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		indent     string
		value      any
		want       string
		wantPretty string
		wantErr    string
		wantErrIs  []error
	}{
		{
			name:  "simple object",
			value: map[string]int{"age": 30},
			want:  "{\n  \"age\": 30\n}\n",
		},
		{
			name:   "uses prefix and indent",
			prefix: "// ",
			indent: "\t",
			value:  map[string]int{"age": 30},
			want:   "{\n// \t\"age\": 30\n// }\n",
		},
		{
			name:  "implements json.Marshaler",
			value: &mockJSONMarshaler{data: []byte(`{"age":30}`)},
			want:  "{\n  \"age\": 30\n}\n",
		},
		{
			name:      "error from json.Marshaler",
			value:     &mockJSONMarshaler{err: errors.New("marshal error!!1")},
			wantErrIs: []error{Err, ErrFailed},
		},
		{
			name:      "invalid value",
			value:     make(chan int),
			wantErr:   "render: failed: json: unsupported type: chan int",
			wantErrIs: []error{Err, ErrFailed},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSON{
				Prefix: tt.prefix,
				Indent: tt.indent,
			}
			var buf bytes.Buffer

			err := j.RenderPretty(&buf, tt.value)
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

func TestJSON_Formats(t *testing.T) {
	h := &JSON{}

	assert.Equal(t, []string{"json"}, h.Formats())
}
