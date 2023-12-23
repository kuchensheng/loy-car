package def

import (
	"github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"time"
)

func init()  {
	if e := rpio.Open();e != nil {
		logrus.Fatal("无法打开rpio,%v",e)
	}
	logrus.Infof("rpio已打开")
}
type machine struct {
	name string
	pwm uint8
	in1 uint8
	in2 uint8
}

var machines = []machine{
	{"1",18,27,22},
}

type pwmRunController struct {
	Run bool
}

func pwmRun(pwmP *rpio.Pin,cycleLen uint32,run *pwmRunController)  {
	for run.Run {
		for i := uint32(0); i < cycleLen; i++ { // increasing brightness
			logrus.Infof("increasing brightness:%d",i)
			pwmP.DutyCycle(i, cycleLen)
			time.Sleep(time.Second/time.Duration(cycleLen))
		}
		for i := cycleLen; i > 0; i-- { // decreasing brightness
			logrus.Infof("decreasing brightness:%d",i)
			pwmP.DutyCycle(i, cycleLen)
			time.Sleep(time.Second/time.Duration(cycleLen))
		}
	}

}

var cycleLen uint32 = 32
func runDuty(pwmP *rpio.Pin,dutyLen uint32, b *pwmRunController) {
	b.Run = true
	go pwmRun(pwmP,dutyLen,b)
	time.Sleep(5 * time.Second)
	b.Run = false
}
func initPin(pwm,in1,in2 uint8) {
	pwmP := rpio.Pin(pwm)
	pin1 := rpio.Pin(in1)
	pin2 := rpio.Pin(in2)
	logrus.Infof("pwmP:%v,pin1:%v,pin2:%v",pwmP,pin1,pin2)

	rpio.WritePin(pin1,rpio.Low)
	rpio.WritePin(pin2,rpio.High)

	pwmP.Mode(rpio.Pwm)
	pwmP.Freq(64000)
	pwmP.DutyCycle(0,32)
	logrus.Infof("最高转速...")
	// the LED will be blinking at 2000Hz
	// (source frequency divided by cycle length => 64000/32 = 2000)

	// five times smoothly fade in and out
	for i := 0; i < 5; i++ {
		for i := uint32(0); i < 32; i++ { // increasing brightness
			pwmP.DutyCycle(i, 32)
			time.Sleep(time.Second/32)
		}
		for i := uint32(32); i > 0; i-- { // decreasing brightness
			pwmP.DutyCycle(i, 32)
			time.Sleep(time.Second/32)
		}
	}

	logrus.Info("停止...")
	rpio.WritePin(pin1,rpio.Low)
	pwmP.DutyCycle(0, 32)
}

func GoForwardWithPWM() {
	for _, m := range machines {
		initPin(m.pwm,m.in1,m.in2)
	}
}
