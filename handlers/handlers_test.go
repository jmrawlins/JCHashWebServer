package handlers

import "sync"

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

	Lock sync.Mutex
}

func (m *HashDataStoreMock) GetNextId() (id uint64, err error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	return m.GetNextIdResult.I, m.GetNextIdResult.E
}
func (m *HashDataStoreMock) StoreHash(id uint64, hash string) error {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.StoreHash_id = id
	m.StoreHash_hash = hash
	return m.StoreHashResult
}
func (m *HashDataStoreMock) GetHash(id uint64) (string, error) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.GetHash_id = id
	return m.GetHashResult.H, m.GetHashResult.E
}
func (m *HashDataStoreMock) GetAllHashes() map[uint64]string {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	return m.GetAllHashesResult
}
