package transport

import (
	"io"

	"github.com/mniak/bcpinpad"
	"github.com/mniak/bcpinpad/utils"
)

type rawPinpad struct {
	rw io.ReadWriter
}

func NewPinpad(rw io.ReadWriter) *rawPinpad {
	return &rawPinpad{
		rw: rw,
	}
}

func (pp *rawPinpad) SendData(payload []byte) error {
	bytes := utils.NewBytesBuilder().
		AddByte(bcpinpad.SYN).
		AddBytes(payload).
		AddByte(bcpinpad.ETB, byte(0x77), byte(0x5e)).
		Bytes()

	_, err := pp.rw.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
