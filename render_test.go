package render_test

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
)

type mockWriter struct {
	WriteErr error
	buf      bytes.Buffer
}

var _ io.Writer = (*mockWriter)(nil)

func (mw *mockWriter) Write(p []byte) (n int, err error) {
	if mw.WriteErr != nil {
		return 0, mw.WriteErr
	}

	return mw.buf.Write(p)
}

func (mw *mockWriter) String() string {
	return mw.buf.String()
}

type mockRenderer struct {
	output string
	err    error
}

var _ render.FormatRenderer = (*mockRenderer)(nil)

func (m *mockRenderer) Render(w io.Writer, _ any) error {
	_, err := w.Write([]byte(m.output))

	if m.err != nil {
		return m.err
	}

	return err
}

func TestDefaultJSON(t *testing.T) {
	assert.Equal(t, &render.JSON{Pretty: true}, render.DefaultJSON)
}

func TestDefaultXML(t *testing.T) {
	assert.Equal(t, &render.XML{Pretty: true}, render.DefaultXML)
}

func TestDefaultYAML(t *testing.T) {
	assert.Equal(t, &render.YAML{Indent: 2}, render.DefaultYAML)
}

func TestDefaultText(t *testing.T) {
	assert.Equal(t, &render.Text{}, render.DefaultText)
}

func TestDefaultBinary(t *testing.T) {
	assert.Equal(t, &render.Binary{}, render.DefaultBinary)
}

func TestDefaultRenderer(t *testing.T) {
	want := &render.Renderer{
		Formats: map[string]render.FormatRenderer{
			"json": render.DefaultJSON,
			"text": render.DefaultText,
			"yaml": render.DefaultYAML,
		},
	}

	assert.Equal(t, want, render.DefaultRenderer)
}

type renderFormatTestCase struct {
	name      string
	writeErr  error
	format    string
	value     any
	want      string
	wantErr   string
	wantErrIs []error
	wantPanic string
}

// "binary" format.
var binaryFormattestCases = []renderFormatTestCase{
	{
		name:   "binary format with binary marshaler",
		format: "binary",
		value:  &mockBinaryMarshaler{data: []byte("test string")},
		want:   "test string",
	},
	{
		name:      "binary format without binary marshaler",
		format:    "binary",
		value:     struct{}{},
		wantErr:   "render: unsupported format: binary",
		wantErrIs: []error{render.Err, render.ErrUnsupportedFormat},
	},
	{
		name:   "binary format with error marshaling",
		format: "binary",
		value: &mockBinaryMarshaler{
			data: []byte("test string"),
			err:  errors.New("marshal error!!1"),
		},
		wantErr:   "render: failed: marshal error!!1",
		wantErrIs: []error{render.Err, render.ErrFailed},
	},
	{
		name:      "binary format with error writing to writer",
		format:    "binary",
		writeErr:  errors.New("write error!!1"),
		value:     &mockBinaryMarshaler{data: []byte("test string")},
		wantErr:   "render: failed: write error!!1",
		wantErrIs: []error{render.Err, render.ErrFailed},
	},
	{
		name:      "binary format with invalid type",
		format:    "binary",
		value:     make(chan int),
		wantErr:   "render: unsupported format: binary",
		wantErrIs: []error{render.Err, render.ErrUnsupportedFormat},
	},
}

// "json" format.
var jsonFormatTestCases = []renderFormatTestCase{
	{
		name:   "json format",
		format: "json",
		value:  map[string]int{"age": 30},
		want:   "{\n  \"age\": 30\n}\n",
	},
	{
		name:   "json format with json marshaler",
		format: "json",
		value:  &mockJSONMarshaler{data: []byte(`{"age":30}`)},
		want:   "{\n  \"age\": 30\n}\n",
	},
	{
		name:      "json format with error from json marshaler",
		format:    "json",
		value:     &mockJSONMarshaler{err: errors.New("marshal error!!1")},
		wantErrIs: []error{render.Err},
	},
	{
		name:      "json format with error writing to writer",
		format:    "json",
		writeErr:  errors.New("write error!!1"),
		value:     map[string]int{"age": 30},
		wantErr:   "render: failed: write error!!1",
		wantErrIs: []error{render.Err, render.ErrFailed},
	},
	{
		name:      "json format with invalid type",
		format:    "json",
		value:     make(chan int),
		wantErr:   "render: failed: json: unsupported type: chan int",
		wantErrIs: []error{render.Err, render.ErrFailed},
	},
}

// "text" format.
var textFormatTestCases = []renderFormatTestCase{
	{
		name:   "text format with fmt.Stringer",
		format: "text",
		value:  &mockStringer{value: "test string"},
		want:   "test string",
	},
	{
		name:   "text format with io.WriterTo",
		format: "text",
		value:  &mockWriterTo{value: "test string"},
		want:   "test string",
	},
	{
		name:      "text format without fmt.Stringer or io.WriterTo",
		format:    "text",
		value:     struct{}{},
		wantErr:   "render: unsupported format: text",
		wantErrIs: []error{render.Err, render.ErrUnsupportedFormat},
	},
	{
		name:      "text format with error writing to writer",
		format:    "text",
		writeErr:  errors.New("write error!!1"),
		value:     &mockStringer{value: "test string"},
		wantErr:   "render: failed: write error!!1",
		wantErrIs: []error{render.Err, render.ErrFailed},
	},
	{
		name:   "text format with error from io.WriterTo",
		format: "text",
		value: &mockWriterTo{
			value: "test string",
			err:   errors.New("WriteTo error!!1"),
		},
		wantErr:   "render: failed: WriteTo error!!1",
		wantErrIs: []error{render.Err, render.ErrFailed},
	},
	{
		name:      "text format with invalid type",
		format:    "text",
		value:     make(chan int),
		wantErr:   "render: unsupported format: text",
		wantErrIs: []error{render.Err, render.ErrUnsupportedFormat},
	},
}

// "xml" format.
var xmlFormatTestCases = []renderFormatTestCase{
	{
		name:   "xml format",
		format: "xml",
		value: struct {
			XMLName xml.Name `xml:"user"`
			Age     int      `xml:"age"`
		}{Age: 30},
		want: "<user>\n  <age>30</age>\n</user>",
	},
	{
		name:   "xml format with xml.Marshaler",
		format: "xml",
		value:  &mockXMLMarshaler{elm: "test string"},
		want:   "<mockXMLMarshaler>test string</mockXMLMarshaler>",
	},
	{
		name:      "xml format with error from xml.Marshaler",
		format:    "xml",
		value:     &mockXMLMarshaler{err: errors.New("marshal error!!1")},
		wantErr:   "render: failed: marshal error!!1",
		wantErrIs: []error{render.Err, render.ErrFailed},
	},
	{
		name:     "xml format with error writing to writer",
		format:   "xml",
		writeErr: errors.New("write error!!1"),
		value: struct {
			XMLName xml.Name `xml:"user"`
			Age     int      `xml:"age"`
		}{Age: 30},
		wantErr:   "render: failed: write error!!1",
		wantErrIs: []error{render.Err, render.ErrFailed},
	},
	{
		name:      "xml format with invalid value",
		format:    "xml",
		value:     make(chan int),
		wantErr:   "render: failed: xml: unsupported type: chan int",
		wantErrIs: []error{render.Err, render.ErrFailed},
	},
}

// "yaml" format.
var yamlFormatTestCases = []renderFormatTestCase{
	{
		name:   "yaml format",
		format: "yaml",
		value:  map[string]int{"age": 30},
		want:   "age: 30\n",
	},
	{
		name:   "yaml format with nested structure",
		format: "yaml",
		value: map[string]any{
			"user": map[string]any{
				"age":  30,
				"name": "John Doe",
			},
		},
		want: "user:\n  age: 30\n  name: John Doe\n",
	},
	{
		name:   "yaml format with yaml.Marshaler",
		format: "yaml",
		value:  &mockYAMLMarshaler{val: map[string]int{"age": 30}},
		want:   "age: 30\n",
	},
	{
		name:      "yaml format with error from yaml.Marshaler",
		format:    "yaml",
		value:     &mockYAMLMarshaler{err: errors.New("mock error")},
		wantErr:   "render: failed: mock error",
		wantErrIs: []error{render.Err, render.ErrFailed},
	},
	{
		name:      "yaml format with error writing to writer",
		format:    "yaml",
		writeErr:  errors.New("write error!!1"),
		value:     map[string]int{"age": 30},
		wantErr:   "render: failed: yaml: write error: write error!!1",
		wantErrIs: []error{render.Err, render.ErrFailed},
	},
	{
		name:      "yaml format with invalid type",
		format:    "yaml",
		value:     make(chan int),
		wantPanic: "cannot marshal type: chan int",
	},
}

func TestRender(t *testing.T) {
	tests := []renderFormatTestCase{}
	tests = append(tests, jsonFormatTestCases...)
	tests = append(tests, textFormatTestCases...)
	tests = append(tests, yamlFormatTestCases...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &mockWriter{WriteErr: tt.writeErr}

			var err error
			var panicRes any
			func() {
				defer func() {
					if r := recover(); r != nil {
						panicRes = r
					}
				}()
				err = render.Render(w, tt.format, tt.value)
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

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		formats   []string
		want      *render.Renderer
		wantErr   string
		wantErrIs []error
	}{
		{
			name:      "no formats",
			formats:   []string{},
			wantErr:   "render: no formats specified",
			wantErrIs: []error{render.Err},
		},
		{
			name:    "single format",
			formats: []string{"json"},
			want: &render.Renderer{
				Formats: map[string]render.FormatRenderer{
					"json": render.DefaultJSON,
				},
			},
		},
		{
			name:    "multiple formats",
			formats: []string{"json", "text", "yaml"},
			want: &render.Renderer{
				Formats: map[string]render.FormatRenderer{
					"json": render.DefaultJSON,
					"text": render.DefaultText,
					"yaml": render.DefaultYAML,
				},
			},
		},
		{
			name:      "invalid format",
			formats:   []string{"json", "text", "invalid"},
			wantErr:   "render: unsupported format: invalid",
			wantErrIs: []error{render.Err, render.ErrUnsupportedFormat},
		},
		{
			name:    "duplicate format",
			formats: []string{"json", "text", "json"},
			want: &render.Renderer{
				Formats: map[string]render.FormatRenderer{
					"json": render.DefaultJSON,
					"text": render.DefaultText,
				},
			},
		},
		{
			name:    "all formats",
			formats: []string{"json", "text", "yaml", "xml", "binary"},
			want: &render.Renderer{
				Formats: map[string]render.FormatRenderer{
					"json":   render.DefaultJSON,
					"text":   render.DefaultText,
					"yaml":   render.DefaultYAML,
					"xml":    render.DefaultXML,
					"binary": render.DefaultBinary,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := render.New(tt.formats...)

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

func TestMustNew(t *testing.T) {
	tests := []struct {
		name      string
		formats   []string
		want      *render.Renderer
		wantErr   string
		wantErrIs []error
		wantPanic string
	}{
		{
			name:      "no formats",
			formats:   []string{},
			wantPanic: "render: no formats specified",
		},
		{
			name:    "single format",
			formats: []string{"json"},
			want: &render.Renderer{
				Formats: map[string]render.FormatRenderer{
					"json": render.DefaultJSON,
				},
			},
		},
		{
			name:    "multiple formats",
			formats: []string{"json", "text", "yaml"},
			want: &render.Renderer{
				Formats: map[string]render.FormatRenderer{
					"json": render.DefaultJSON,
					"text": render.DefaultText,
					"yaml": render.DefaultYAML,
				},
			},
		},
		{
			name:      "invalid format",
			formats:   []string{"json", "text", "invalid"},
			wantPanic: "render: unsupported format: invalid",
		},
		{
			name:    "duplicate format",
			formats: []string{"json", "text", "json"},
			want: &render.Renderer{
				Formats: map[string]render.FormatRenderer{
					"json": render.DefaultJSON,
					"text": render.DefaultText,
				},
			},
		},
		{
			name:    "all formats",
			formats: []string{"json", "text", "yaml", "xml", "binary"},
			want: &render.Renderer{
				Formats: map[string]render.FormatRenderer{
					"json":   render.DefaultJSON,
					"text":   render.DefaultText,
					"yaml":   render.DefaultYAML,
					"xml":    render.DefaultXML,
					"binary": render.DefaultBinary,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *render.Renderer
			var err error
			var panicRes any
			func() {
				defer func() {
					if r := recover(); r != nil {
						panicRes = r
					}
				}()
				got = render.MustNew(tt.formats...)
			}()

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
