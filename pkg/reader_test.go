package pkg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdReader_ReadString(t *testing.T) {
	t.Parallel()

	t.Run("ignore same input", func(t *testing.T) {
		t.Parallel()

		stdout := new(bytes.Buffer)
		target := NewBashOutputReader(stdout)

		stdout.Write([]byte("Hello\n"))

		_, line, err := target.ReadString()
		assert.Nil(t, err)
		assert.Equal(t, "Hello", line)

		stdout.Write([]byte("Hello\nWorld"))

		_, line, err = target.ReadString()
		assert.Nil(t, err)
		assert.Equal(t, "World", line)
	})
}
