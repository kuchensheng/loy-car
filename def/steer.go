package def

import (
	"github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"math"
	"time"
)

// dutyLen : 0 = 0°，5 = 45°，10 = 90°，15 = 135°，20 = 180°. dutyLen=1,degree=9°。degree=1°，dutyLen=1/9
var steerReq = 200
var steerCycleLen uint32 = 25

var middleCycleLen uint32 = 15

var currentSteerDegree float64 = 90
var zeroCycleLen uint32 = 5

func getSteerDutyLen(degree float64, left bool) uint32 {
	newCurrentDegree := currentSteerDegree
	if left {
		newCurrentDegree = newCurrentDegree - degree
	} else {
		newCurrentDegree = newCurrentDegree + degree
	}
	if newCurrentDegree < 0 {
		newCurrentDegree = 0
	}
	if newCurrentDegree > 180 {
		newCurrentDegree = 180
	}
	dutyLen := math.Ceil(newCurrentDegree / 9)

	currentSteerDegree = newCurrentDegree

	return uint32(dutyLen) + zeroCycleLen
}

type ultrasonicMache struct {
	trigger rpio.Pin
	echo    rpio.Pin
}

var steerMache *machine
var ultrasonic *ultrasonicMache

func init() {
	steerMache = &machine{
		name: "舵机",
		pwm:  rpio.Pin(7),
	}
	steerMache.pwm.Mode(rpio.Pwm)
	steerMache.pwm.Freq(steerReq)

	steerMache.run(currentSteerDegree, true)

	ultrasonic = &ultrasonicMache{
		trigger: rpio.Pin(22),
		echo:    rpio.Pin(21),
	}
	go ultrasonic.run()

}

func (u *ultrasonicMache) run() {
	u.trigger.Low()
	u.echo.Low()
	time.Sleep(2 * time.Microsecond)
	for {
		u.trigger.High()
		currentTime := time.Now()
		go func(now time.Time) {
			//等待echo返回
			for u.echo.Read() == 0 {
				//持续等待
			}
			//得到高电位信息，得到毫米
			distance := time.Now().Sub(now).Milliseconds() * 170
			logrus.Infof("距离:%d", distance)
		}(currentTime)
		time.Sleep(10 * time.Microsecond)
		u.trigger.Low()
		time.Sleep(2 * time.Microsecond)
	}
}
func (m *machine) run(degree float64, left bool) {
	m.pwm.DutyCycle(getSteerDutyLen(degree, left), steerCycleLen)
}

// Left 向左偏转
func Left(degree float64) {
	steerMache.run(degree, true)
}

func Right(degree float64) {
	steerMache.run(degree, false)
}

func Reset() {
	steerMache.pwm.DutyCycle(middleCycleLen, steerCycleLen)
}
