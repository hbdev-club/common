package http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	//
	"github.com/google/uuid"
	//
	"github.com/hbdev-club/common/logger"
)

var (
	requestIdKey = "X-Request-Id"
)

type Middleware func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// RequestMiddleware 请求中间件，记录请求路径、参数、耗时等
func RequestMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestAt := time.Now()
		requestId := r.Header.Get(requestIdKey)
		if requestId == "" {
			requestId = uuid.New().String()
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, logger.RequestIdKey, requestId)

		baseMDC := &logger.BaseMDC{RequestId: requestId}
		requestMDC := &logger.RequestMDC{BaseMDC: baseMDC}
		parsedURL, err := url.Parse(r.RequestURI)
		if err != nil {
			logger.Log.WithCtx(ctx).Error(fmt.Sprintf("Error parse URI: %v", err))
		} else {
			requestMDC.RequestUri = parsedURL.Path
			requestMDC.RequestQuery = parsedURL.RawQuery
			// Ignore ready request
			if strings.HasSuffix(requestMDC.RequestUri, "/ready") {
				baseMDC.RequestId = logger.Ignore
				ctx = context.WithValue(ctx, logger.RequestIdKey, logger.Ignore)
			}
		}

		logger.Log.WithMDC(requestMDC).Info(fmt.Sprintf("Request method:%s, url:%s", r.Method, r.URL))
		defer func() {
			duration := time.Since(requestAt).Milliseconds()

			responseMDC := &logger.ResponseMDC{BaseMDC: baseMDC, ResponseDuration: duration}
			logger.Log.WithMDC(responseMDC).Info(fmt.Sprintf("Response duration:%vms", duration))
		}()

		w.Header().Set(requestIdKey, requestId)

		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
