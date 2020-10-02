package datalink

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/mniak/bcpinpad"
	"github.com/mniak/bcpinpad/datalink/entangled"
	"github.com/mniak/bcpinpad/transport"
	"github.com/mniak/bcpinpad/utils"
	"github.com/stellar/go/crc16"
	"github.com/stretchr/testify/assert"
)

func TestReceiverShouldBeTransportReceiver(t *testing.T) {
	r := bytes.NewReader([]byte("dummy bytes"))
	rec := NewReceiver(r)
	var a transport.Receiver = rec
	_ = a
}

func TestReceiveWellFormattedMessage(t *testing.T) {
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

			recv := NewReceiver(alice)

			bytes := utils.NewBytesBuilder().
				AddByte(bcpinpad.SYN).
				AddString(d.text).
				AddByte(bcpinpad.ETB).
				AddByte(d.crc1, d.crc2).
				Bytes()

			go func() {
				bob.Write(bytes)
			}()
			text, err := recv.Receive()
			assert.NoError(t, err, "scan raised error")
			assert.Equal(t, d.text, text, "scan text is invalid")
		})
	}
}

func TestReceiveWellFormattedMessage_WithCANInTheBeginning(t *testing.T) {
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

			recv := NewReceiver(alice)

			bytes := utils.NewBytesBuilder().
				AddByte(bcpinpad.CAN).
				AddByte(bcpinpad.SYN).
				AddString(d.text).
				AddByte(bcpinpad.ETB).
				AddByte(d.crc1, d.crc2).
				Bytes()

			go bob.Write(bytes)
			text, err := recv.Receive()
			assert.NoError(t, err, "scan raised error")
			assert.Equal(t, d.text, text, "scan text is invalid")
		})
	}
}

func TestReceiveWithoutData(t *testing.T) {
	alice, _ := entangled.EntangledReadWriters()

	recv := NewReceiver(alice)

	text, err := recv.Receive()
	assert.Error(t, err, "scan should raise error")
	assert.Equal(t, err, io.EOF, "scan error should be EOF")
	assert.Empty(t, text, "scan text should be empty")
}

func TestReceiveWithWrongCRC(t *testing.T) {

	alice, bob := entangled.EntangledReadWriters()

	recv := NewReceiver(alice)

	bytes := utils.NewBytesBuilder().
		AddByte(bcpinpad.SYN).
		AddString("ABCDEFG").
		AddByte(bcpinpad.ETB).
		AddByte(0x11, 0x22).
		Bytes()

	go bob.Write(bytes)
	_, err := recv.Receive()
	assert.Error(t, err, "scan should raise error")
	assert.Equal(t, err, crc16.ErrInvalidChecksum, "scan error should be due to CRC")
}

func TestReceiveWithByteOutOfRange(t *testing.T) {
	testData := []byte{
		0x00, 0x11, 0x19,
		0x90, 0xa0, 0xf0,
	}
	for _, b := range testData {
		t.Run(fmt.Sprintf("%x", b), func(t *testing.T) {
			alice, bob := entangled.EntangledReadWriters()

			recv := NewReceiver(alice)

			bytes := utils.NewBytesBuilder().
				AddByte(bcpinpad.SYN).
				AddString("ABCD").
				AddByte(b).
				AddString("EFGH").
				AddByte(bcpinpad.ETB).
				AddByte(0x11, 0x22).
				Bytes()

			go bob.Write(bytes)
			_, err := recv.Receive()
			assert.Error(t, err, "scan should raise error")
			assert.Equal(t, err, ErrBytesOutOfRange, "scan error should be due to bytes out of range")
		})
	}
}

func TestReceiveWithPayloadLength0(t *testing.T) {

	alice, bob := entangled.EntangledReadWriters()

	recv := NewReceiver(alice)

	bytes := utils.NewBytesBuilder().
		AddByte(bcpinpad.SYN).
		AddByte(bcpinpad.ETB).
		AddByte(0x11, 0x22).
		Bytes()

	go bob.Write(bytes)
	_, err := recv.Receive()
	assert.Error(t, err, "scan should raise error")
	assert.Equal(t, err, ErrMessageTooShort, "scan error should be due to payload too short")
}

func TestReceiveWithPayloadLengthGreaterThan1024(t *testing.T) {

	alice, bob := entangled.EntangledReadWriters()

	recv := NewReceiver(alice)

	bytes := utils.NewBytesBuilder().
		AddByte(bcpinpad.SYN).
		AddString(strings.Repeat("a", 1024+1)).
		AddByte(bcpinpad.ETB).
		AddByte(0x11, 0x22).
		Bytes()

	go bob.Write(bytes)
	_, err := recv.Receive()
	assert.Error(t, err, "scan should raise error")
	assert.Equal(t, err, ErrMessageTooLong, "scan error should be due to payload too long")
}

func TestReceiveACKorNAK_WhenReadByte_ShouldReturnAccordingly(t *testing.T) {
	testData := map[string]struct {
		byte
		ack bool
		err error
	}{
		"ACK": {bcpinpad.ACK, true, nil},
		"NAK": {bcpinpad.NAK, false, nil},
		"ETB": {bcpinpad.ETB, false, ErrProtocolViolation},
		"X":   {'X', false, ErrProtocolViolation},
		"8":   {8, false, ErrProtocolViolation},
	}
	for name, d := range testData {
		t.Run(name, func(t *testing.T) {
			alice, bob := entangled.EntangledReadWriters()
			recv := NewReceiver(alice)

			go bob.Write([]byte{d.byte})
			ack, err := recv.ReadACKorNAK()
			assert.Equal(t, d.ack, ack)
			assert.Equal(t, d.err, err)
		})
	}
}
