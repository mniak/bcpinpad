package datalink

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/mniak/bcpinpad"
	"github.com/stellar/go/crc16"
)

var (
	ErrProtocolViolation = errors.New("protocol violation. expecting byte SYN (0x16), ACK (0x06) or NAK (0x15)")
	ErrBytesOutOfRange   = errors.New("found bytes out of the range 0x20-0x7f")
	ErrMessageTooShort   = errors.New("the message payload is too short. it should be at least 1 in length")
	ErrMessageTooLong    = errors.New("the message payload is too long. it should be at most 1024 in length")
	ErrACK               = errors.New("ACK received")
	ErrNAK               = errors.New("NAK received")
)

func PayloadSplitter(data []byte, atEOF bool) (int, []byte, error) {
	if len(data) == 0 {
		if atEOF {
			return 0, nil, io.EOF
		} else {
			// Request more data.
			return 0, nil, nil
		}
	}

	if data[0] == bcpinpad.CAN {
		advance, token, err := PayloadSplitter(data[1:], atEOF)
		return advance + 1, token, err
	}

	if data[0] == bcpinpad.ACK {
		return 1, []byte{bcpinpad.ACK}, ErrACK
	}
	if data[0] == bcpinpad.NAK {
		return 1, nil, ErrNAK
	}

	if data[0] != bcpinpad.SYN {
		return 0, nil, ErrProtocolViolation
	}
	if i := bytes.IndexByte(data, bcpinpad.ETB); i >= 0 {
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

func NewScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Split(PayloadSplitter)
	return scanner
}
