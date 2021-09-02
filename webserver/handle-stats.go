package webserver

import (
	"net/http"
)

func handleStats(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Stats!\n"))
}
