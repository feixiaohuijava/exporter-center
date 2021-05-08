package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var loggerIns *logrus.Logger

func init() {
	loggerIns = Logger("logs/exporter-center.log")
	fmt.Println(loggerIns)
}

func Infoln(args ...interface{}) {
	loggerIns.Infoln(args...)
}

func Errorln(args ...interface{}) {
	loggerIns.Errorln(args...)
}
