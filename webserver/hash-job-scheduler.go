package webserver

import (
	"github.com/jmrawlins/JCHashWebServer/hash"
	"github.com/jmrawlins/JCHashWebServer/hash/datastore"
)

type HashJobScheduler struct {
	ds datastore.DataStore
}

func NewHashJobScheduler(ds datastore.DataStore) HashJobScheduler {
	return HashJobScheduler{ds}
}

func (scheduler *HashJobScheduler) Schedule(request hash.HashCreateRequest) {
	go scheduler.ds.StoreHash(request.Id, request.Password) // TODO pass in error channel
}
