package datastore

import (
	"fmt"
	"sync"
)

/*
MemoryHashDataStore provides an in-memory implementation of the data store interface.
As such instead of performing io it owns the resources directly
and uses its own mutexes to lock down synchronized operations.

It is expected to be instantiated by NewMemoryHashDataStore
*/
type MemoryHashDataStore struct {
	nextId       uint64
	hashes       map[uint64]string
	stats        ServerStats
	idLock       *sync.Mutex
	hashesLock   *sync.Mutex
	statsLock    *sync.Mutex
	nextIdCalled bool
}

func NewMemoryHashDataStore() *MemoryHashDataStore {
	ds := MemoryHashDataStore{}
	ds.idLock = &sync.Mutex{}
	ds.hashesLock = &sync.Mutex{}
	ds.statsLock = &sync.Mutex{}
	ds.hashes = make(map[uint64]string)
	ds.stats = make(ServerStats)
	return &ds
}

func (ds *MemoryHashDataStore) GetNextId() (uint64, error) {
	ds.idLock.Lock()
	defer ds.idLock.Unlock()

	// No more manually setting the next id
	ds.nextIdCalled = true

	ds.nextId += 1

	return ds.nextId, nil
}

func (ds *MemoryHashDataStore) SetLastId(id uint64) error {
	ds.idLock.Lock()
	defer ds.idLock.Unlock()

	if ds.nextIdCalled {
		return fmt.Errorf("setLastId called after already assigning ids. Ignoring.")
	}

	ds.nextId = id
	return nil
}

func (ds *MemoryHashDataStore) StoreHash(id uint64, hash string) error {
	ds.hashesLock.Lock()
	defer ds.hashesLock.Unlock()
	if _, ok := ds.hashes[id]; ok {
		return fmt.Errorf("WARN: Not allowed to overwrite an existing hash (%d)", id)
	}
	ds.hashes[id] = hash
	return nil
}

func (ds *MemoryHashDataStore) GetHash(id uint64) (string, error) {
	ds.hashesLock.Lock()
	defer ds.hashesLock.Unlock()

	var value string
	var ok bool
	if value, ok = ds.hashes[id]; !ok {
		return "", fmt.Errorf("invalid hash id: No hash associated with ID '%v'", id)
	}

	return value, nil
}

func (ds *MemoryHashDataStore) GetAllHashes() map[uint64]string {
	ds.idLock.Lock()
	defer ds.idLock.Unlock()
	return ds.hashes
}

func (ds *MemoryHashDataStore) StoreRequestTime(uri string, ms int64) error {
	ds.statsLock.Lock()
	defer ds.statsLock.Unlock()

	reqStat, ok := ds.stats[uri]
	if !ok {
		ds.stats[uri] = RequestStats{URI: uri, Total: 1, Average: float64(ms)}
	} else {
		reqStat.Average = (float64(reqStat.Total)*reqStat.Average + float64(ms)) / float64(reqStat.Total+1)
		reqStat.Total += 1
		ds.stats[uri] = reqStat
	}

	return nil
}

func (ds *MemoryHashDataStore) GetStats() (ServerStats, error) {
	ds.statsLock.Lock()
	defer ds.statsLock.Unlock()

	stats := make(ServerStats)
	for key, val := range ds.stats {
		stats[key] = val
	}
	return stats, nil
}

func (ds *MemoryHashDataStore) GetUriStats(uri string) (RequestStats, error) {
	ds.statsLock.Lock()
	defer ds.statsLock.Unlock()
	_, ok := ds.stats[uri]
	if !ok {
		// Return a zero entry for the address that has no requests
		return RequestStats{URI: uri, Total: 0, Average: 0}, nil
	}

	return ds.stats[uri], nil
}
