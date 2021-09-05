package handlers

import (
	"io"
	"net/http"
)

type ShutdownHandler struct {
	shutdown chan<- struct{}
}

func NewShutdownHandler(shutdown chan<- struct{}) *ShutdownHandler {
	return &ShutdownHandler{shutdown}
}
func (handler ShutdownHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	io.WriteString(resp, "Shutdown!\n")
	handler.shutdown <- struct{}{}
}
