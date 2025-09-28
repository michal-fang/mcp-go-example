package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 预定义的时间格式常量
const (
	TimeFormatISO8601 = "2006-01-02T15:04:05Z07:00"
	TimeFormatHuman   = "2006-01-02 15:04:05"
	TimeFormatWithMS  = "2006-01-02 15:04:05.000"
	TimeFormatCompact = "20060102-150405"
)

type Options struct {
	Level         string // debug, info, warn, error, panic, fatal
	Format        string // json, text
	Output        string // stdout, file, both
	TimeFormat    string // iso8601, human, with-ms, etc.
	Filepath      string // 文件名
	MaxSize       int    // MB
	MaxBackups    int    // 备份文件数
	MaxAge        int    // 备份最大days
	EnableCaller  bool   // 是否启用文件名和行号
	DisableStdout bool   // 是否禁用 stdout 输出
}

func NewOptions() *Options {
	return &Options{
		Level:         viper.GetString("log.level"),
		Format:        viper.GetString("log.format"),
		Output:        viper.GetString("log.output"),
		TimeFormat:    viper.GetString("log.timeFormat"),
		Filepath:      viper.GetString("log.filepath"),
		MaxSize:       viper.GetInt("log.maxSize"),
		MaxBackups:    viper.GetInt("log.maxBackups"),
		MaxAge:        viper.GetInt("log.maxAge"),
		EnableCaller:  viper.GetBool("log.enableCaller"),
		DisableStdout: viper.GetBool("log.disableStdout"),
	}
}

func NewDefaultOptions() *Options {
	return &Options{
		Level:         "info",
		Format:        "json",
		Output:        "both",
		TimeFormat:    "human",
		Filepath:      "logs/app.log",
		MaxSize:       10,
		MaxBackups:    7,
		MaxAge:        30,
		EnableCaller:  true,
		DisableStdout: false,
	}
}

func getTimeFormat(format string) string {
	switch format {
	case "iso8601", "ISO8601":
		return TimeFormatISO8601
	case "human", "human-readable":
		return TimeFormatHuman
	case "with-ms", "milliseconds":
		return TimeFormatWithMS
	case "compact":
		return TimeFormatCompact
	case "rfc3339":
		return time.RFC3339
	case "rfc1123":
		return time.RFC1123
	default:
		return TimeFormatHuman // 默认格式
	}
}

func setLogFormatter(opts *Options, l *logrus.Logger) {
	timeFormat := getTimeFormat(opts.TimeFormat)
	// 设置格式
	switch opts.Format {
	case "json":
		l.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: getCallerPrettifier(),
			TimestampFormat:  timeFormat,
		})
	case "text":
		// 判断是否输出到终端（TTY）
		l.SetFormatter(&logrus.TextFormatter{
			ForceColors:               false,
			DisableColors:             true,
			FullTimestamp:             true,
			EnvironmentOverrideColors: true,
			CallerPrettyfier:          getCallerPrettifier(),
		})
	default:
		l.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:             true,
			TimestampFormat:           timeFormat,
			ForceColors:               false,
			DisableColors:             true,
			EnvironmentOverrideColors: true,
			CallerPrettyfier:          getCallerPrettifier(),
		})
	}
}

func setupOutput(opts *Options, l *logrus.Logger) {
	var writers []io.Writer
	// 根据配置添加输出目标
	switch opts.Output {
	case "stdout":
		if !opts.DisableStdout {
			writers = append(writers, os.Stdout)
		}
	case "file":
		if opts.Filepath != "" {
			fileOutput := setupFileOutput(opts)
			writers = append(writers, fileOutput)
		}
	case "both", "":
		writers = append(writers, os.Stdout)
		if opts.Filepath != "" {
			fileOutput := setupFileOutput(opts)
			writers = append(writers, fileOutput)
		}
	}
	if len(writers) > 1 {
		mw := io.MultiWriter(writers...)
		l.SetOutput(mw)
	}
}

func setupFileOutput(cfg *Options) io.Writer {
	// 确保日志目录存在
	dir := filepath.Dir(cfg.Filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logrus.Warnf("Failed to create log directory: %v, using stdout only", err)
		return os.Stdout
	}
	return &lumberjack.Logger{
		Filename:   cfg.Filepath,
		MaxSize:    cfg.MaxSize,    // MB
		MaxBackups: cfg.MaxBackups, // 保留的旧文件数量
		MaxAge:     cfg.MaxAge,     // 天数
		Compress:   true,           // 是否压缩
		LocalTime:  true,           // 使用本地时间
	}
}

func getCallerPrettifier() func(*runtime.Frame) (function string, file string) {
	return func(frame *runtime.Frame) (function string, file string) {
		// 只保留文件名，不包含完整路径
		_, filename := filepath.Split(frame.File)
		return "", fmt.Sprintf("%s:%d", filename, frame.Line)
	}
}
