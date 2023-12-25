package def

import (
	"github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"sync"
	"time"
)

const (
	freq = 64000
)

func init() {
	if e := rpio.Open(); e != nil {
		logrus.Fatal("无法打开rpio,%v", e)
	}
	logrus.Infof("rpio已打开")
}

type machine struct {
	// name 组合名
	name string
	// pwm pwm引脚，一般为GPIO12、GPIO13、GPIO18、GPIO19，树莓派引脚说明：
	pwm uint8
	// in1 直流电机的输入引脚
	in1 uint8
	// in2 直流电机的输入引脚
	in2 uint8
}

var (
	frontLeft  = machine{"左前轮", 12, 20, 21}
	frontRight = machine{"前右轮", 19, 26, 06}
	rearLeft   = machine{"后左轮", 18, 23, 24}
	rearRight  = machine{"后右轮", 13, 22, 27}
	// machines 配置好的4个直流电机引脚
	machines = []machine{frontLeft, frontRight, rearLeft, rearRight}
)

var (
	// cycleLen PWM 周期长度
	cycleLen uint32 = 32
	// DutyLen PWM 占空长度，值越大，表示电压越高,最大值不超过cycleLen = 32
	DutyLen = cycleLen

	//最大偏转角度
	maxDegree uint32 = 40
)

var (
	runningCh = make(chan bool)
)

func running(pwm, in1, in2 uint8, dutyLen uint32, forward bool, ch chan bool) {
	if dutyLen > cycleLen {
		dutyLen = cycleLen
	}
	pwmP := rpio.Pin(pwm)
	pin1 := rpio.Pin(in1)
	pin2 := rpio.Pin(in2)
	//控制正转还是倒转
	if forward {
		rpio.WritePin(pin1, rpio.High)
		rpio.WritePin(pin2, rpio.Low)
	} else {
		rpio.WritePin(pin1, rpio.Low)
		rpio.WritePin(pin2, rpio.High)
	}
	defer func() {
		pin1.Low()
		logrus.Infof("设置[%v]为低电位", pin1)
		pin2.Low()
		logrus.Infof("设置[%v]为低电位", pin2)
		pwmP.Low()
		logrus.Infof("设置[%v]为低电位", pwmP)
	}()
	logrus.Infof("pwmP:%v,pin1:%v,pin2:%v", pwmP, pin1, pin2)
	pwmP.Mode(rpio.Output)
	rpio.WritePin(pwmP, rpio.High)
	time.Sleep(2 * time.Second)
	pwmP.Mode(rpio.Pwm)
	pwmP.Freq(freq)
	logrus.Infof("转速:%.2f %", float32(dutyLen/cycleLen)*100)
	pwmP.DutyCycle(dutyLen, cycleLen)
	logrus.Info("持续运行，等待停止信号...")
	time.Sleep(10 * time.Second)
	//<-ch
	logrus.Info("停止运行")
}

// goForward 以给定功率前进，功率值： dutyLen / cycleLen
func goForward(m machine, dutyLen uint32) {
	logrus.Infof("转子名称：%s", m.name)
	running(m.pwm, m.in1, m.in2, dutyLen, true, runningCh)
}

// goForward 以给定功率后退，功率值： dutyLen / cycleLen
func goInvert(m machine, dutyLen uint32) {
	running(m.pwm, m.in1, m.in2, dutyLen, false, runningCh)
}

var calculateDegree = func(degree uint32, rate float32) uint32 {
	return uint32(float32(DutyLen) * float32(degree/maxDegree) * rate)
}

type pair struct {
	m machine
	r float32
}

// wheel 转弯的处理方法
func wheel(degree uint32, ps [4]pair, callback func(m machine, dutyLen uint32)) {
	if degree > 40 {
		degree = 40
	}
	if degree <= 0 {
		return
	}
	var wg sync.WaitGroup
	wg.Add(4)
	var wheelRun = func(m machine, d uint32, r float32) {
		callback(m, calculateDegree(d, r))
		wg.Done()
	}
	for _, p := range ps {
		go wheelRun(p.m, degree, p.r)
	}
	wg.Wait()
}
