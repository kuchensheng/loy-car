package main

import (
	"github.com/stianeikeland/go-rpio/v4"
	"os"
	"time"
)

func main() {
	err := rpio.Open()
	if err != nil {
		os.Exit(1)
	}
	defer rpio.Close()

	pin := rpio.Pin(18)
	rpio.WritePin(rpio.Pin(22),rpio.High)
	rpio.WritePin(rpio.Pin(27),rpio.Low)
	pin.Mode(rpio.Pwm)
	pin.Freq(64000)
	pin.DutyCycle(0, 32)
	// the LED will be blinking at 2000Hz
	// (source frequency divided by cycle length => 64000/32 = 2000)

	// five times smoothly fade in and out
	for {
		for i := uint32(0); i < 32; i++ { // increasing brightness
			pin.DutyCycle(i, 32)
			time.Sleep(time.Second/32)
		}
		for i := uint32(32); i > 0; i-- { // decreasing brightness
			pin.DutyCycle(i, 32)
			time.Sleep(time.Second/32)
		}
	}
}
