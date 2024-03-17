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

var _ render.Renderer = (*mockRenderer)(nil)

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

func TestDefaultWriterTo(t *testing.T) {
	assert.Equal(t, &render.WriterTo{}, render.DefaultWriterTo)
}

func TestDefaultStringer(t *testing.T) {
	assert.Equal(t, &render.Stringer{}, render.DefaultStringer)
}

func TestDefaultText(t *testing.T) {
	want := &render.MultiRenderer{
		Renderers: []render.Renderer{
			&render.Stringer{},
			&render.WriterTo{},
		},
	}

	assert.Equal(t, want, render.DefaultText)
}

func TestDefaultBinary(t *testing.T) {
	assert.Equal(t, &render.Binary{}, render.DefaultBinary)
}

func TestDefaultRenderer(t *testing.T) {
	want := &render.FormatRenderer{
		Renderers: map[string]render.Renderer{
			"bin":    render.DefaultBinary,
			"binary": render.DefaultBinary,
			"json":   render.DefaultJSON,
			"plain":  render.DefaultText,
			"text":   render.DefaultText,
			"txt":    render.DefaultText,
			"xml":    render.DefaultXML,
			"yaml":   render.DefaultYAML,
			"yml":    render.DefaultYAML,
		},
	}

	assert.Equal(t, want, render.DefaultRenderer)
}

func TestRender(t *testing.T) {
	tests := []struct {
		name      string
		writeErr  error
		format    string
		value     any
		want      string
		wantErr   string
		wantErrIs []error
		wantPanic string
	}{
		// "bin" format.
		{
			name:   "bin format with binary marshaler",
			format: "bin",
			value:  &mockBinaryMarshaler{data: []byte("test string")},
			want:   "test string",
		},
		{
			name:      "bin format without binary marshaler",
			format:    "bin",
			value:     struct{}{},
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		{
			name:   "bin format with error marshaling",
			format: "bin",
			value: &mockBinaryMarshaler{
				data: []byte("test string"),
				err:  errors.New("marshal error!!1"),
			},
			wantErr:   "render: marshal error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "bin format with error writing to writer",
			format:    "bin",
			writeErr:  errors.New("write error!!1"),
			value:     &mockBinaryMarshaler{data: []byte("test string")},
			wantErr:   "render: write error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "bin format with invalid type",
			format:    "bin",
			value:     make(chan int),
			wantErr:   "render: cannot render",
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		// "binary" format.
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
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		{
			name:   "binary format with error marshaling",
			format: "binary",
			value: &mockBinaryMarshaler{
				data: []byte("test string"),
				err:  errors.New("marshal error!!1"),
			},
			wantErr:   "render: marshal error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "binary format with error writing to writer",
			format:    "binary",
			writeErr:  errors.New("write error!!1"),
			value:     &mockBinaryMarshaler{data: []byte("test string")},
			wantErr:   "render: write error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "binary format with invalid type",
			format:    "binary",
			value:     make(chan int),
			wantErr:   "render: cannot render",
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		// "json" format.
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
			wantErr:   "render: write error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "json format with invalid type",
			format:    "json",
			value:     make(chan int),
			wantErr:   "render: json: unsupported type: chan int",
			wantErrIs: []error{render.Err},
		},
		// "plain" format.
		{
			name:   "plain format with fmt.Stringer",
			format: "plain",
			value:  &mockStringer{value: "test string"},
			want:   "test string",
		},
		{
			name:   "plain format with io.WriterTo",
			format: "plain",
			value:  &mockWriterTo{value: "test string"},
			want:   "test string",
		},
		{
			name:      "plain format without fmt.Stringer or io.WriterTo",
			format:    "plain",
			value:     struct{}{},
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		{
			name:      "plain format with error writing to writer",
			format:    "plain",
			writeErr:  errors.New("write error!!1"),
			value:     &mockStringer{value: "test string"},
			wantErr:   "render: write error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:   "plain format with error from io.WriterTo",
			format: "plain",
			value: &mockWriterTo{
				value: "test string",
				err:   errors.New("WriteTo error!!1"),
			},
			wantErr:   "render: WriteTo error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "plain format with invalid type",
			format:    "plain",
			value:     make(chan int),
			wantErr:   "render: cannot render",
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		// "text" format.
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
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		{
			name:      "text format with error writing to writer",
			format:    "text",
			writeErr:  errors.New("write error!!1"),
			value:     &mockStringer{value: "test string"},
			wantErr:   "render: write error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:   "text format with error from io.WriterTo",
			format: "text",
			value: &mockWriterTo{
				value: "test string",
				err:   errors.New("WriteTo error!!1"),
			},
			wantErr:   "render: WriteTo error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "text format with invalid type",
			format:    "text",
			value:     make(chan int),
			wantErr:   "render: cannot render",
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		// "txt" format.
		{
			name:   "txt format with fmt.Stringer",
			format: "txt",
			value:  &mockStringer{value: "test string"},
			want:   "test string",
		},
		{
			name:   "txt format with io.WriterTo",
			format: "txt",
			value:  &mockWriterTo{value: "test string"},
			want:   "test string",
		},
		{
			name:      "txt format without fmt.Stringer or io.WriterTo",
			format:    "txt",
			value:     struct{}{},
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		{
			name:      "txt format with error writing to writer",
			format:    "txt",
			writeErr:  errors.New("write error!!1"),
			value:     &mockStringer{value: "test string"},
			wantErr:   "render: write error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:   "txt format with error from io.WriterTo",
			format: "txt",
			value: &mockWriterTo{
				value: "test string",
				err:   errors.New("WriteTo error!!1"),
			},
			wantErr:   "render: WriteTo error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "txt format with invalid type",
			format:    "txt",
			value:     make(chan int),
			wantErr:   "render: cannot render",
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		// "xml" format.
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
			wantErr:   "render: marshal error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:     "xml format with error writing to writer",
			format:   "xml",
			writeErr: errors.New("write error!!1"),
			value: struct {
				XMLName xml.Name `xml:"user"`
				Age     int      `xml:"age"`
			}{Age: 30},
			wantErr:   "render: write error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "xml format with invalid value",
			format:    "xml",
			value:     make(chan int),
			wantErr:   "render: xml: unsupported type: chan int",
			wantErrIs: []error{render.Err},
		},
		// "yaml" format.
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
			wantErr:   "render: mock error",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "yaml format with error writing to writer",
			format:    "yaml",
			writeErr:  errors.New("write error!!1"),
			value:     map[string]int{"age": 30},
			wantErr:   "render: yaml: write error: write error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "yaml format with invalid type",
			format:    "yaml",
			value:     make(chan int),
			wantPanic: "cannot marshal type: chan int",
		},
		// "yml" format.
		{
			name:   "yml format",
			format: "yml",
			value:  map[string]int{"age": 30},
			want:   "age: 30\n",
		},
		{
			name:   "yml format with nested structure",
			format: "yml",
			value: map[string]any{
				"user": map[string]any{
					"age":  30,
					"name": "John Doe",
				},
			},
			want: "user:\n  age: 30\n  name: John Doe\n",
		},
		{
			name:   "yml format with yaml.Marshaler",
			format: "yml",
			value:  &mockYAMLMarshaler{val: map[string]int{"age": 30}},
			want:   "age: 30\n",
		},
		{
			name:      "yml format with error from yaml.Marshaler",
			format:    "yml",
			value:     &mockYAMLMarshaler{err: errors.New("mock error")},
			wantErr:   "render: mock error",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "yml format with error writing to writer",
			format:    "yml",
			writeErr:  errors.New("write error!!1"),
			value:     map[string]int{"age": 30},
			wantErr:   "render: yaml: write error: write error!!1",
			wantErrIs: []error{render.Err},
		},
		{
			name:      "yml format with invalid type",
			format:    "yml",
			value:     make(chan int),
			wantPanic: "cannot marshal type: chan int",
		},
	}
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
