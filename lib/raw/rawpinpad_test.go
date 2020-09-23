package raw

import (
	"sync"
	"testing"
	"time"

	"github.com/mniak/ppabecs/lib/utils"
	"github.com/stretchr/testify/assert"
)

const TimeoutDuration = 2 * time.Second
const ToleranceDuration = 50 * time.Millisecond

func TestSend_WhenReceiveACK_ShouldStopSending(t *testing.T) {
	payload := []byte("OPN000")
	expectedBytes := utils.NewBytesBuilder().
		AddByte(SYN).
		AddBytes(payload).
		AddByte(ETB, byte(0x77), byte(0x5e)).
		Bytes()

	alice, bob := utils.EntangledReadWriters()
	pp := NewPinpad(alice)

	startTime := time.Now()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		err := pp.SendData(payload)
		assert.NoError(t, err, "error while sending data")
		wg.Done()
	}()
	go func() {
		bob.Write([]byte{ACK})
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
