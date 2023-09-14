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
	//fileName = "logs/application.log"

	defaultLevel  = zapcore.InfoLevel
	encoderConfig zapcore.EncoderConfig
	log           *Logger
)

type Logger struct {
	*zap.Logger
}

// WithCtx 从 ctx 提取 request_id，并记录在 mdc 中
func (l *Logger) WithCtx(ctx context.Context) *zap.Logger {
	return l.WithMDC(&BaseMDC{RequestId: fmt.Sprint(ctx.Value(RequestIdKey))})
}

func (l *Logger) WithMDC(mdc MDCMarshaler) *zap.Logger {
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

func InitLog(appName_, fileName_ string, defaultLevel_ zapcore.Level) {
	if appName_ != "" {
		appName = appName_
	}
	if fileName_ != "" {
		fileName = fileName_
	}
	defaultLevel = zapcore.InfoLevel
	//
	buildEncoderConfig()
	log = &Logger{zap.New(
		zapcore.NewTee(
			genFileOutputCore(), genConsoleOutputCore(),
		),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.PanicLevel),
	)}
	log.Info(fmt.Sprintf("Init log successfully, app name: %v, file name: %v, default level: %v",
		appName_, fileName_, defaultLevel_))
}
