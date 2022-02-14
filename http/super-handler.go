package http

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jmrawlins/JCHashWebServer/datastore"
)

type SuperHandler struct {
	handler     http.Handler
	sds         datastore.StatsDataStore
	interceptor UnaryServerInterceptor
}

func NewSuperHandler(handler http.Handler, sds datastore.StatsDataStore, interceptor UnaryServerInterceptor) *SuperHandler {
	return &SuperHandler{handler, sds, interceptor}
}
func (sh SuperHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// Pre-processing (start timer)
	startTime := time.Now().UnixMilli()
	strUri := strings.TrimLeft(req.URL.Path, "/")
	log.Println("Received request at:", strUri)

	if sh.interceptor != nil {
		// Redirect the handler through the interceptor chain
		handler := func(context context.Context, resp http.ResponseWriter, req *http.Request) {
			sh.handler.ServeHTTP(resp, req)
		}
		sh.interceptor(context.TODO(), resp, req, nil, handler)
	} else {
		// Call the real handler directly
		sh.handler.ServeHTTP(resp, req)
	}

	// Post-processing (stop timer and send timing info)
	endTime := time.Now().UnixMilli()
	sh.sds.StoreRequestTime(req.URL.Path, endTime-startTime)
}
