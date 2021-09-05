package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmrawlins/JCHashWebServer/datastore"
)

type StatsHandler struct {
	ds datastore.StatsDataStore
}

func NewStatsHandler(ds datastore.StatsDataStore) *StatsHandler {
	return &StatsHandler{ds}
}

func (handler StatsHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	errorString := "Unable to retrieve stats: "
	if req.URL.RawQuery == "" {
		if err := json.NewEncoder(resp).Encode(handler.ds.GetUriStats("/hash")); err != nil {
			http.Error(resp, errorString+err.Error(), http.StatusServiceUnavailable)
		}
	} else {
		stats, err := handler.ds.GetStats()
		if err != nil {
			http.Error(resp, errorString+err.Error(), http.StatusServiceUnavailable)
		}
		if err := json.NewEncoder(resp).Encode(stats); err != nil {
			http.Error(resp, errorString+err.Error(), http.StatusInternalServerError)
		}
	}
}
