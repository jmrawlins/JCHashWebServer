package datastore

type Stats struct {
	Total   uint64  `json:"total"`
	Average float64 `json:"average"`
}

// TODO ADD Per-endpoint statistics gathering
type StatsDataStore interface {
	StoreRequestTime(ms int64)
	GetStats() Stats
}
