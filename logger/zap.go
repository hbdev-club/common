package logger

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type customEncoder struct {
	zapcore.Encoder
}

type customCore struct {
	zapcore.Core
}

func (c customEncoder) Clone() zapcore.Encoder {
	return c.Encoder.Clone()
}

func (c customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	return c.Encoder.EncodeEntry(entry, nil)
}

func NewCustomEncoder(cfg zapcore.EncoderConfig) (enc zapcore.Encoder) {
	return customEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg),
	}
}

func (c *customCore) With(fields []zapcore.Field) zapcore.Core {
	return c.Core.With(nil)
}

func newCustomCore(enc zapcore.Encoder, ws zapcore.WriteSyncer, enabler zapcore.LevelEnabler) zapcore.Core {
	return &customCore{
		zapcore.NewCore(enc, ws, enabler),
	}
}
