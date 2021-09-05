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
		fmt.Fprintf(resp, "Up and running...")
	} else if hashId, err := strconv.ParseUint(strUri, 10, 64); err != nil {
		log.Println("404")
	} else {
		var hashValue string
		var err error
		if hashValue, err = handler.ds.GetHash(hashId); err != nil {
			fmt.Fprint(resp, "404 hash not defined for ", hashId)
		} else {
			log.Printf("{ \"id\":%d, \"hash\":\"%s\"\n", hashId, hashValue)
			fmt.Fprint(resp, hashValue)
		}
	}
}
