package log

import (
	"fmt"
	stdlog "log"
	"os"
	"time"

	"go.uber.org/zap"
)

var (
	DefaultLogger *zap.Logger
)

func init() {
	stdlog.SetFlags(stdlog.Lshortfile)

	proccessName := os.Args[0]
	logName := fmt.Sprintf("%s-%s.log", proccessName, time.Now().Format("20060102"))

	var err error
	DefaultLogger, err = NewLogger(logName)
	if nil != err {
		stdlog.Fatalln(err)
	}

}

//flush bufferd log
//NOTE: call Flush() before proccess exit
func Flush() error {
	return DefaultLogger.Sync()
}

func NewLogger(name string) (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		name,
		"stdout",
	} //both output to file and stdout
	return cfg.Build(zap.AddCallerSkip(1))
}

func Debug(msg string, fields ...zap.Field) {
	DefaultLogger.Debug(msg, fields...)
}

func Debugf(template string, args ...interface{}) {
	DefaultLogger.Sugar().Debugf(template, args...)
}

func Info(msg string, fields ...zap.Field) {
	DefaultLogger.Info(msg, fields...)
}

func Infof(template string, args ...interface{}) {
	DefaultLogger.Sugar().Infof(template, args...)
}

func Warn(msg string, fields ...zap.Field) {
	DefaultLogger.Warn(msg, fields...)
}

func Warnf(template string, args ...interface{}) {
	DefaultLogger.Sugar().Warnf(template, args...)
}

func Error(msg string, fields ...zap.Field) {
	DefaultLogger.Error(msg, fields...)
}

func Errorf(template string, args ...interface{}) {
	DefaultLogger.Sugar().Errorf(template, args...)
}

func Fatal(msg string, fields ...zap.Field) {
	DefaultLogger.Sync()
	DefaultLogger.Fatal(msg, fields...)
}

func Fatalf(template string, args ...interface{}) {
	DefaultLogger.Sync()
	DefaultLogger.Sugar().Fatalf(template, args...)
}
