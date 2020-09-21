package raw

import (
	"bytes"
	"fmt"
	"io"

	"github.com/stellar/go/crc16"
)

func PayloadScanner(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		// Request more data.
		return 0, nil, nil
	}

	if data[0] == CAN {
		advance, token, err := PayloadScanner(data[1:], atEOF)
		return advance + 1, token, err
	}

	if data[0] != SYN {
		return 0, nil, fmt.Errorf("protocol violation. expecting SYN (0x16) but received %x", data[0])
	}
	if i := bytes.IndexByte(data[1:], ETB); i >= 0 {
		payload := data[1 : i+1]

		// If data does not have space for CRC
		totalLength := 1 + len(payload) + 1 + 2
		if len(data) < totalLength {
			// Request more data.
			return 0, nil, nil
		}

		// Validate byte range
		for _, b := range payload {
			if b < 0x20 || b > 0x7f {
				return 0, nil, fmt.Errorf("protocol violation. expecting chars in the range 0x20-0x7f but received %x", data[0])
			}
		}

		crc := data[i+2 : i+4]

		//Validate CRC (bytes must be reversed)
		crcErr := crc16.Validate(data[1:i+2], []byte{crc[1], crc[0]})

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
