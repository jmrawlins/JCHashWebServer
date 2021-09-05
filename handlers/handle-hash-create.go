package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/jmrawlins/JCHashWebServer/datastore"
)

type HashCreateHandler struct {
	ds datastore.HashDataStore
	wg *sync.WaitGroup
}

func NewHashCreateHandler(ds datastore.HashDataStore, wg *sync.WaitGroup) *HashCreateHandler {
	return &HashCreateHandler{ds, wg}
}

func (handler HashCreateHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := handler.ds.GetNextId()
	if err != nil {
		http.Error(resp, "error creating hash:"+err.Error(), http.StatusServiceUnavailable)
	}
	scheduleHashJob(handler.wg, handler.ds, id, req.FormValue("password"))

	fmt.Fprintf(resp, "%v", id)
}

func scheduleHashJob(wg *sync.WaitGroup, ds datastore.HashDataStore, id uint64, password string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ds.StoreHash(id, password)
	}()
}
