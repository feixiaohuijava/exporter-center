package logs

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var Log *logrus.Logger

type LogFormatter struct{}

func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	var file string
	var len int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		len = entry.Caller.Line
	}
	msg := fmt.Sprintf("%s devops-go %s:%d GOID:%d %s %s\n", timestamp, file, len, getGID(), strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func Logger(path string) *logrus.Logger {
	if Log != nil {
		return Log
	}
	writer, _ := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithRotationTime(20*time.Second), // 日志切割时间间隔
	)

	pathMap := lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.PanicLevel: writer,
	}

	Log = logrus.New()
	Log.SetReportCaller(true)
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		//&logrus.JSONFormatter{},
		new(LogFormatter),
	))
	Log.SetFormatter(new(LogFormatter))
	return Log
}

func LoggerHandler() gin.HandlerFunc {
	//logger := Logger("logs/bigdata-metric-exporter.log")
	return func(context *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		context.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := context.Request.Method
		// 请求路由
		reqUri := context.Request.RequestURI
		// 状态码
		statusCode := context.Writer.Status()
		// 请求IP
		clientIP := context.ClientIP()
		loggerIns.Infoln(
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri)
	}
}
