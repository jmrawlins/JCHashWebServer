package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmrawlins/JCHashWebServer/datastore"
)

type StatsHandler struct {
	Ds datastore.StatsDataStore
}

func (handler StatsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewEncoder(resp).Encode(handler.Ds.GetStats()); err != nil {
		http.Error(resp, "Unable to retrieve stats: "+err.Error(), http.StatusServiceUnavailable)
	}
}
