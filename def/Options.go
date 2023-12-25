package def

import (
	"github.com/sirupsen/logrus"
	"sync"
)

// LeftForward 左转 左前转子减速到70%，左后转子减速到65%，右前转子减速到85%，右后转子全速运转
// degree 左转读书，最大为40°
func LeftForward(degree uint32) {
	ps := [4]pair{
		{frontLeft, 0.7},
		{rearLeft, 0.65},
		{frontRight, 0.85},
		{rearRight, 1},
	}
	wheel(degree, ps, goForward)
}

// RightForward 右转 右前转子减速到70%，右后转子减速到65%，左前转子减速到8%，左后转子全速运转
func RightForward(degree uint32) {
	ps := [4]pair{
		{frontRight, 0.7},
		{rearRight, 0.65},
		{frontLeft, 0.85},
		{rearLeft, 1},
	}
	wheel(degree, ps, goForward)
}

// DiveForward 俯冲 前2转子减速到80%，后2转子全速运转
func DiveForward(degree uint32) {
	ps := [4]pair{
		{frontRight, 0.8},
		{frontLeft, 0.8},
		{rearRight, 1},
		{rearLeft, 1},
	}
	wheel(degree, ps, goForward)
}

// ClimbForward 攀爬 前2转子全速运转，后2转子减速到80%
func ClimbForward(degree uint32) {
	ps := [4]pair{
		{frontRight, 1},
		{frontLeft, 1},
		{rearRight, 0.8},
		{rearLeft, 0.8},
	}
	wheel(degree, ps, goForward)
}

// GoForwardWithPWM 直行
func GoForwardWithPWM(duty uint32) {
	DutyLen = duty
	var wg sync.WaitGroup
	wg.Add(len(machines))
	for _, mc := range machines {
		go func(m machine) {
			logrus.Infof("转子[%s]启动...", m.name)
			goForward(m, duty)
			wg.Done()
		}(mc)
	}
	wg.Wait()
}

// LeftInvert 左转 左前转子减速到70%，左后转子减速到65%，右前转子减速到85%，右后转子全速运转
// degree 左转读书，最大为40°
func LeftInvert(degree uint32) {
	ps := [4]pair{
		{frontLeft, 0.7},
		{rearLeft, 0.65},
		{frontRight, 0.85},
		{rearRight, 1},
	}
	wheel(degree, ps, goInvert)
}

// RightInvert 右转 右前转子减速到70%，右后转子减速到65%，左前转子减速到8%，左后转子全速运转
func RightInvert(degree uint32) {
	ps := [4]pair{
		{frontRight, 0.7},
		{rearRight, 0.65},
		{frontLeft, 0.85},
		{rearLeft, 1},
	}
	wheel(degree, ps, goInvert)
}

// DiveInvert 俯冲 前2转子减速到80%，后2转子全速运转
func DiveInvert(degree uint32) {
	ps := [4]pair{
		{frontRight, 0.8},
		{frontLeft, 0.8},
		{rearRight, 1},
		{rearLeft, 1},
	}
	wheel(degree, ps, goInvert)
}

// ClimbInvert 攀爬 后2转子全速运转，后2转子减速到80%
func ClimbInvert(degree uint32) {
	ps := [4]pair{
		{frontRight, 1},
		{frontLeft, 1},
		{rearRight, 0.8},
		{rearLeft, 0.8},
	}
	wheel(degree, ps, goInvert)
}

// GoInvertWithPWM 后退 直行
func GoInvertWithPWM(dutyLen uint32) {
	DutyLen = dutyLen
	var wg sync.WaitGroup
	wg.Add(len(machines))
	for _, mc := range machines {
		go func(m machine) {
			logrus.Infof("转子[%s]启动...", m.name)
			goInvert(m, dutyLen)
			wg.Done()
		}(mc)
	}
	wg.Wait()
}

// StopWithPWM 停止
func StopWithPWM() {
	runningCh <- true
}
