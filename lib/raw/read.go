package raw

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/mniak/ppabecs/lib/utils"
)

type Response interface{}

// type DataResponse struct{}
// type InterruptResponse struct{}
// func (r io.Reader) ReadMessage() Response {
// 	br := bufio.NewReader(r)
// 	br.ReadByte()
// }
// type customReader struct {
// 	r  io.Reader
// 	br bufio.Reader
// }
// func (cr customReader) start() {
// 	go func() {
// 		buff := make([]byte, 1024)
// 		for {
// 			cr.r.Read(buff)
// 			// cr.buff.
// 		}
// 	}()
// }
const Timeout = 2 * time.Second

var (
	Aborted = errors.New("aborted")
)

func ReadByte(br bufio.Reader) (byte, error) {
	ctx, _ := context.WithTimeout(context.Background(), Timeout)
	result, err := utils.RunWithContext(ctx, func() (interface{}, error) {
		for {
			b, err := br.ReadByte()
			if err == io.EOF {
				continue
			}
			if err != nil {
				return nil, err
			}
			return b, nil
		}
	})
	if err == nil {
		return result.(byte), err
	}
	return byte(0), err
}
func ReadSYN(r bufio.Reader) error {
	b, err := ReadByte(r)
	if err != nil {
		return err
	} else if b != SYN {
		return fmt.Errorf("protocol violation. expecting SYN (0x16) but received %x", b)
	}
	return nil
}

// func scanETB(data []byte, atEOF bool) (advance int, token []byte, err error) {
// 	// Skip leading spaces.
// 	start := 0
// 	for width := 0; start < len(data); start += width {
// 		var r rune
// 		r, width = utf8.DecodeRune(data[start:])
// 		if !isSpace(r) {
// 			break
// 		}
// 	}
// 	// Scan until space, marking end of word.
// 	for width, i := 0, start; i < len(data); i += width {
// 		var r rune
// 		r, width = utf8.DecodeRune(data[i:])
// 		if isSpace(r) {
// 			return i + width, data[start:i], nil
// 		}
// 	}
// 	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
// 	if atEOF && len(data) > start {
// 		return len(data), data[start:], nil
// 	}
// 	// Request more data.
// 	return start, nil, nil
// }
// func ReadPayload(br bufio.Reader) (byte, error) {
// 	ctx, _ := context.WithTimeout(context.Background(), Timeout)
// 	result, err := utils.RunWithContext(ctx, func() (interface{}, error) {
// 		for {
// 			bufio.ScanLines()
// 			if err == io.EOF {
// 				continue
// 			}
// 			if err != nil {
// 				return nil, err
// 			}
// 			return b, nil
// 		}
// 	})
// 	if err == nil {
// 		return result.(byte), err
// 	}
// 	return byte(0), err
// }

// func ReadPayload(r bufio.Reader) ([]byte, error) {
// 	ctx, _ := context.WithTimeout(context.Background(), Timeout)
// 	result, err := utils.RunWithContext(ctx, func() (interface{}, error) {
// 		buffer := make([]byte, 1)
// 		n, err := r.Read(buffer)
// 		if err != nil {
// 			return nil, err
// 		} else if n < 1 {
// 			return nil, errors.New("no byte was read")
// 		} else {
// 			return buffer[0], nil
// 		}
// 	})
// 	if err == nil {
// 		return result.([]byte), err
// 	}
// 	return byte(0), err
// }
// func ReadResponse(r bufio.Reader) (*Response, error) {
// 	err := ReadSYN(r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	payload, err := ReadPayload(r)
// 	if err != nil {
// 		return nil, err
// 	}
// }
