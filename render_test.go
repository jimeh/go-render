package render

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"
	"testing"

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

var _ FormatRenderer = (*mockRenderer)(nil)

func (m *mockRenderer) Render(w io.Writer, _ any) error {
	_, err := w.Write([]byte(m.output))

	if m.err != nil {
		return m.err
	}

	return err
}

type renderFormatTestCase struct {
	name        string
	writeErr    error
	formats     []string
	value       any
	valueFunc   func() any
	want        string
	wantPretty  string
	wantCompact string
	wantErr     string
	wantErrIs   []error
	wantPanic   string
}

// "binary" format.
var binaryFormattestCases = []renderFormatTestCase{
	{
		name:    "with binary marshaler",
		formats: []string{"binary", "bin"},
		value:   &mockBinaryMarshaler{data: []byte("test string")},
		want:    "test string",
	},
	{
		name:      "without binary marshaler",
		formats:   []string{"binary", "bin"},
		value:     struct{}{},
		wantErr:   "render: unsupported format: {{format}}",
		wantErrIs: []error{Err, ErrUnsupportedFormat},
	},
	{
		name:    "with error marshaling",
		formats: []string{"binary", "bin"},
		value: &mockBinaryMarshaler{
			data: []byte("test string"),
			err:  errors.New("marshal error!!1"),
		},
		wantErr:   "render: failed: marshal error!!1",
		wantErrIs: []error{Err, ErrFailed},
	},
	{
		name:      "with error writing to writer",
		formats:   []string{"binary", "bin"},
		writeErr:  errors.New("write error!!1"),
		value:     &mockBinaryMarshaler{data: []byte("test string")},
		wantErr:   "render: failed: write error!!1",
		wantErrIs: []error{Err, ErrFailed},
	},
	{
		name:      "with invalid type",
		formats:   []string{"binary", "bin"},
		value:     make(chan int),
		wantErr:   "render: unsupported format: {{format}}",
		wantErrIs: []error{Err, ErrUnsupportedFormat},
	},
}

// "json" format.
var jsonFormatTestCases = []renderFormatTestCase{
	{
		name:        "with map",
		formats:     []string{"json"},
		value:       map[string]int{"age": 30},
		wantPretty:  "{\n  \"age\": 30\n}\n",
		wantCompact: "{\"age\":30}\n",
	},
	{
		name:        "with json marshaler",
		formats:     []string{"json"},
		value:       &mockJSONMarshaler{data: []byte(`{"age":30}`)},
		wantPretty:  "{\n  \"age\": 30\n}\n",
		wantCompact: "{\"age\":30}\n",
	},
	{
		name:      "with error from json marshaler",
		formats:   []string{"json"},
		value:     &mockJSONMarshaler{err: errors.New("marshal error!!1")},
		wantErrIs: []error{Err},
	},
	{
		name:      "with error writing to writer",
		formats:   []string{"json"},
		writeErr:  errors.New("write error!!1"),
		value:     map[string]int{"age": 30},
		wantErr:   "render: failed: write error!!1",
		wantErrIs: []error{Err, ErrFailed},
	},
	{
		name:      "with invalid type",
		formats:   []string{"json"},
		value:     make(chan int),
		wantErr:   "render: failed: json: unsupported type: chan int",
		wantErrIs: []error{Err, ErrFailed},
	},
}

// "text" format.
var textFormatTestCases = []renderFormatTestCase{
	{
		name:      "nil",
		formats:   []string{"text", "txt", "plain"},
		value:     nil,
		wantErr:   "render: unsupported format: {{format}}",
		wantErrIs: []error{Err, ErrUnsupportedFormat},
	},
	{
		name:    "byte slice",
		formats: []string{"text", "txt", "plain"},
		value:   []byte("test byte slice"),
		want:    "test byte slice",
	},
	{
		name:    "nil byte slice",
		formats: []string{"text", "txt", "plain"},
		value:   []byte(nil),
		want:    "",
	},
	{
		name:    "empty byte slice",
		formats: []string{"text", "txt", "plain"},
		value:   []byte{},
		want:    "",
	},
	{
		name:    "rune slice",
		formats: []string{"text", "txt", "plain"},
		value:   []rune{'r', 'u', 'n', 'e', 's', '!', ' ', 'y', 'e', 's'},
		want:    "runes! yes",
	},
	{
		name:    "string",
		formats: []string{"text", "txt", "plain"},
		value:   "test string",
		want:    "test string",
	},
	{
		name:    "int",
		formats: []string{"text", "txt", "plain"},
		value:   int(42),
		want:    "42",
	},
	{
		name:    "int8",
		formats: []string{"text", "txt", "plain"},
		value:   int8(43),
		want:    "43",
	},
	{
		name:    "int16",
		formats: []string{"text", "txt", "plain"},
		value:   int16(44),
		want:    "44",
	},
	{
		name:    "int32",
		formats: []string{"text", "txt", "plain"},
		value:   int32(45),
		want:    "45",
	},
	{
		name:    "int64",
		formats: []string{"text", "txt", "plain"},
		value:   int64(46),
		want:    "46",
	},
	{
		name:    "uint",
		formats: []string{"text", "txt", "plain"},
		value:   uint(47),
		want:    "47",
	},
	{
		name:    "uint8",
		formats: []string{"text", "txt", "plain"},
		value:   uint8(48),
		want:    "48",
	},
	{
		name:    "uint16",
		formats: []string{"text", "txt", "plain"},
		value:   uint16(49),
		want:    "49",
	},
	{
		name:    "uint32",
		formats: []string{"text", "txt", "plain"},
		value:   uint32(50),
		want:    "50",
	},
	{
		name:    "uint64",
		formats: []string{"text", "txt", "plain"},
		value:   uint64(51),
		want:    "51",
	},
	{
		name:    "float32",
		formats: []string{"text", "txt", "plain"},
		value:   float32(3.14),
		want:    "3.14",
	},
	{
		name:    "float64",
		formats: []string{"text", "txt", "plain"},
		value:   float64(3.14159),
		want:    "3.14159",
	},
	{
		name:    "bool true",
		formats: []string{"text", "txt", "plain"},
		value:   true,
		want:    "true",
	},
	{
		name:    "bool false",
		formats: []string{"text", "txt", "plain"},
		value:   false,
		want:    "false",
	},
	{
		name:    "implements fmt.Stringer",
		formats: []string{"text", "txt", "plain"},
		value:   &mockStringer{value: "test string"},
		want:    "test string",
	},
	{
		name:      "error writing to writer with fmt.Stringer",
		formats:   []string{"text", "txt", "plain"},
		writeErr:  errors.New("write error!!1"),
		value:     &mockStringer{value: "test string"},
		wantErr:   "render: failed: write error!!1",
		wantErrIs: []error{Err, ErrFailed},
	},
	{
		name:    "implements io.WriterTo",
		formats: []string{"text", "txt", "plain"},
		value:   &mockWriterTo{value: "test string"},
		want:    "test string",
	},
	{
		name:    "io.WriterTo error",
		formats: []string{"text", "txt", "plain"},
		value: &mockWriterTo{
			value: "test string",
			err:   errors.New("WriteTo error!!1"),
		},
		wantErr:   "render: failed: WriteTo error!!1",
		wantErrIs: []error{Err, ErrFailed},
	},
	{
		name:      "implements io.Reader",
		formats:   []string{"text", "txt", "plain"},
		valueFunc: func() any { return &mockReader{value: "reader string"} },
		want:      "reader string",
	},
	{
		name:    "io.Reader error",
		formats: []string{"text", "txt", "plain"},
		value: &mockReader{
			value: "reader string",
			err:   errors.New("Read error!!1"),
		},
		wantErr:   "render: failed: Read error!!1",
		wantErrIs: []error{Err, ErrFailed},
	},
	{
		name:    "error",
		formats: []string{"text", "txt", "plain"},
		value:   errors.New("this is an error"),
		want:    "this is an error",
	},
	{
		name:      "does not implement any supported type/interface",
		formats:   []string{"text", "txt", "plain"},
		value:     struct{}{},
		wantErr:   "render: unsupported format: {{format}}",
		wantErrIs: []error{Err, ErrUnsupportedFormat},
	},
}

// "xml" format.
var xmlFormatTestCases = []renderFormatTestCase{
	{
		name:    "xml format",
		formats: []string{"xml"},
		value: struct {
			XMLName xml.Name `xml:"user"`
			Age     int      `xml:"age"`
		}{Age: 30},
		wantPretty:  "<user>\n  <age>30</age>\n</user>",
		wantCompact: "<user><age>30</age></user>",
	},
	{
		name:    "xml format with xml.Marshaler",
		formats: []string{"xml"},
		value:   &mockXMLMarshaler{elm: "test string"},
		want:    "<mockXMLMarshaler>test string</mockXMLMarshaler>",
	},
	{
		name:      "xml format with error from xml.Marshaler",
		formats:   []string{"xml"},
		value:     &mockXMLMarshaler{err: errors.New("marshal error!!1")},
		wantErr:   "render: failed: marshal error!!1",
		wantErrIs: []error{Err, ErrFailed},
	},
	{
		name:     "xml format with error writing to writer",
		formats:  []string{"xml"},
		writeErr: errors.New("write error!!1"),
		value: struct {
			XMLName xml.Name `xml:"user"`
			Age     int      `xml:"age"`
		}{Age: 30},
		wantErr:   "render: failed: write error!!1",
		wantErrIs: []error{Err, ErrFailed},
	},
	{
		name:      "xml format with invalid value",
		formats:   []string{"xml"},
		value:     make(chan int),
		wantErr:   "render: failed: xml: unsupported type: chan int",
		wantErrIs: []error{Err, ErrFailed},
	},
}

// "yaml" format.
var yamlFormatTestCases = []renderFormatTestCase{
	{
		name:    "yaml format with map",
		formats: []string{"yaml", "yml"},
		value:   map[string]int{"age": 30},
		want:    "age: 30\n",
	},
	{
		name:    "yaml format with nested structure",
		formats: []string{"yaml", "yml"},
		value: map[string]any{
			"user": map[string]any{
				"age":  30,
				"name": "John Doe",
			},
		},
		want: "user:\n  age: 30\n  name: John Doe\n",
	},
	{
		name:    "yaml format with yaml.Marshaler",
		formats: []string{"yaml", "yml"},
		value:   &mockYAMLMarshaler{val: map[string]int{"age": 30}},
		want:    "age: 30\n",
	},
	{
		name:      "yaml format with error from yaml.Marshaler",
		formats:   []string{"yaml", "yml"},
		value:     &mockYAMLMarshaler{err: errors.New("mock error")},
		wantErr:   "render: failed: mock error",
		wantErrIs: []error{Err, ErrFailed},
	},
	{
		name:      "yaml format with error writing to writer",
		formats:   []string{"yaml", "yml"},
		writeErr:  errors.New("write error!!1"),
		value:     map[string]int{"age": 30},
		wantErr:   "render: failed: yaml: write error: write error!!1",
		wantErrIs: []error{Err, ErrFailed},
	},
	{
		name:      "yaml format with invalid type",
		formats:   []string{"yaml", "yml"},
		value:     make(chan int),
		wantPanic: "cannot marshal type: chan int",
	},
}

func TestPretty(t *testing.T) {
	tests := []renderFormatTestCase{}
	tests = append(tests, jsonFormatTestCases...)
	tests = append(tests, textFormatTestCases...)
	tests = append(tests, yamlFormatTestCases...)

	for _, tt := range tests {
		for _, format := range tt.formats {
			t.Run(format+" format "+tt.name, func(t *testing.T) {
				w := &mockWriter{WriteErr: tt.writeErr}

				value := tt.value
				if tt.valueFunc != nil {
					value = tt.valueFunc()
				}

				var err error
				var panicRes any
				func() {
					defer func() {
						if r := recover(); r != nil {
							panicRes = r
						}
					}()
					err = Pretty(w, format, value)
				}()

				got := w.String()
				var want string
				if tt.wantPretty == "" && tt.wantCompact == "" {
					want = tt.want
				} else {
					want = tt.wantPretty
				}

				if tt.wantPanic != "" {
					assert.Equal(t, tt.wantPanic, panicRes)
				}

				if tt.wantErr != "" {
					wantErr := strings.ReplaceAll(
						tt.wantErr, "{{format}}", format,
					)
					assert.EqualError(t, err, wantErr)
				}
				for _, e := range tt.wantErrIs {
					assert.ErrorIs(t, err, e)
				}

				if tt.wantPanic == "" &&
					tt.wantErr == "" && len(tt.wantErrIs) == 0 {
					assert.NoError(t, err)
					assert.Equal(t, want, got)
				}
			})
		}
	}
}

func TestCompact(t *testing.T) {
	tests := []renderFormatTestCase{}
	tests = append(tests, jsonFormatTestCases...)
	tests = append(tests, textFormatTestCases...)
	tests = append(tests, yamlFormatTestCases...)

	for _, tt := range tests {
		for _, format := range tt.formats {
			t.Run(format+" format "+tt.name, func(t *testing.T) {
				w := &mockWriter{WriteErr: tt.writeErr}

				value := tt.value
				if tt.valueFunc != nil {
					value = tt.valueFunc()
				}

				var err error
				var panicRes any
				func() {
					defer func() {
						if r := recover(); r != nil {
							panicRes = r
						}
					}()
					err = Compact(w, format, value)
				}()

				got := w.String()
				var want string
				if tt.wantPretty == "" && tt.wantCompact == "" {
					want = tt.want
				} else {
					want = tt.wantCompact
				}

				if tt.wantPanic != "" {
					assert.Equal(t, tt.wantPanic, panicRes)
				}

				if tt.wantErr != "" {
					wantErr := strings.ReplaceAll(
						tt.wantErr, "{{format}}", format,
					)
					assert.EqualError(t, err, wantErr)
				}
				for _, e := range tt.wantErrIs {
					assert.ErrorIs(t, err, e)
				}

				if tt.wantPanic == "" &&
					tt.wantErr == "" && len(tt.wantErrIs) == 0 {
					assert.NoError(t, err)
					assert.Equal(t, want, got)
				}
			})
		}
	}
}
