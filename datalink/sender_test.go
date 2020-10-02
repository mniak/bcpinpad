package datalink

import (
	"testing"

	"github.com/mniak/bcpinpad/transport"
)

func TestSenderShouldBeTransportSender(t *testing.T) {
	r := sender{}
	var a transport.Sender = r
	_ = a
}
