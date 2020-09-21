package raw

import (
	"io"
	"time"
)

func EntangledReadWriters() (io.ReadWriter, io.ReadWriter) {
	ar, aw := io.Pipe()
	br, bw := io.Pipe()
	alice := ConcreteReadWriter{ar, bw}
	bob := ConcreteReadWriter{br, aw}
	return alice, bob
	// var b1 bytes.Buffer
	// var b2 bytes.Buffer
	// return &b1, &b2
}

func Time(fn func()) time.Duration {
	start := time.Now()
	fn()
	duration := time.Since(start)
	return duration
}

type ConcreteReadWriter struct {
	io.Reader
	io.Writer
}

// func ReadChannel(r io.Reader) (chan byte, chan error) {
// 	ch := make(chan byte)
// 	che := make(chan error)
// 	go func(chan byte) {
// 		buffer := make([]byte, 1)
// 		for ch != nil {
// 			_, err := r.Read(buffer)
// 			if err != nil {
// 				che <- err
// 				break
// 			}
// 		}

// 	}(ch)
// 	return ch, che
// }

// func ReadWithTimeout(r io.Reader, bufferSize int) ([]byte, error) {
// 	ch := make(chan []byte)
// 	che := make(chan error)
// 	go func() {
// 		buffer := make([]byte, bufferSize)
// 		_, err := r.Read(buffer)
// 		if err != nil {
// 			che <- err
// 		} else {
// 			buffer.
// 		}
// 	}()
// 	select {
// 	case data := <-ch:
// 		return data, nil
// 	case err := <-che:
// 		return []byte{}, err
// 	}
// }
