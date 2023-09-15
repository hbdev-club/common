package logger

import (
	"context"
	"fmt"
	"os"
	//
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	appName = "unknown"

	fileName = "/www/app/logs/application.log"

	defaultLevel  = zapcore.InfoLevel
	encoderConfig zapcore.EncoderConfig
	Log           *Logger
)

type Logger struct {
	*zap.Logger
}

// WithCtx 从 ctx 提取 request_id，并记录在 mdc 中
func (l *Logger) WithCtx(ctx context.Context) *zap.Logger {
	return l.WithMDC(&BaseMDC{RequestId: fmt.Sprint(ctx.Value(RequestIdKey))})
}

func (l *Logger) WithMDC(mdc MDCMarshaler) *zap.Logger {
	requestId := mdc.GetRequestId()
	if requestId == Ignore {
		return l.With()
	}
	return l.With(zap.Object("mdc", mdc))
}

func buildEncoderConfig() {
	encoderConfig = zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "logger_name",
		MessageKey:     "message",
		StacktraceKey:  "exception",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// genFileOutputCore 往文件输出 json 格式日志
func genFileOutputCore() zapcore.Core {
	writer := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100, // MB
		MaxBackups: 50,
		MaxAge:     30, // Days
		Compress:   false,
	}
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	jsonEncoder.AddString("@version", "go_1")
	jsonEncoder.AddString("app_name", appName)

	return zapcore.NewCore(
		jsonEncoder,
		zapcore.AddSync(writer),
		defaultLevel,
	)
}

// genConsoleOutputCore 往控制台输出 plain-text 格式日志
func genConsoleOutputCore() zapcore.Core {
	consoleEncoder := NewCustomEncoder(encoderConfig)

	return newCustomCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		defaultLevel,
	)
}

// InitLog
//
//	Params example:
//	InitLog("service-name", "/www/app/logs/application.log", "INFO")
func InitLog(appName_, fileName_, level_ string) *Logger {
	if Log != nil {
		return Log
	}
	if appName_ != "" {
		appName = appName_
	}
	if fileName_ != "" {
		fileName = fileName_
	}
	if level_ == "" {
		defaultLevel = zapcore.InfoLevel
	} else {
		var err error
		defaultLevel, err = zapcore.ParseLevel(level_)
		if err != nil {
			panic("invalid log level")
		}
	}
	//
	buildEncoderConfig()
	Log = &Logger{zap.New(
		zapcore.NewTee(
			genFileOutputCore(), genConsoleOutputCore(),
		),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.PanicLevel),
	)}
	Log.Info(fmt.Sprintf("Init log successfully, app name: %v, file name: %v, level: %v",
		appName_, fileName_, level_))

	return Log
}
