package gw

import (
	"bytes"
	"testing"

	"gotest.tools/assert"
)

func TestCommandsReader(t *testing.T) {
	for _, tt := range []struct {
		name     string
		commands []string
		err      string
	}{
		{
			name: "base",
			commands: []string{
				"INFO {}\r\n",
				"MSG test 1 3\r\n123\r\n",
				"MSG test 1 3\r\n1\r\n\r\n",
				"PUB test 3\r\n1\r\n\r\n",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			for _, s := range tt.commands {
				buf.WriteString(s)
			}
			reader := NewCommandsReader(&buf)
			for _, expected := range tt.commands {
				next, err := reader.nextCommand()
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, expected, string(next))
			}
			_, err := reader.nextCommand()
			if err == nil {
				t.Fatal("Expected an error")
			}
			if tt.err != "" {
				assert.Equal(t, tt.err, err.Error())
			} else {
				assert.Equal(t, "EOF", err.Error())
			}
		})
	}
}
