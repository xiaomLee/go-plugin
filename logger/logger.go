package logger

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
	"time"
)

const RequestId = "X-Request-Id"

var (
	DefaultLoggerFormat = textFormatter

	textFormatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	}

	jsonFormatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	}
)

var (
	_appName      = ""
	_formatter    logrus.Formatter
	_reportCaller = false
)

type Config struct {
	AppName string `json:"app_name" mapstructure:"app_name"`
	// panic=0, fatal=1, error=2, warning=3, info=4, debug=5, trace=6
	Level int `json:"level" mapstructure:"level"`
	Formatter string `json:"formatter" mapstructure:"formatter"`
	Dir string `json:"dir" mapstructure:"dir"`
	Filename string `json:"filename" mapstructure:"filename"`
	// max keep days. MaxAge and RotationCount cannot be both set
	MaxAge int `json:"max_age" mapstructure:"max_age"`
	// rotate size. M
	RotationSize int64 `json:"rotation_size" mapstructure:"rotation_size"`
	// rotate time. second, default 24*time.Hour
	RotationTime int `json:"rotation_time" mapstructure:"rotation_time"`
	// rotate count
	RotationCount uint `json:"rotation_count" mapstructure:"rotation_count"`
}

// ConfigLogger 配置业务日志系统
func ConfigLogger(cfg Config) {
	logrus.SetLevel(logrus.Level(cfg.Level))
	_appName = cfg.AppName
	if cfg.Formatter == "json" {
		_formatter = jsonFormatter
	} else {
		_formatter = textFormatter
	}
	logrus.SetFormatter(_formatter)

	// 设置打印文件名 行号, 递归runtime.Caller()性能消耗大
	//logrus.SetReportCaller(true)
	// 自己封装获取 file:function:line
	_reportCaller = true

	// if both set, give maxAge a sane default
	if cfg.MaxAge > 0 && cfg.RotationCount > 0 {
		cfg.MaxAge = -1
	}
	if cfg.RotationTime ==0{
		cfg.RotationTime = 86400
	}
	if cfg.RotationSize == 0 {
		// default 1GB
		cfg.RotationSize=1024
	}
	configLocalFs(cfg.Dir, cfg.Filename,
		rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge)*24*time.Hour),
		rotatelogs.WithRotationSize(cfg.RotationSize*1024*1024),
		rotatelogs.WithRotationTime(time.Duration(cfg.RotationTime)*time.Second),
		rotatelogs.WithRotationCount(cfg.RotationCount),
		)
}

func SetReportCaller(include bool) {
	_reportCaller = include
}

//配置本地文件系统并按周期分割
func configLocalFs(logPath string, logFileName string, options ...rotatelogs.Option) {
	linkPath := path.Join(logPath, logFileName)
	options = append(options,rotatelogs.WithLinkName(linkPath))

	filename := strings.TrimSuffix(logFileName, path.Ext(logFileName))
	writer, _ := rotatelogs.New(path.Join(logPath, filename+"-%Y%m%d%H%M%S"+".log"),options...)

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, _formatter)
	logrus.AddHook(lfHook)
}

// ---------------------------------------------------------------------------------------------------------------------
// pkg func

func Error(msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["app"] = _appName
	if _reportCaller {
		fields["function"] = getCaller(1)
	}
	logrus.WithFields(fields).Error(msg)
}

func Info(msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["app"] = _appName
	if _reportCaller {
		fields["function"] = getCaller(1)
	}
	logrus.WithFields(fields).Info(msg)
}

func Debug(msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["app"] = _appName
	if _reportCaller {
		fields["function"] = getCaller(1)
	}
	logrus.WithFields(fields).Debug(msg)
}

func Warning(msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["app"] = _appName
	if _reportCaller {
		fields["function"] = getCaller(1)
	}
	logrus.WithFields(fields).Warn(msg)
}

func Fatal(msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["app"] = _appName
	if _reportCaller {
		fields["function"] = getCaller(1)
	}
	logrus.WithFields(fields).Fatal(msg)
}

func ErrorWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	genCommon(ctx, fields)
	logrus.WithFields(fields).Error(msg)
}

func InfoWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	genCommon(ctx, fields)
	logrus.WithFields(fields).Info(msg)
}

func DebugWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	genCommon(ctx, fields)
	logrus.WithFields(fields).Debug(msg)
}

func WarningWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	genCommon(ctx, fields)
	logrus.WithFields(fields).Warn(msg)
}

func FatalWithContext(ctx context.Context, msg string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	genCommon(ctx, fields)
	logrus.WithFields(fields).Fatal(msg)
}

func getCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip + 1)
	//fmt.Println("getCaller", pc, file, line, ok)
	if !ok {
		return ""
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}

	function := ""
	if pc != 0 {
		frames := runtime.CallersFrames([]uintptr{pc})
		frame, _ := frames.Next()
		function = frame.Function
	}

	return fmt.Sprintf("%s:%s:%d", file, path.Base(function), line)
}

func genCommon(ctx context.Context, fields map[string]interface{}) {
	fields["app"] = _appName
	if _reportCaller {
		fields["function"] = getCaller(2)
	}
	fields["context"] = ctx
	fields[RequestId] = ctx.Value(RequestId)
}
