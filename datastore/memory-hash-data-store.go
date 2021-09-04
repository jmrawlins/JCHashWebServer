package datastore

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jmrawlins/JCHashWebServer/hash"
)

type MemoryHashDataStore struct {
	nextId     hash.HashId
	hashes     map[hash.HashId]string
	stats      Stats
	idLock     *sync.Mutex
	hashesLock *sync.Mutex
	statsLock  *sync.Mutex
}

func NewMemoryDataStore() *MemoryHashDataStore {
	ds := MemoryHashDataStore{}
	ds.idLock = &sync.Mutex{}
	ds.hashesLock = &sync.Mutex{}
	ds.statsLock = &sync.Mutex{}
	ds.hashes = make(map[hash.HashId]string)
	return &ds
}

func (ds *MemoryHashDataStore) GetNextId() (hash.HashId, error) {
	ds.idLock.Lock()
	defer ds.idLock.Unlock()
	ds.nextId += 1

	return ds.nextId, nil
}

func (ds *MemoryHashDataStore) StoreHash(id hash.HashId, password string) error {
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

func (ds *MemoryHashDataStore) GetHash(id hash.HashId) (string, error) {
	ds.hashesLock.Lock()
	defer ds.hashesLock.Unlock()

	var value string
	var ok bool
	if value, ok = ds.hashes[id]; !ok {
		return "", fmt.Errorf("invalid hash id: No hash associated with ID '%v'", id)
	}

	return value, nil
}

func (ds *MemoryHashDataStore) GetAllHashes() *map[hash.HashId]string {
	ds.idLock.Lock()
	defer ds.idLock.Unlock()
	return &ds.hashes
}

func (ds *MemoryHashDataStore) StoreRequestTime(ms int64) {
	ds.statsLock.Lock()
	defer ds.statsLock.Unlock()

	ds.stats.Average = (float64(ds.stats.Total)*ds.stats.Average + float64(ms)) / float64(ds.stats.Total+1)
	ds.stats.Total += 1
	json.NewEncoder(os.Stdout).Encode(ds.stats)
}

func (ds *MemoryHashDataStore) GetStats() Stats {
	return ds.stats
}
