package render_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/jimeh/go-render"
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
		name      string
		pretty    bool
		prefix    string
		indent    string
		value     interface{}
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name:   "simple object without pretty",
			pretty: false,
			value:  map[string]int{"age": 30},
			want:   "{\"age\":30}\n",
		},
		{
			name:   "simple object with pretty",
			pretty: true,
			value:  map[string]int{"age": 30},
			want:   "{\n  \"age\": 30\n}\n",
		},
		{
			name:   "pretty with prefix and indent",
			pretty: true,
			prefix: "// ",
			indent: "\t",
			value:  map[string]int{"age": 30},
			want:   "{\n// \t\"age\": 30\n// }\n",
		},
		{
			name:   "prefix and indent without pretty",
			pretty: false,
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
			wantErrIs: []error{render.Err, render.ErrFailed},
		},
		{
			name:      "invalid value",
			pretty:    false,
			value:     make(chan int),
			wantErr:   "render: failed: json: unsupported type: chan int",
			wantErrIs: []error{render.Err, render.ErrFailed},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &render.JSON{
				Pretty: tt.pretty,
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
