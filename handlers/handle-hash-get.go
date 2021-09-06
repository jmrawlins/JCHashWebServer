package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmrawlins/JCHashWebServer/datastore"
)

type HashGetHandler struct {
	ds datastore.HashDataStore
}

func NewHashGetHandler(ds datastore.HashDataStore) *HashGetHandler {
	return &HashGetHandler{ds}
}

func (handler HashGetHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the id from the uri
	strUri := strings.TrimLeft(req.URL.Path, "/")
	log.Println("Received request at:", strUri)

	// Get endpoint as an int, if possible
	if len(strUri) == 0 {
		fmt.Fprintf(resp, "Welcome to JCHashWebServer! Consider trying our other endpoints:\n/hash\n/#\n/stats\n/stats?all\n/shutdown")
	} else if hashId, err := strconv.ParseUint(strUri, 10, 64); err != nil {
		http.Error(resp, "Not found", http.StatusNotFound)
		log.Println("404")
	} else {
		var hashValue string
		var err error
		if hashValue, err = handler.ds.GetHash(hashId); err != nil {
			errMsg := fmt.Sprint("Hash not defined for ", hashId)
			http.Error(resp, errMsg, http.StatusNotFound)
			fmt.Fprintf(resp, errMsg, hashId)
		} else {
			fmt.Fprint(resp, hashValue)
		}
	}
}
