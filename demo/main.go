package main

import (
	"bufio"
	"fmt"
	"log"
	"time"

	"github.com/mniak/ppabecs/lib/raw"
	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{
		Name:        "COM4",
		Baud:        19200,
		ReadTimeout: 2 * time.Second,
	}

	fmt.Println("Reading...")
	// serialPort, err := serial.Open("/dev/ttyACM0", mode)
	serialPort, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 4; i++ {

		// buffer := make([]byte, 4)
		// n, err := serialPort.Read(buffer)

		// fmt.Println(n, err, buffer)

		// --------------------------
		scanner := bufio.NewScanner(serialPort)
		scanner.Split(raw.PayloadScanner)

		ok := scanner.Scan()
		err = scanner.Err()

		fmt.Println(ok, err)
	}

	// ll, err := raw.NewPinpad(s)
	// fmt.Println(ll, err)
}
