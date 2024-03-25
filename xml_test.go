package render

import (
	"bytes"
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockXMLMarshaler struct {
	elm string
	err error
}

var _ xml.Marshaler = (*mockXMLMarshaler)(nil)

func (mxm *mockXMLMarshaler) MarshalXML(
	e *xml.Encoder,
	start xml.StartElement,
) error {
	err := e.EncodeElement(mxm.elm, start)

	if mxm.err != nil {
		return mxm.err
	}

	return err
}

func TestXML_Render(t *testing.T) {
	tests := []struct {
		name      string
		prefix    string
		indent    string
		value     any
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "simple object",
			value: struct {
				XMLName xml.Name `xml:"user"`
				Age     int      `xml:"age"`
			}{Age: 30},
			want: `<user><age>30</age></user>`,
		},
		{
			name:   "ignores indent without pretty",
			prefix: "//",
			indent: "\t",
			value: struct {
				XMLName xml.Name `xml:"user"`
				Age     int      `xml:"age"`
			}{Age: 30},
			want: `<user><age>30</age></user>`,
		},
		{
			name:  "implements xml.Marshaler",
			value: &mockXMLMarshaler{elm: "test string"},
			want:  "<mockXMLMarshaler>test string</mockXMLMarshaler>",
		},
		{
			name:      "error from xml.Marshaler",
			value:     &mockXMLMarshaler{err: errors.New("mock error")},
			wantErr:   "render: failed: mock error",
			wantErrIs: []error{Err, ErrFailed},
		},
		{
			name:      "invalid value",
			value:     make(chan int),
			wantErr:   "render: failed: xml: unsupported type: chan int",
			wantErrIs: []error{Err, ErrFailed},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &XML{
				Prefix: tt.prefix,
				Indent: tt.indent,
			}
			var buf bytes.Buffer

			err := x.Render(&buf, tt.value)
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

func TestXML_RenderPretty(t *testing.T) {
	tests := []struct {
		name      string
		prefix    string
		indent    string
		value     any
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name: "simple object",
			value: struct {
				XMLName xml.Name `xml:"user"`
				Age     int      `xml:"age"`
			}{Age: 30},
			want: "<user>\n  <age>30</age>\n</user>",
		},
		{
			name:   "uses prefix and indent",
			prefix: "//",
			indent: "\t",
			value: struct {
				XMLName xml.Name `xml:"user"`
				Age     int      `xml:"age"`
			}{Age: 30},
			want: "//<user>\n//\t<age>30</age>\n//</user>",
		},
		{
			name:  "implements xml.Marshaler",
			value: &mockXMLMarshaler{elm: "test string"},
			want:  "<mockXMLMarshaler>test string</mockXMLMarshaler>",
		},
		{
			name:      "error from xml.Marshaler",
			value:     &mockXMLMarshaler{err: errors.New("mock error")},
			wantErr:   "render: failed: mock error",
			wantErrIs: []error{Err, ErrFailed},
		},
		{
			name:      "invalid value",
			value:     make(chan int),
			wantErr:   "render: failed: xml: unsupported type: chan int",
			wantErrIs: []error{Err, ErrFailed},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &XML{
				Prefix: tt.prefix,
				Indent: tt.indent,
			}
			var buf bytes.Buffer

			err := x.RenderPretty(&buf, tt.value)
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

func TestXML_Formats(t *testing.T) {
	h := &XML{}

	assert.Equal(t, []string{"xml"}, h.Formats())
}
