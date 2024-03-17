package render_test

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXML_Render(t *testing.T) {
	tests := []struct {
		name      string
		pretty    bool
		prefix    string
		indent    string
		value     any
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name:   "simple object without pretty",
			pretty: false,
			value: struct {
				XMLName xml.Name `xml:"user"`
				Age     int      `xml:"age"`
			}{Age: 30},
			want: `<user><age>30</age></user>`,
		},
		{
			name:   "simple object with pretty",
			pretty: true,
			indent: "    ",
			value: struct {
				XMLName xml.Name `xml:"user"`
				Age     int      `xml:"age"`
			}{Age: 30},
			want: "<user>\n    <age>30</age>\n</user>",
		},
		{
			name:   "with prefix and indent",
			pretty: true,
			prefix: "//",
			indent: "\t",
			value: struct {
				XMLName xml.Name `xml:"user"`
				Age     int      `xml:"age"`
			}{Age: 30},
			want: "//<user>\n//\t<age>30</age>\n//</user>",
		},
		{
			name:      "invalid value",
			pretty:    false,
			value:     make(chan int),
			wantErr:   "render: xml: unsupported type: chan int",
			wantErrIs: []error{render.Err},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &render.XML{
				Pretty: tt.pretty,
				Prefix: tt.prefix,
				Indent: tt.indent,
			}

			var buf bytes.Buffer
			err := x.Render(&buf, tt.value)

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
