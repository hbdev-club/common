package logger

import "go.uber.org/zap/zapcore"

var (
	RequestIdKey = "request_id"
)

type MDCMarshaler interface {
	MarshalLogObject(enc zapcore.ObjectEncoder) error
}

type BaseMDC struct {
	RequestId string
}

type RequestMDC struct {
	RequestId    string
	RequestUri   string
	RequestQuery string
}

type ResponseMDC struct {
	RequestId        string
	ResponseDuration int64
}

func (mdc *BaseMDC) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("request_id", mdc.RequestId)

	return nil
}

func (mdc *RequestMDC) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("request_id", mdc.RequestId)
	enc.AddString("request_uri", mdc.RequestUri)
	enc.AddString("request_query", mdc.RequestQuery)

	return nil
}

func (mdc *ResponseMDC) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("request_id", mdc.RequestId)
	enc.AddInt64("response_duration", mdc.ResponseDuration)

	return nil
}
