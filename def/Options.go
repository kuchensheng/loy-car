package def

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var wheelFunc = func(degree, rate float32) float32 {
	return (float32(maxDegree) - degree) / float32(maxDegree) * rate
}

// LeftForward 左转 左前转子减速到70%，左后转子减速到65%，右前转子减速到85%，右后转子全速运转
// degree 左转读书，最大为40°
func LeftForward(degree uint32) {
	ps := [4]pair{
		{frontLeft, wheelFunc(float32(degree), 0.75)},
		{rearLeft, wheelFunc(float32(degree), 0.65)},
		{frontRight, wheelFunc(float32(degree), 0.9)},
		{rearRight, 1},
	}
	wheel(ps)
}

// RightForward 右转 右前转子减速到70%，右后转子减速到65%，左前转子减速到8%，左后转子全速运转
func RightForward(degree uint32) {
	ps := [4]pair{
		{frontRight, wheelFunc(float32(degree), 0.75)},
		{rearRight, wheelFunc(float32(degree), 0.65)},
		{frontLeft, wheelFunc(float32(degree), 0.9)},
		{rearLeft, 1},
	}
	wheel(ps)
}

// DiveForward 俯冲 前2转子减速到80%，后2转子全速运转
func DiveForward(degree uint32) {
	ps := [4]pair{
		{frontRight, wheelFunc(float32(degree), 0.8)},
		{frontLeft, wheelFunc(float32(degree), 0.8)},
		{rearRight, 1},
		{rearLeft, 1},
	}
	wheel(ps)
}

// ClimbForward 攀爬 前2转子全速运转，后2转子减速到80%
func ClimbForward(degree uint32) {
	ps := [4]pair{
		{frontRight, 1},
		{frontLeft, 1},
		{rearRight, wheelFunc(float32(degree), 0.8)},
		{rearLeft, wheelFunc(float32(degree), 0.8)},
	}
	wheel(ps)
}

// GoForwardWithPWM 直行
func GoForwardWithPWM(duty uint32) {
	DutyLen = duty
	var wg sync.WaitGroup
	wg.Add(len(machines))
	for _, mc := range machines {
		go func(m *machine, dutyLen uint32) {
			logrus.Infof("转子[%s]启动...", m.name)
			goForward(m, dutyLen)
			wg.Done()
		}(mc, duty)
	}
	wg.Wait()
}

// LeftInvert 左转 左前转子减速到70%，左后转子减速到65%，右前转子减速到85%，右后转子全速运转
// degree 左转读书，最大为40°
func LeftInvert(degree uint32) {
	ps := [4]pair{
		{frontLeft, wheelFunc(float32(degree), 0.7)},
		{rearLeft, wheelFunc(float32(degree), 0.65)},
		{frontRight, wheelFunc(float32(degree), 0.85)},
		{rearRight, 1},
	}
	wheel(ps)
}

// RightInvert 右转 右前转子减速到70%，右后转子减速到65%，左前转子减速到80%，左后转子全速运转
func RightInvert(degree uint32) {
	ps := [4]pair{
		{frontRight, wheelFunc(float32(degree), 0.7)},
		{rearRight, wheelFunc(float32(degree), 0.65)},
		{frontLeft, wheelFunc(float32(degree), 0.8)},
		{rearLeft, 1},
	}
	wheel(ps)
}

// DiveInvert 俯冲 前2转子减速到80%，后2转子全速运转
func DiveInvert(degree uint32) {
	ps := [4]pair{
		{frontRight, wheelFunc(float32(degree), 0.8)},
		{frontLeft, wheelFunc(float32(degree), 0.8)},
		{rearRight, 1},
		{rearLeft, 1},
	}
	wheel(ps)
}

// ClimbInvert 攀爬 后2转子全速运转，后2转子减速到80%
func ClimbInvert(degree uint32) {
	ps := [4]pair{
		{frontRight, 1},
		{frontLeft, 1},
		{rearRight, wheelFunc(float32(degree), 0.8)},
		{rearLeft, wheelFunc(float32(degree), 0.8)},
	}
	wheel(ps)
}

// GoInvertWithPWM 后退 直行
func GoInvertWithPWM(dutyLen uint32) {
	DutyLen = dutyLen
	var wg sync.WaitGroup
	wg.Add(len(machines))
	for _, mc := range machines {
		go func(m *machine) {
			logrus.Infof("转子[%s]倒转...", m.name)
			goInvert(m, dutyLen)
			wg.Done()
		}(mc)
	}
	wg.Wait()
}

// StopWithPWM 停止前进
func StopWithPWM() {
	for _, m := range machines {
		m.state = false
		m.signal <- true
	}
}
