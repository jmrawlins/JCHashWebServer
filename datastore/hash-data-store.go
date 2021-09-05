package datastore

type HashDataStore interface {
	GetNextId() (id uint64, err error)
	StoreHash(id uint64, hash string) error
	GetHash(id uint64) (string, error)
	GetAllHashes() map[uint64]string
}
