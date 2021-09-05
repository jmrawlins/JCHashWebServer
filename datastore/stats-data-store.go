package datastore

type RequestStats struct {
	URI     string  `json:"request"`
	Total   uint64  `json:"total"`
	Average float64 `json:"average"`
}

type ServerStats map[string]*RequestStats

type StatsDataStore interface {
	StoreRequestTime(uri string, ms int64)
	GetStats() (string, error)
	GetUriStats(uri string) RequestStats
}
