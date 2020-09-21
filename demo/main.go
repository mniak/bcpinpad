package main

import (
	"fmt"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 19200}
	s, err := serial.OpenPort(c)
	fmt.Println(s, err)

	// ll, err := raw.NewPinpad(s)
	// fmt.Println(ll, err)
}
