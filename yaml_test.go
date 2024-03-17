package render_test

import (
	"bytes"
	"testing"

	"github.com/jimeh/go-render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYAML_Render(t *testing.T) {
	tests := []struct {
		name      string
		indent    int
		value     interface{}
		want      string
		wantPanic string
	}{
		{
			name:  "simple object default indent",
			value: map[string]int{"age": 30},
			want:  "age: 30\n",
		},
		{
			name:   "nested structure",
			indent: 0, // This will use the default indent of 2 spaces
			value: map[string]any{
				"user": map[string]any{
					"age":  30,
					"name": "John Doe",
				},
			},
			want: "user:\n  age: 30\n  name: John Doe\n",
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
			name:      "invalid value",
			indent:    0,
			value:     make(chan int),
			wantPanic: "cannot marshal type: chan int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &render.YAML{
				Indent: tt.indent,
			}

			var buf bytes.Buffer
			var err error
			var panicRes any
			func() {
				defer func() {
					if r := recover(); r != nil {
						panicRes = r
					}
				}()
				err = j.Render(&buf, tt.value)
			}()

			if tt.wantPanic != "" {
				assert.Equal(t, tt.wantPanic, panicRes)
			} else {
				require.NoError(t, err)
				got := buf.String()
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
