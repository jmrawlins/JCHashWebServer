package datastore

import "github.com/jmrawlins/JCHashWebServer/hash"

type DataStore interface {
	GetNextId() (id hash.HashId, err error)
	StoreHash(id hash.HashId, hash string) error
	GetHash(id hash.HashId) (string, error)
	GetAllHashes() *map[hash.HashId]string
}
