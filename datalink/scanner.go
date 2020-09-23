package datalink

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/mniak/ppabecs"
	"github.com/stellar/go/crc16"
)

var (
	ErrBytesOutOfRange = errors.New("found bytes out of the range 0x20-0x7f")
	ErrMessageTooShort = errors.New("the message payload is too short. it should be at least 1 in length")
	ErrMessageTooLong  = errors.New("the message payload is too long. it should be at most 1024 in length")
)

func PayloadSplit(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		// Request more data.
		return 0, nil, nil
	}

	if data[0] == ppabecs.CAN {
		advance, token, err := PayloadSplit(data[1:], atEOF)
		return advance + 1, token, err
	}

	if data[0] != ppabecs.SYN {
		return 0, nil, fmt.Errorf("protocol violation. expecting SYN (0x16) but received %x", data[0])
	}
	if i := bytes.IndexByte(data, ppabecs.ETB); i >= 0 {
		if i < 2 {
			return 0, nil, ErrMessageTooShort
		}
		if i > 1024+1 {
			return 0, nil, ErrMessageTooLong
		}
		payload := data[1:i]

		// If data does not have space for CRC
		totalLength := 1 + len(payload) + 1 + 2
		if len(data) < totalLength {
			// Request more data.
			return 0, nil, nil
		}

		// Validate byte range
		for _, b := range payload {
			if b < 0x20 || b > 0x7f {
				return 0, nil, ErrBytesOutOfRange
			}
		}

		crc := data[i+1 : i+1+2]

		//Validate CRC (bytes must be reversed)
		crcErr := crc16.Validate(data[1:i+1], []byte{crc[1], crc[0]})

		// We have a full terminated line.
		return totalLength, payload, crcErr
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		// return len(data), data, nil
		return 0, nil, io.EOF
	}
	// Request more data.
	return 0, nil, nil
}
