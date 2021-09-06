package datastore

type RequestStats struct {
	URI     string  `json:"request"`
	Total   uint64  `json:"total"`
	Average float64 `json:"average"`
}

type ServerStats map[string]*RequestStats

type StatsDataStore interface {
	StoreRequestTime(uri string, ms int64) error
	GetStats() (string, error)
	GetUriStats(uri string) (RequestStats, error)
}

type StatsDataStoreMock struct {
	StoreRequestTimeResult error
	GetStatsResult         struct {
		S string
		E error
	}
	GetUriStatsResults struct {
		S RequestStats
		E error
	}

	StoreRequestTime_uri string
	StoreRequestTime_ms  int64

	GetUriStats_uri string
}

func (m *StatsDataStoreMock) StoreRequestTime(uri string, ms int64) error {
	m.StoreRequestTime_uri = uri
	m.StoreRequestTime_ms = ms
	return m.StoreRequestTimeResult
}
func (m *StatsDataStoreMock) GetStats() (string, error) {
	return m.GetStatsResult.S, m.GetStatsResult.E
}
func (m *StatsDataStoreMock) GetUriStats(uri string) (RequestStats, error) {
	m.GetUriStats_uri = uri
	return m.GetUriStatsResults.S, m.GetUriStatsResults.E
}
