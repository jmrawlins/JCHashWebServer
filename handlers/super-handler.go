package handlers

import (
	"net/http"
	"time"

	"github.com/jmrawlins/JCHashWebServer/datastore"
)

type SuperHandler struct {
	handler http.Handler
	sds     datastore.StatsDataStore
}

func NewSuperHandler(handler http.Handler, sds datastore.StatsDataStore) *SuperHandler {
	return &SuperHandler{handler, sds}
}
func (sh SuperHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Pre-processing (start timer)
	startTime := time.Now().UnixMilli()

	// Call the real handler
	sh.handler.ServeHTTP(resp, req)

	// Post-processing (stop timer and send timing info)
	endTime := time.Now().UnixMilli()
	sh.sds.StoreRequestTime(req.URL.Path, endTime-startTime)
}
