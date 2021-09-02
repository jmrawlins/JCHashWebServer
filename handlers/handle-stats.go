package handlers

import (
	"net/http"
)

type StatsHandler struct {
}

func (handler StatsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Stats!\n"))
}
