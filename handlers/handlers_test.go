package handlers

type HashDataStoreMock struct {
	GetNextIdResult struct {
		I uint64
		E error
	}
	StoreHashResult error
	GetHashResult   struct {
		H string
		E error
	}
	GetAllHashesResult map[uint64]string

	StoreHash_id   uint64
	StoreHash_hash string

	GetHash_id uint64
}

func (m *HashDataStoreMock) GetNextId() (id uint64, err error) {
	return m.GetNextIdResult.I, m.GetNextIdResult.E
}
func (m *HashDataStoreMock) StoreHash(id uint64, hash string) error {
	m.StoreHash_id = id
	m.StoreHash_hash = hash
	return m.StoreHashResult
}
func (m *HashDataStoreMock) GetHash(id uint64) (string, error) {
	m.GetHash_id = id
	return m.GetHashResult.H, m.GetHashResult.E
}
func (m *HashDataStoreMock) GetAllHashes() map[uint64]string { return m.GetAllHashesResult }
