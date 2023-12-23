package def

import (
	"periph.io/x/conn/v3/physic"

	//"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3/bcm283x"
	"strings"
	"sync"
	"time"

	// gpio2 "github.com/brian-armstrong/gpio"
	"github.com/sirupsen/logrus"
	"periph.io/x/conn/v3/gpio"
)

const (
	PWMA = "GPIO18"
	AIN1 = "GPIO27"
	AIN2 = "GPIO22"
)

const (
	PWMB = "15"
	BIN1 = "16"
	BIN2 = "17"
)

const (
	PWMC = "26"
	CIN1 = "19"
	CIN2 = "13"
)

const (
	PWMD = "20"
	DIN1 = "25"
	DIN2 = "GPIO12"
)

var bcmPins = []*bcm283x.Pin{
	bcm283x.GPIO1,bcm283x.GPIO2, bcm283x.GPIO3,bcm283x.GPIO4,bcm283x.GPIO5,bcm283x.GPIO6,bcm283x.GPIO7,bcm283x.GPIO8,bcm283x.GPIO9,bcm283x.GPIO10,
	bcm283x.GPIO11,bcm283x.GPIO12, bcm283x.GPIO13,bcm283x.GPIO14,bcm283x.GPIO15,bcm283x.GPIO16,bcm283x.GPIO17,bcm283x.GPIO18,bcm283x.GPIO19,bcm283x.GPIO20,
	bcm283x.GPIO21,bcm283x.GPIO22, bcm283x.GPIO23,bcm283x.GPIO24,bcm283x.GPIO25,bcm283x.GPIO26,bcm283x.GPIO27,bcm283x.GPIO28,bcm283x.GPIO29,bcm283x.GPIO30,
	bcm283x.GPIO31,bcm283x.GPIO32, bcm283x.GPIO33,bcm283x.GPIO34,bcm283x.GPIO35,bcm283x.GPIO36,bcm283x.GPIO37,bcm283x.GPIO38,bcm283x.GPIO39,bcm283x.GPIO40,
}

var forward = func(pwmP, in1P, in2P gpio.PinIO) {
	logrus.Info("执行前进方法....")
	logrus.Infof("pwmP:%v",pwmP)
	_ = pwmP.In(pwmP.Pull(),gpio.NoEdge)
	_ = pwmP.Out(gpio.High)
	defer pwmP.Out(gpio.Low)
	in1P.Out(gpio.High)
	in2P.Out(gpio.Low)
	time.Sleep(2 * time.Second)
	duty, _ := gpio.ParseDuty("20%")
	print(10 * physic.MicroHertz,duty)
	//pwmP.PWM(duty,10 * physic.RPM)
	logrus.Info("倒转...")
	in1P.Out(gpio.Low)
	in2P.Out(gpio.High)
	time.Sleep(2 * time.Second)
}

func getBcmPinByName(name string) *bcm283x.Pin {
	for _, pin := range bcmPins {
		if pin.Name() == name {
			return pin
		} else if strings.HasSuffix(pin.Name(),name) {
			return pin
		}
	}
	return bcm283x.GPIO46
}

func run(pwm, in1, in2 string, callback func(pwmP, in1P, in2P gpio.PinIO)) {
	logrus.Info("开始执行run方法...")
	pwmP := getBcmPinByName(pwm)
	in1P := getBcmPinByName(in1)
	in2P := getBcmPinByName(in2)
	logrus.Infof("排出引脚，%s,%s,%s,%v,%v,%v", pwm, in1, in2, pwmP, in1P, in2P)
	defer pwmP.Out(gpio.Low)
	defer in1P.Out(gpio.Low)
	defer in2P.Out(gpio.Low)


	callback(pwmP, in1P, in2P)
}

var pinSlice = [][]string{
	{PWMA, AIN1, AIN2},
	// {PWMB, BIN1, BIN2},
	// {PWMC, CIN1, CIN2},
	// {PWMD, DIN1, DIN2},

}

func GoForward() {
	var wg sync.WaitGroup
	wg.Add(len(pinSlice))
	for edgeIndex := 0; edgeIndex < len(pinSlice); edgeIndex++ {
		pins := pinSlice[edgeIndex]
		go func() {
			run(pins[0], pins[1], pins[2], forward)
			wg.Done()
		}()
	}
	wg.Wait()

}
