package render_test

import (
	"bytes"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			indent: "  ",
			value:  map[string]int{"age": 30},
			want:   "{\n  \"age\": 30\n}\n",
		},
		{
			name:   "with prefix and indent",
			pretty: true,
			prefix: "// ",
			indent: "\t",
			value:  map[string]int{"age": 30},
			want:   "{\n// \t\"age\": 30\n// }\n",
		},
		{
			name:      "invalid value",
			pretty:    false,
			value:     make(chan int),
			wantErr:   "render: json: unsupported type: chan int",
			wantErrIs: []error{render.Err},
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

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tt.wantErr)
			}
			for _, e := range tt.wantErrIs {
				assert.ErrorIs(t, err, e)
			}

			if tt.wantErr == "" && len(tt.wantErrIs) == 0 {
				require.NoError(t, err)
				got := buf.String()
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
