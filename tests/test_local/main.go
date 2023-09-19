package main

import (
	"github.com/hbdev-club/common/logger"
	"net/http"
	"os"
	//
	customhttp "github.com/hbdev-club/common/http"
)

var (
	log = logger.InitLog("service-test-local", "run.log", "")
)

func ready(w http.ResponseWriter, r *http.Request) {
	log.WithCtx(r.Context()).Info("Ready !!!")
	_, err := w.Write([]byte("ready"))
	if err != nil {
		return
	}
}

func ready2(w http.ResponseWriter, r *http.Request) {
	log.WithCtx(r.Context()).Info("Ready !!!")
	_, err := w.Write([]byte("ready"))
	if err != nil {
		return
	}
}

func main() {
	mux := customhttp.NewCustomServeMux()
	middlewares := []customhttp.Middleware{
		customhttp.RequestMiddleware,
	}
	mux.Use(middlewares...)
	mux.HandleFunc("/ready", ready)
	mux.HandleFunc("/ready2", ready2)

	server := &http.Server{
		Addr:    ":10002",
		Handler: mux,
	}
	log.Info("http://127.0.0.1:10002")
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err.Error())
		return
	}

	os.Exit(0)
}
