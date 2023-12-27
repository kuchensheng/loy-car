package def

import (
	"github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"sync"
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
	pwm rpio.Pin
	// in1 直流电机的输入引脚
	in1 rpio.Pin
	// in2 直流电机的输入引脚
	in2 rpio.Pin
	// state 发动机状态，true表示在运行，false未运行
	state bool
	// forward 发动机转动方向,true表示正转，false倒转
	forward bool
	// signal 信号通道，当收到信号时，转子停止转动
	signal chan bool
}

var (
	frontLeft  = &machine{"左前轮", rpio.Pin(12), rpio.Pin(20), rpio.Pin(16), false, true, make(chan bool)}
	frontRight = &machine{"前右轮", rpio.Pin(19), rpio.Pin(06), rpio.Pin(26), false, true, make(chan bool)}
	rearLeft   = &machine{"后左轮", rpio.Pin(18), rpio.Pin(23), rpio.Pin(24), false, true, make(chan bool)}
	rearRight  = &machine{"后右轮", rpio.Pin(13), rpio.Pin(22), rpio.Pin(27), false, true, make(chan bool)}
	// machines 配置好的4个直流电机引脚
	machines = []*machine{frontLeft, frontRight, rearLeft, rearRight}
)

var (
	// cycleLen PWM 周期长度
	cycleLen uint32 = 32
	// DutyLen PWM 占空长度，值越大，表示电压越高,最大值不超过cycleLen = 32
	DutyLen = cycleLen

	//最大偏转角度
	maxDegree uint32 = 40

	currentDegree uint32 = 0
)

func (m *machine) running(dutyLen uint32) {
	if dutyLen > cycleLen {
		dutyLen = cycleLen
	}

	if m.state {
		logrus.Infof("转子[%s]的状态为[%v],发送停止信息", m.name, m.state)
		m.signal <- true
	}
	logrus.Infof("将转子[%s]的状态设置为运行中", m.name)
	m.state = true
	//控制正转还是倒转
	if m.forward {
		rpio.WritePin(m.in1, rpio.High)
		rpio.WritePin(m.in2, rpio.Low)
	} else {
		rpio.WritePin(m.in1, rpio.Low)
		rpio.WritePin(m.in2, rpio.High)
	}
	defer func(p ...*rpio.Pin) {
		m.state = false
		for _, pin := range p {
			logrus.Infof("设置[%v]为低电位", *pin)
			pin.Mode(rpio.Output)
			pin.Low()
		}
	}(&m.in1, &m.in2, &m.pwm)
	logrus.Infof("machine:%v", m)
	m.pwm.Mode(rpio.Pwm)
	m.pwm.Freq(freq)

	logrus.Infof("转速:%.2f %", float32(dutyLen/cycleLen)*100)
	m.pwm.DutyCycle(dutyLen, cycleLen)
	logrus.Info("持续运行，等待停止信号...")
loop:
	for {
		select {
		case <-m.signal:
			logrus.Infof("转子[%s]收到停止信号,携程即将停止运行", m.name)
			break loop
		default:
			//持续运行
		}
	}

	logrus.Infof("转子[%s]停止运行", m.name)
}

// goForward 以给定功率前进，功率值： dutyLen / cycleLen
func goForward(m *machine, dutyLen uint32) {
	logrus.Infof("前进：转子名称：%s", m.name)
	m.forward = true
	m.running(dutyLen)
}

// goForward 以给定功率后退，功率值： dutyLen / cycleLen
func goInvert(m *machine, dutyLen uint32) {
	logrus.Infof("后退：转子名称：%s", m.name)
	m.forward = false
	m.running(dutyLen)
}

type pair struct {
	m *machine
	//r 根据偏转角度计算各个轮子的运动速度
	r float32
}

func newPair(m *machine) pair {
	return pair{m, 1}
}

// wheel 转弯的处理方法
func wheel(ps [4]pair) {
	var wg sync.WaitGroup
	wg.Add(len(ps))
	var wheelRun = func(m *machine, d uint32) {
		m.running(d)
		wg.Done()
	}
	for _, p := range ps {
		go wheelRun(p.m, uint32(p.r))
	}
	wg.Wait()
}
