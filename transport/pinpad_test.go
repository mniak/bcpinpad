package transport

import (
	"sync"
	"testing"
	"time"

	"github.com/mniak/bcpinpad"
	"github.com/mniak/bcpinpad/datalink/entangled"
	"github.com/mniak/bcpinpad/mocks"
	"github.com/mniak/bcpinpad/utils"
	"github.com/stretchr/testify/assert"
)

const TimeoutDuration = 2 * time.Second
const ToleranceDuration = 50 * time.Millisecond

func TestSendData_WhenReceiveACK_ShouldStopRetrying(t *testing.T) {
	payload := []byte("OPN000")
	expectedBytes := utils.NewBytesBuilder().
		AddByte(bcpinpad.SYN).
		AddBytes(payload).
		AddByte(bcpinpad.ETB, byte(0x77), byte(0x5e)).
		Bytes()

	alice, bob := entangled.EntangledReadWriters()
	mockReceiver := new(mocks.DataReceiver)
	mockReceiver.On("ReadACKorNAK").Return(true, nil)
	pp := NewPinpad(alice, mockReceiver)

	startTime := time.Now()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		ack, err := pp.sendData(payload)
		assert.True(t, ack, "Expected ACK")
		assert.NoError(t, err, "error while sending data")
		wg.Done()
	}()
	go func() {
		bob.Write([]byte{bcpinpad.ACK})
	}()
	time.Sleep(100 * time.Millisecond)
	wg.Wait()
	duration := time.Since(startTime)

	assert.True(t, duration < 100*time.Millisecond+ToleranceDuration, "function did not complete in the specified time duration")

	recvBuffer := make([]byte, len(expectedBytes))
	recvCount, err := bob.Read(recvBuffer)

	assert.NoError(t, err)
	assert.Equal(t, len(expectedBytes), recvCount, "wrong number of bytes sent")
	assert.EqualValues(t, expectedBytes, recvBuffer, "wrong bytes sent")
}

func TestSendData_WhenDoesNotReceiveReplyAndTimeout_ShouldAddCANandAbort(t *testing.T) {
	payload := []byte("OPN000")
	expectedBytes := utils.NewBytesBuilder().
		AddByte(bcpinpad.SYN).
		AddBytes(payload).
		AddByte(bcpinpad.ETB, byte(0x77), byte(0x5e)).
		AddByte(bcpinpad.CAN).
		Bytes()

	alice, bob := entangled.EntangledReadWriters()
	mockReceiver := new(mocks.DataReceiver)
	mockReceiver.On("ReadACKorNAK").Return(false, bcpinpad.ErrTimeout)
	pp := NewPinpad(alice, mockReceiver)

	startTime := time.Now()
	{
		ack, err := pp.sendData(payload)
		assert.False(t, ack, "OK (ACK) not expected")
		assert.Equal(t, bcpinpad.ErrTimeout, err)
	}
	duration := time.Since(startTime)

	assert.True(t, duration < TimeoutDuration+ToleranceDuration, "function did not complete in the specified time duration")

	recvBuffer := make([]byte, len(expectedBytes))
	recvCount, err := bob.Read(recvBuffer)

	assert.NoError(t, err)
	assert.Equal(t, len(expectedBytes), recvCount, "wrong number of bytes sent")
	assert.EqualValues(t, expectedBytes, recvBuffer, "wrong bytes sent")
}

func TestSend(t *testing.T) {
	requestPayload := []byte("REQUEST PAYLOAD")
	responsePayload := []byte("RESPONSE PAYLOAD")

	alice, _ := entangled.EntangledReadWriters()
	mockReceiver := new(mocks.DataReceiver)
	mockReceiver.On("ReadACKorNAK").Return(true, nil)
	mockReceiver.On("Receive").Return(responsePayload, nil)
	pp := NewPinpad(alice, mockReceiver)

	actualResponsePayload, err := pp.Send(requestPayload)
	assert.NoError(t, err)
	assert.EqualValues(t, responsePayload, actualResponsePayload, "response payload does not match with expectation")
}
