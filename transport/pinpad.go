package transport

import (
	"errors"
	"io"

	"github.com/mniak/bcpinpad"
	"github.com/mniak/bcpinpad/utils"
	"github.com/stellar/go/crc16"
)

var (
	NegativeAcknowledgement = errors.New("negative acknowledgement received")
)

type Pinpad interface {
	Send(payload []byte) ([]byte, error)
}

type pinpad struct {
	rw       io.ReadWriter
	receiver DataReceiver
}

func NewPinpad(rw io.ReadWriter, receiver DataReceiver) *pinpad {
	return &pinpad{
		rw:       rw,
		receiver: receiver,
	}
}

func (pp *pinpad) sendData(payload []byte) (bool, error) {

	crc := crc16.Checksum(append(payload, bcpinpad.ETB))
	bytes := utils.NewBytesBuilder().
		AddByte(bcpinpad.SYN).
		AddBytes(payload).
		AddByte(bcpinpad.ETB).
		AddByte(crc[1], crc[0]).
		Bytes()

	_, err := pp.rw.Write(bytes)
	if err != nil {
		return false, err
	}

	var ack bool
	for i := 0; i < 3; i++ {
		ack, err = pp.receiver.ReadACKorNAK()
		if err != bcpinpad.ErrTimeout {
			break
		}
	}
	if err != nil {
		pp.rw.Write([]byte{bcpinpad.CAN})
		return false, err
	}
	return ack, nil
}

func (pp *pinpad) Send(payload []byte) ([]byte, error) {
	ok, err := pp.sendData(payload)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, NegativeAcknowledgement
	}

	recvbytes, err := pp.receiver.Receive()
	// if err != nil {
	// 	return nil, err
	// }
	return recvbytes, err
}
