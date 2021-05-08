package tool

import (
	"time"
)

func GetCurTimeTsp() time.Time {
	currentTime := time.Now()
	return currentTime
}

func GetBeforeMinTimeTsp(value int) time.Time {
	beforeOneMinuteTime := time.Now().Add(-time.Minute * time.Duration(value))
	return beforeOneMinuteTime
}
