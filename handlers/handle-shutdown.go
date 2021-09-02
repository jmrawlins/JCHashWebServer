package handlers

import (
	"io"
	"net/http"
)

type ShutdownHandler struct {
	ShutdownChannel chan<- bool
}

func (handler *ShutdownHandler) HandleShutdown(resp http.ResponseWriter, req *http.Request) {
	io.WriteString(resp, "Shutdown!\n")
	handler.ShutdownChannel <- true
}
