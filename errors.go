package bcpinpad

import "errors"

var (
	ErrCancelled = errors.New("cancelled")
	ErrTimeout   = errors.New("timed out")
)
