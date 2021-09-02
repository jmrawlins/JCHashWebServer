package services

import (
	"github.com/jmrawlins/JCHashWebServer/datastore/hashdatastore"
	"github.com/jmrawlins/JCHashWebServer/hash"
)

type HashJobScheduler struct {
	ds hashdatastore.HashDataStore
}

func NewHashJobScheduler(ds hashdatastore.HashDataStore) HashJobScheduler {
	return HashJobScheduler{ds}
}

func (scheduler *HashJobScheduler) Schedule(request hash.HashCreateRequest) {
	go scheduler.ds.StoreHash(request.Id, request.Password) // TODO pass in error channel
}
