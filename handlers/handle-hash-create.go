package handlers

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

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

	if strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data") {
		if err := req.ParseMultipartForm(2048); err != nil {
			errMsg := fmt.Sprint("Error parsing multipart form data: ", err.Error())
			log.Println(errMsg)
			http.Error(resp, errMsg, http.StatusBadRequest)
			return
		}
	} else if strings.HasPrefix(req.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		if err := req.ParseForm(); err != nil {
			errMsg := fmt.Sprint("Error urlencoded form data:", err)
			log.Println(errMsg)
			http.Error(resp, errMsg, http.StatusBadRequest)
			return
		}
	}

	fmt.Fprint(resp, id)
	scheduleHashJob(handler.wg, handler.ds, id, req.FormValue("password"))
}

func scheduleHashJob(wg *sync.WaitGroup, ds datastore.HashDataStore, id uint64, password string) {
	wg.Add(1)
	time.AfterFunc(5*time.Second, func() {
		defer wg.Done()
		hash := sha512.Sum512([]byte(password))
		hashB64Str := base64.StdEncoding.EncodeToString(hash[:])
		ds.StoreHash(id, hashB64Str)
	})

}
