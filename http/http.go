package http

import (
	"net/http"
)

// CustomServeMux 扩展了 ServeMux，增加了中间件
type CustomServeMux struct {
	*http.ServeMux
	middlewares []Middleware
}

func NewCustomServeMux() *CustomServeMux {
	return &CustomServeMux{
		ServeMux: http.NewServeMux(),
	}
}

func (m *CustomServeMux) Use(middlewares ...Middleware) {
	m.middlewares = append(m.middlewares, middlewares...)
}

func (m *CustomServeMux) Handle(pattern string, handler http.Handler) {
	newHandler := applyMiddlewares(handler, m.middlewares...)
	m.ServeMux.Handle(pattern, newHandler)
}

func (m *CustomServeMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	newHandler := applyMiddlewares(handler, m.middlewares...)
	m.ServeMux.Handle(pattern, newHandler)
}
