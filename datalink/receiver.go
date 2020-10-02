package datalink

import (
	"bufio"
	"io"
)

type receiver struct {
	scanner *bufio.Scanner
}

func newSimpleReceiver(r io.Reader) receiver {
	return receiver{
		scanner: NewScanner(r),
	}
}

func (r receiver) Receive() (string, error) {
	ok := r.scanner.Scan()
	err := r.scanner.Err()
	bytes := r.scanner.Text()
	if err == nil && !ok {
		return string(bytes), io.EOF
	}
	return string(bytes), err
}

func (r receiver) ReadACKorNAK() (bool, error) {
	return false, nil
}
