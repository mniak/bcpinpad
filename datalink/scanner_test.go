package datalink

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/mniak/bcpinpad"
	"github.com/mniak/bcpinpad/datalink/entangled"
	"github.com/mniak/bcpinpad/utils"
	"github.com/stellar/go/crc16"
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
			alice, bob := entangled.EntangledReadWriters()

			scanner := bufio.NewScanner(alice)
			scanner.Split(PayloadSplitter)

			bytes := utils.NewBytesBuilder().
				AddByte(bcpinpad.SYN).
				AddString(d.text).
				AddByte(bcpinpad.ETB).
				AddByte(d.crc1, d.crc2).
				Bytes()

			go func() {
				bob.Write(bytes)
			}()
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
			alice, bob := entangled.EntangledReadWriters()

			scanner := bufio.NewScanner(alice)
			scanner.Split(PayloadSplitter)

			bytes := utils.NewBytesBuilder().
				AddByte(bcpinpad.CAN).
				AddByte(bcpinpad.SYN).
				AddString(d.text).
				AddByte(bcpinpad.ETB).
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

func TestScanWithoutData(t *testing.T) {
	alice, _ := entangled.EntangledReadWriters()

	scanner := bufio.NewScanner(alice)
	scanner.Split(PayloadSplitter)

	assert.False(t, scanner.Scan(), "scan should fail")
	assert.NoError(t, scanner.Err(), "scan should not raise error")
}

func TestScanWithWrongCRC(t *testing.T) {

	alice, bob := entangled.EntangledReadWriters()

	scanner := bufio.NewScanner(alice)
	scanner.Split(PayloadSplitter)

	bytes := utils.NewBytesBuilder().
		AddByte(bcpinpad.SYN).
		AddString("ABCDEFG").
		AddByte(bcpinpad.ETB).
		AddByte(0x11, 0x22).
		Bytes()

	go bob.Write(bytes)
	assert.False(t, scanner.Scan(), "scan should fail")
	assert.Error(t, scanner.Err(), "scan should raise error")
	assert.Equal(t, scanner.Err(), crc16.ErrInvalidChecksum, "scan error should be due to CRC")
}

func TestScanWithByteOutOfRange(t *testing.T) {
	testData := []byte{
		0x00, 0x11, 0x19,
		0x90, 0xa0, 0xf0,
	}
	for _, b := range testData {
		t.Run(fmt.Sprintf("%x", b), func(t *testing.T) {
			alice, bob := entangled.EntangledReadWriters()

			scanner := bufio.NewScanner(alice)
			scanner.Split(PayloadSplitter)

			bytes := utils.NewBytesBuilder().
				AddByte(bcpinpad.SYN).
				AddString("ABCD").
				AddByte(b).
				AddString("EFGH").
				AddByte(bcpinpad.ETB).
				AddByte(0x11, 0x22).
				Bytes()

			go bob.Write(bytes)
			assert.False(t, scanner.Scan(), "scan should fail")
			assert.Error(t, scanner.Err(), "scan should raise error")
			assert.Equal(t, scanner.Err(), ErrBytesOutOfRange, "scan error should be due to bytes out of range")
		})
	}
}

func TestScanWithPayloadLength0(t *testing.T) {

	alice, bob := entangled.EntangledReadWriters()

	scanner := bufio.NewScanner(alice)
	scanner.Split(PayloadSplitter)

	bytes := utils.NewBytesBuilder().
		AddByte(bcpinpad.SYN).
		AddByte(bcpinpad.ETB).
		AddByte(0x11, 0x22).
		Bytes()

	go bob.Write(bytes)
	assert.False(t, scanner.Scan(), "scan should fail")
	assert.Error(t, scanner.Err(), "scan should raise error")
	assert.Equal(t, scanner.Err(), ErrMessageTooShort, "scan error should be due to payload too short")
}

func TestScanWithPayloadLengthGreaterThan1024(t *testing.T) {

	alice, bob := entangled.EntangledReadWriters()

	scanner := bufio.NewScanner(alice)
	scanner.Split(PayloadSplitter)

	bytes := utils.NewBytesBuilder().
		AddByte(bcpinpad.SYN).
		AddString(strings.Repeat("a", 1024+1)).
		AddByte(bcpinpad.ETB).
		AddByte(0x11, 0x22).
		Bytes()

	go bob.Write(bytes)
	assert.False(t, scanner.Scan(), "scan should fail")
	assert.Error(t, scanner.Err(), "scan should raise error")
	assert.Equal(t, scanner.Err(), ErrMessageTooLong, "scan error should be due to payload too long")
}
