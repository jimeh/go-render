package render_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
)

type mockStringer struct {
	value string
}

var _ fmt.Stringer = (*mockStringer)(nil)

func (ms *mockStringer) String() string {
	return ms.value
}

type mockWriterTo struct {
	value string
	err   error
}

var _ io.WriterTo = (*mockWriterTo)(nil)

func (m *mockWriterTo) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(m.value))

	if m.err != nil {
		return int64(n), m.err
	}

	return int64(n), err
}

type mockReader struct {
	value  string
	cursor int
	err    error
}

var _ io.Reader = (*mockReader)(nil)

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.err != nil {
		return 0, m.err
	}

	if len(m.value) == 0 {
		return 0, io.EOF
	}

	n = copy(p, m.value[m.cursor:])
	m.cursor += n

	if m.cursor >= len(m.value) {
		return n, io.EOF
	}

	return n, nil
}

func Test_mockReader_Read(t *testing.T) {
	mr := &mockReader{value: "test string"}

	b1 := make([]byte, 5)
	n1, err := mr.Read(b1)

	assert.NoError(t, err)
	assert.Equal(t, 5, n1)
	assert.Equal(t, "test ", string(b1))

	b2 := make([]byte, 5)
	n2, err := mr.Read(b2)

	assert.NoError(t, err)
	assert.Equal(t, 5, n2)
	assert.Equal(t, "strin", string(b2))

	b3 := make([]byte, 5)
	n3, err := mr.Read(b3)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 1, n3)
	assert.Equal(t, []byte{byte('g'), 0, 0, 0, 0}, b3)
}

func TestText_Render(t *testing.T) {
	tests := []struct {
		name      string
		writeErr  error
		value     any
		want      string
		wantErr   string
		wantErrIs []error
	}{
		{
			name:      "nil",
			value:     nil,
			wantErr:   "render: cannot render: <nil>",
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
		{
			name:  "byte slice",
			value: []byte("test byte slice"),
			want:  "test byte slice",
		},
		{
			name:  "nil byte slice",
			value: []byte(nil),
			want:  "",
		},
		{
			name:  "empty byte slice",
			value: []byte{},
			want:  "",
		},
		{
			name:  "rune slice",
			value: []rune{'r', 'u', 'n', 'e', 's', '!', ' ', 'y', 'e', 's'},
			want:  "runes! yes",
		},
		{
			name:  "string",
			value: "test string",
			want:  "test string",
		},
		{name: "int", value: int(42), want: "42"},
		{name: "int8", value: int8(43), want: "43"},
		{name: "int16", value: int16(44), want: "44"},
		{name: "int32", value: int32(45), want: "45"},
		{name: "int64", value: int64(46), want: "46"},
		{name: "uint", value: uint(47), want: "47"},
		{name: "uint8", value: uint8(48), want: "48"},
		{name: "uint16", value: uint16(49), want: "49"},
		{name: "uint32", value: uint32(50), want: "50"},
		{name: "uint64", value: uint64(51), want: "51"},
		{name: "float32", value: float32(3.14), want: "3.14"},
		{name: "float64", value: float64(3.14159), want: "3.14159"},
		{name: "bool true", value: true, want: "true"},
		{name: "bool false", value: false, want: "false"},
		{
			name:  "implements fmt.Stringer",
			value: &mockStringer{value: "test string"},
			want:  "test string",
		},
		{
			name:      "error writing to writer with fmt.Stringer",
			writeErr:  errors.New("write error!!1"),
			value:     &mockStringer{value: "test string"},
			wantErr:   "render: failed: write error!!1",
			wantErrIs: []error{render.Err, render.ErrFailed},
		},
		{
			name:  "implements io.WriterTo",
			value: &mockWriterTo{value: "test string"},
			want:  "test string",
		},
		{
			name: "io.WriterTo error",
			value: &mockWriterTo{
				value: "test string",
				err:   errors.New("WriteTo error!!1"),
			},
			wantErr:   "render: failed: WriteTo error!!1",
			wantErrIs: []error{render.Err, render.ErrFailed},
		},
		{
			name:  "implements io.Reader",
			value: &mockReader{value: "reader string"},
			want:  "reader string",
		},
		{
			name: "io.Reader error",
			value: &mockReader{
				value: "reader string",
				err:   errors.New("Read error!!1"),
			},
			wantErr:   "render: failed: Read error!!1",
			wantErrIs: []error{render.Err, render.ErrFailed},
		},
		{
			name:  "error",
			value: errors.New("this is an error"),
			want:  "this is an error",
		},
		{
			name:      "does not implement any supported type/interface",
			value:     struct{}{},
			wantErr:   "render: cannot render: struct {}",
			wantErrIs: []error{render.Err, render.ErrCannotRender},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &render.Text{}
			w := &mockWriter{WriteErr: tt.writeErr}

			err := s.Render(w, tt.value)
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

func ptr[T any](v T) *T {
	return &v
}
