package entangled

import (
	"bufio"
	"io"
	"time"
)

func EntangledReadWriters() (ReadWriter, ReadWriter) {
	ar, aw := io.Pipe()
	br, bw := io.Pipe()

	alice := NewReadWriter(ar, bw)
	bob := NewReadWriter(br, aw)
	return alice, bob
}

func countTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	duration := time.Since(start)
	return duration
}

type ReadWriter struct {
	*bufio.ReadWriter
}

func NewReadWriter(r io.Reader, w io.Writer) ReadWriter {
	return ReadWriter{
		bufio.NewReadWriter(
			bufio.NewReader(r),
			bufio.NewWriter(w),
		),
	}
}

func (rw ReadWriter) Read(p []byte) (int, error) {
	return 0, nil
}
