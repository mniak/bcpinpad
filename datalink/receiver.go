package datalink

import (
	"bufio"
	"io"
)

type dataReceiver struct {
	scanner *bufio.Scanner
}

func NewDataReceiver(r io.Reader) dataReceiver {
	return dataReceiver{
		scanner: NewScanner(r),
	}
}

func (r dataReceiver) Receive() (string, error) {
	ok := r.scanner.Scan()
	err := r.scanner.Err()
	bytes := r.scanner.Text()
	if err == nil && !ok {
		return string(bytes), io.EOF
	}
	return string(bytes), err
}

func (r dataReceiver) ReadACKorNAK() (bool, error) {
	_, err := r.Receive()
	switch err {
	case ErrACK:
		return true, nil
	case ErrNAK:
		return false, nil
	}
	return false, err
}
