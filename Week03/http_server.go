package main

import (
	"net/http"
	"context"
	"fmt"
	"time"
)

type HttpServer struct {
	s       *http.Server
	handler  http.Handler
	ctx      context.Context
}

func NewServer(address string, mux http.Handler) *HttpServer {
	h := &HttpServer{ctx: context.Background()}
	h.s = &http.Server{
		Addr:            address,
		WriteTimeout:    time.Second * 3,
		Handler:         mux,
	}
	return h
}

func (h *HttpServer) Start() error {
	fmt.Println("httpServer Start")
	return h.s.ListenAndServe()
}

func (h *HttpServer) Stop() error {
	err := h.s.Shutdown(h.ctx)
	fmt.Println("httpServer Stop")
	return err
}