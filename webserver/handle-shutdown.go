package webserver

import (
	"io"
	"net/http"
)

type ShutdownHandler struct {
	shutdownChannel chan<- bool
}

func (handler *ShutdownHandler) handleShutdown(resp http.ResponseWriter, req *http.Request) {
	io.WriteString(resp, "Shutdown!\n")
	handler.shutdownChannel <- true
}
