package handlers

import (
	"net/http"
)

func HandleStats(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Stats!\n"))
}
