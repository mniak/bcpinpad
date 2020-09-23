package entangled

import (
	"bytes"
	"io"
	"time"
)

const ReadTimeout = 2 * time.Second

type CustomPipe struct {
	buffer *bytes.Buffer
}

func NewPipe() CustomPipe {
	return CustomPipe{
		buffer: bytes.NewBuffer(make([]byte, 0)),
	}
}
func (p CustomPipe) Write(bytes []byte) (int, error) {
	return p.buffer.Write(bytes)
}
func (p CustomPipe) Read(bytes []byte) (int, error) {
	n := 0
	n, err := p.buffer.Read(bytes)
	start := time.Now()
	for err == io.EOF && n < len(bytes) && ReadTimeout > time.Since(start) {
		time.Sleep(400 * time.Millisecond)
		n, err = p.buffer.Read(bytes[n:])
	}
	return n, err
}

func EntangledReadWriters() (ReadWriter, ReadWriter) {
	b1 := NewPipe()
	b2 := NewPipe()

	alice := NewReadWriter(b1, b2)
	bob := NewReadWriter(b2, b1)
	return alice, bob
}

func countTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	duration := time.Since(start)
	return duration
}

type ReadWriter struct {
	io.Reader
	io.Writer
}

func NewReadWriter(r io.Reader, w io.Writer) ReadWriter {
	return ReadWriter{
		r,
		w,
	}
}
