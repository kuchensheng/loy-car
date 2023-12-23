package main

import (
	"github.com/kuchensheng/loy-car/def"
	"periph.io/x/conn/v3/driver/driverreg"
)

func main() {
	_, _ = driverreg.Init()
	def.GoForwardWithPWM()
}
