package log

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	WithField(key string, value interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}

type LogrusLogger struct {
	log *logrus.Logger
}

// Logger 日志接口
var _ Logger = &LogrusLogger{}
var (
	mu  sync.Mutex
	std = &LogrusLogger{log: logrus.New()}
)

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = NewLogger(opts)
}

func NewLogger(opts *Options) *LogrusLogger {
	if opts == nil {
		opts = NewDefaultOptions()
	}
	// 创建logrus.Logger
	l := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(opts.Level)
	if err != nil {
		level = logrus.DebugLevel
	} else {
		l.SetLevel(level)
	}
	// 设置日志格式
	setLogFormatter(opts, l)
	//设置调用者信息
	l.SetReportCaller(opts.EnableCaller)
	//设置输出目标
	setupOutput(opts, l)
	return &LogrusLogger{log: l}
}

func Debugln(args ...interface{}) {
	std.log.Debugln(args...)
}

func (l LogrusLogger) Debugln(args ...interface{}) {
	l.log.Debugln(args...)
}

func Infoln(args ...interface{}) {
	std.log.Infoln(args...)
}

func (l LogrusLogger) Infoln(args ...interface{}) {
	l.log.Infoln(args...)
}

func Warnln(args ...interface{}) {
	std.log.Warnln(args...)
}

func (l LogrusLogger) Warnln(args ...interface{}) {
	l.log.Warnln(args...)
}

func Errorln(args ...interface{}) {
	std.log.Errorln(args...)
}

func (l LogrusLogger) Errorln(args ...interface{}) {
	l.log.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	std.log.Fatalln(args...)
}

func (l LogrusLogger) Fatalln(args ...interface{}) {
	l.log.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	std.log.Panicln(args...)
}

func (l LogrusLogger) Panicln(args ...interface{}) {
	l.log.Panicln(args...)
}

func WithField(key string, value interface{}) {
	std.log.WithField(key, value)
}

func (l LogrusLogger) WithField(key string, value interface{}) {
	l.log.WithField(key, value)
}

func Info(args ...interface{}) {
	std.log.Info(args...)
}

func (l LogrusLogger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func Debug(args ...interface{}) {
	std.log.Debug(args...)
}

func (l LogrusLogger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func Warn(args ...interface{}) {
	std.log.Warn(args...)
}

func (l LogrusLogger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func Error(args ...interface{}) {
	std.log.Error(args...)
}

func (l LogrusLogger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func Fatal(args ...interface{}) {
	std.log.Fatal(args...)
}

func (l LogrusLogger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func Panic(args ...interface{}) {
	std.log.Panic(args...)
}

func (l LogrusLogger) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

func Debugf(format string, args ...interface{}) {
	std.log.Debugf(format, args...)
}

func (l LogrusLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	std.log.Infof(format, args...)
}

func (l LogrusLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	std.log.Warnf(format, args...)
}

func (l LogrusLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	std.log.Errorf(format, args...)
}

func (l LogrusLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	std.log.Fatalf(format, args...)
}

func (l LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	std.log.Panicf(format, args...)
}

func (l LogrusLogger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}
