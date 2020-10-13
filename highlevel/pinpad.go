package highlevel

import (
	"io"
	"time"

	"github.com/mniak/bcpinpad/datalink"
	"github.com/mniak/bcpinpad/transport"
	"github.com/tarm/serial"
)

type Pinpad struct {
	pp transport.Pinpad
}

func Open(rw io.ReadWriteCloser) (Pinpad, error) {
	tpp := transport.NewPinpad(rw, datalink.NewDataReceiver(rw))
	return Pinpad{
		pp: tpp,
	}, nil
}

func OpenSerial(port string) (Pinpad, error) {
	config := &serial.Config{
		Name:        port,
		Baud:        19200,
		Parity:      serial.ParityNone,
		Size:        8,
		StopBits:    1,
		ReadTimeout: 2 * time.Second,
	}
	s, err := serial.OpenPort(config)
	if err != nil {
		return Pinpad{}, err
	}
	return Open(s)
}
