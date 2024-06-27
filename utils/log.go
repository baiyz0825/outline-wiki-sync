// outline-wiki-sync
//
// @(#)log.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package utils

import (
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = NewLogger()
}

func NewLogger() *logrus.Logger {
	Logger := logrus.New()
	// 日志级别
	logLev, err := logrus.ParseLevel("info")
	Logger.SetLevel(logLev)
	if err != nil {
		Logger.SetLevel(logrus.DebugLevel)
	}

	Logger.SetOutput(os.Stdout)

	Logger.SetReportCaller(true)
	// 格式化
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return frame.Function, fileName
		},
	})
	return Logger
}
