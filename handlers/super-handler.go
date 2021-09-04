package handlers

import (
	"net/http"
	"time"

	"github.com/jmrawlins/JCHashWebServer/datastore"
)

type SuperHandler struct {
	ActualHandler http.Handler
	Sds           datastore.StatsDataStore
}

func (sh SuperHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Pre-processing (start timer)
	startTime := time.Now().UnixMilli()

	// Call the real handler
	sh.ActualHandler.ServeHTTP(resp, req)

	// Post-processing (stop timer and send timing info)
	endTime := time.Now().UnixMilli()
	sh.Sds.StoreRequestTime(endTime - startTime)
}
