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
		stats, err := handler.ds.GetUriStats("/hash")
		if err != nil {
			http.Error(resp, errorString+err.Error(), http.StatusServiceUnavailable)
		}
		encoder := json.NewEncoder(resp)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(stats); err != nil {
			http.Error(resp, errorString+err.Error(), http.StatusInternalServerError)
		}
	} else {
		stats, err := handler.ds.GetStats()
		if err != nil {
			http.Error(resp, errorString+err.Error(), http.StatusServiceUnavailable)
		}
		encoder := json.NewEncoder(resp)
		encoder.SetEscapeHTML(false)

		if err := encoder.Encode(stats); err != nil {
			http.Error(resp, errorString+err.Error(), http.StatusInternalServerError)
		}
	}
}
