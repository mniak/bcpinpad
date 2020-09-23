package transport

import "errors"

type Sender interface {
}

type Receiver interface {
	Receive() (string, error)
}

var (
	ErrCancelled = errors.New("aborted")
	ErrTimeout   = errors.New("aborted")
)
