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
	_, err := r.Receive()
	switch err {
	case ErrACK:
		return true, nil
	case ErrNAK:
		return false, nil
	}
	return false, err
}
