package raw

import (
	"bufio"
	"testing"

	"github.com/mniak/ppabecs/lib/utils"
	"github.com/stretchr/testify/assert"
)

func TestScanWellFormattedMessage(t *testing.T) {
	testData := []struct {
		text string
		crc1 byte
		crc2 byte
	}{
		{"OPN000", 0x77, 0x5e},
		{"AAAAAAAA", 0x9a, 0x63},
	}

	for _, d := range testData {
		t.Run(d.text, func(t *testing.T) {
			alice, bob := EntangledReadWriters()

			scanner := bufio.NewScanner(alice)
			scanner.Split(PayloadScanner)

			bytes := utils.NewBytesBuilder().
				AddByte(SYN).
				AddString(d.text).
				AddByte(ETB).
				AddByte(d.crc1, d.crc2).
				Bytes()

			go bob.Write(bytes)
			assert.True(t, scanner.Scan(), "scan failed")
			assert.NoError(t, scanner.Err(), "scan raised error")
			text := scanner.Text()
			assert.Equal(t, d.text, text, "scan text is invalid")
		})
	}
}

func TestScanWellFormattedMessage_WithCANInTheBeginning(t *testing.T) {
	testData := []struct {
		text string
		crc1 byte
		crc2 byte
	}{
		{"OPN000", 0x77, 0x5e},
		{"AAAAAAAA", 0x9a, 0x63},
	}

	for _, d := range testData {
		t.Run(d.text, func(t *testing.T) {
			alice, bob := EntangledReadWriters()

			scanner := bufio.NewScanner(alice)
			scanner.Split(PayloadScanner)

			bytes := utils.NewBytesBuilder().
				AddByte(CAN).
				AddByte(SYN).
				AddString(d.text).
				AddByte(ETB).
				AddByte(d.crc1, d.crc2).
				Bytes()

			go bob.Write(bytes)
			assert.True(t, scanner.Scan(), "scan failed")
			assert.NoError(t, scanner.Err(), "scan raised error")
			text := scanner.Text()
			assert.Equal(t, d.text, text, "scan text is invalid")
		})
	}
}
