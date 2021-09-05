package datastore

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

/*
MemoryHashDataStore provides an in-memory implementation of the data store interface.
As such instead of performing io it owns the resources directly
and uses its own mutexes to lock down synchronized operations.
*/
type MemoryHashDataStore struct {
	nextId     uint64
	hashes     map[uint64]string
	stats      ServerStats
	idLock     *sync.Mutex
	hashesLock *sync.Mutex
	statsLock  *sync.Mutex
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
	ds.nextId += 1

	return ds.nextId, nil
}

func (ds *MemoryHashDataStore) StoreHash(id uint64, password string) error {
	time.Sleep((5 * time.Second))

	hash := sha512.Sum512([]byte(password))
	hashB64Str := base64.StdEncoding.EncodeToString(hash[:])
	ds.hashesLock.Lock()
	defer ds.hashesLock.Unlock()
	if _, ok := ds.hashes[id]; ok {
		log.Fatalln("Setting an already set id! That should never happen!")
	}
	ds.hashes[id] = hashB64Str
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

func (ds *MemoryHashDataStore) GetAllHashes() *map[uint64]string {
	ds.idLock.Lock()
	defer ds.idLock.Unlock()
	return &ds.hashes
}

func (ds *MemoryHashDataStore) StoreRequestTime(uri string, ms int64) {
	ds.statsLock.Lock()
	defer ds.statsLock.Unlock()

	reqStat, ok := ds.stats[uri]
	if !ok {
		ds.stats[uri] = &RequestStats{URI: uri, Total: 1, Average: float64(ms)}
	} else {
		reqStat.Average = (float64(reqStat.Total)*reqStat.Average + float64(ms)) / float64(reqStat.Total+1)
		reqStat.Total += 1
		ds.stats[uri] = reqStat
	}
}

func (ds *MemoryHashDataStore) GetStats() (string, error) {
	ds.statsLock.Lock()
	defer ds.statsLock.Unlock()
	stats, err := json.Marshal(ds.stats)
	return string(stats), err
}

func (ds *MemoryHashDataStore) GetUriStats(uri string) RequestStats {
	ds.statsLock.Lock()
	defer ds.statsLock.Unlock()
	_, ok := ds.stats[uri]
	if !ok {
		ds.stats[uri] = &RequestStats{URI: uri, Total: 0, Average: 0}
	}

	return *ds.stats[uri]
}
