package datastore

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/jmrawlins/JCHashWebServer/hash"
)

type MemoryDataStore struct {
	nextId hash.HashId
	hashes map[hash.HashId]string
	mutex  *sync.Mutex
}

func NewMemoryDataStore() *MemoryDataStore {
	ds := MemoryDataStore{}
	ds.mutex = &sync.Mutex{}
	ds.hashes = make(map[hash.HashId]string)
	return &ds
}

func (ds *MemoryDataStore) GetNextId() (hash.HashId, error) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.nextId += 1

	return ds.nextId, nil
}

func (ds *MemoryDataStore) StoreHash(id hash.HashId, password string) error {
	time.Sleep((5 * time.Second))
	hash := sha512.Sum512([]byte(password))
	hashB64Str := base64.StdEncoding.EncodeToString(hash[:])
	ds.hashes[id] = hashB64Str
	return nil
}

func (ds *MemoryDataStore) GetHash(id hash.HashId) (string, error) {
	var value string
	var ok bool
	if value, ok = ds.hashes[id]; !ok {
		return "", fmt.Errorf("invalid hash id: No hash associated with ID '%v'", id)
	}
	return value, nil
}

func (ds *MemoryDataStore) GetAllHashes() *map[hash.HashId]string {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	return &ds.hashes
}