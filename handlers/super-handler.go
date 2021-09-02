package handlers

import "net/http"

type SuperHandler struct {
	ActualHandler http.Handler
}

func (sh SuperHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Pre-processing (start timer)

	// Call the real handler
	sh.ActualHandler.ServeHTTP(resp, req)

	// Post-processing (stop timer and send timing info)
}
