package main

import (
	"github.com/kuchensheng/loy-car/def"
	"github.com/sirupsen/logrus"
	"periph.io/x/conn/v3/driver/driverreg"
	"time"
)

var testRunning = func(optionName string, dutyLen, degree uint32, callback func(params uint32)) {
	logrus.Infof("操作:%s", optionName)
	callback(dutyLen)
	time.Sleep(time.Second * 2)
}

func main() {
	_, _ = driverreg.Init()
	//testRunning("开机", 0, 0, def.GoForwardWithPWM)

	logrus.Info("直行，全速前进")
	testRunning("全速前进", 32, 0, def.GoForwardWithPWM)
	//
	//logrus.Info("减速10%，并左转10°")
	//testRunning("减速10%", def.DutyLen-def.DutyLen/uint32(10), 0, def.GoForwardWithPWM)
	//testRunning("左转10°", def.DutyLen, 10, def.LeftForward)
	//
	//logrus.Info("再减速10%，并再左转10°")
	//testRunning("减速10%", def.DutyLen-def.DutyLen/uint32(10), 0, def.GoForwardWithPWM)
	//testRunning("左转10°", def.DutyLen, 10, def.LeftForward)
	//
	//logrus.Info("再减速10%，并再左转10°")
	//testRunning("减速10%", def.DutyLen-def.DutyLen/uint32(10), 0, def.GoForwardWithPWM)
	//testRunning("左转10°", def.DutyLen, 10, def.LeftForward)
	//
	//logrus.Info("右转10°，并加速10%")
	//testRunning("右转10°", def.DutyLen, 10, def.RightForward)
	//testRunning("加速10%", def.DutyLen+def.DutyLen/uint32(10), 0, def.GoForwardWithPWM)
	//
	//logrus.Info("右转10°，并加速10%")
	//testRunning("右转10°", def.DutyLen, 10, def.RightForward)
	//testRunning("加速10%", def.DutyLen+def.DutyLen/uint32(10), 0, def.GoForwardWithPWM)
	//
	//logrus.Info("右转10°，并加速10%")
	//testRunning("右转10°", def.DutyLen, 10, def.RightForward)
	//testRunning("加速10%", def.DutyLen+def.DutyLen/uint32(10), 0, def.GoForwardWithPWM)
	//
	//logrus.Info("停止")
	//testRunning("停止", 0, 0, def.GoForwardWithPWM)
	//
	//logrus.Info("后退")
	//testRunning("后退", 16, 0, def.GoInvertWithPWM)
	//
	//logrus.Info("后退并左转舵")
	//testRunning("左转舵", 16, 10, def.LeftInvert)
	//
	//logrus.Info("后退并右转舵")
	//testRunning("右转舵", 16, 10, def.RightInvert)
	//
	//logrus.Infof("俯冲")
	//testRunning("俯冲", 32, 0, def.DiveForward)
	//
	//logrus.Infof("攀爬")
	//testRunning("攀爬", 32, 0, def.ClimbForward)
}
