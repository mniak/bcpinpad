package datalink

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/mniak/bcpinpad"
	"github.com/mniak/bcpinpad/utils"
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
	} else if b != bcpinpad.SYN {
		return fmt.Errorf("protocol violation. expecting SYN (0x16) but received %x", b)
	}
	return nil
}
