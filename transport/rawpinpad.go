package transport

import (
	"io"
)

type rawPinpad struct {
	rw io.ReadWriter
}

func NewPinpad(rw io.ReadWriter) *rawPinpad {
	return &rawPinpad{
		rw: rw,
	}
}

func (pp *rawPinpad) SendData(bytes []byte) error {
	_, err := pp.rw.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
