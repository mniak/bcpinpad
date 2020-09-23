package utils

import (
	"fmt"
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntangledReadWritersOneShouldReadWhatTheOtherWrites(t *testing.T) {
	alice, bob := EntangledReadWriters()

	data := []struct {
		writer io.ReadWriter
		reader io.ReadWriter
	}{
		{alice, bob},
		{bob, alice},
	}
	for i, d := range data {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			buffer1 := make([]byte, 16)
			_, _ = rand.Read(buffer1)
			go func() {
				d.writer.Write(buffer1)
			}()
			go func() {
				buffer2 := make([]byte, 16)
				d.reader.Read(buffer2)
				assert.EqualValues(t, buffer1, buffer2)
			}()
		})
	}
}
