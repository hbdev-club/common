package logger

import "go.uber.org/zap/zapcore"

var (
	RequestIdKey = "request_id"

	Ignore = "ignore"
)

type MDCMarshaler interface {
	MarshalLogObject(enc zapcore.ObjectEncoder) error
	GetRequestId() string
}

type BaseMDC struct {
	RequestId string
}

type RequestMDC struct {
	*BaseMDC
	//RequestId    string
	RequestUri   string
	RequestQuery string
}

type ResponseMDC struct {
	*BaseMDC
	//RequestId        string
	ResponseDuration int64
}

func (mdc *BaseMDC) GetRequestId() string {
	return mdc.RequestId
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
