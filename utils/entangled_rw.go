package utils

import (
	"io"
	"time"
)

func EntangledReadWriters() (io.ReadWriter, io.ReadWriter) {
	ar, aw := io.Pipe()
	br, bw := io.Pipe()

	alice := readWriter{ar, bw}
	bob := readWriter{br, aw}
	return alice, bob
}

func countTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	duration := time.Since(start)
	return duration
}

type readWriter struct {
	io.Reader
	io.Writer
}
