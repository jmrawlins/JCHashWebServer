package handlers

import (
	"fmt"
	"net/http"

	"github.com/jmrawlins/JCHashWebServer/datastore/hashdatastore"
	"github.com/jmrawlins/JCHashWebServer/hash"
	"github.com/jmrawlins/JCHashWebServer/services"
)

type HashCreateHandler struct {
	Ds        hashdatastore.HashDataStore
	Scheduler services.HashJobScheduler
}

func (handler HashCreateHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	id, _ := handler.Ds.GetNextId() // TODO handle error
	handler.Scheduler.Schedule(hash.HashCreateRequest{Id: id, Password: req.FormValue("password")})

	fmt.Fprintf(resp, "%v", id)
}

func (handler HashCreateHandler) serveHTTP(resp http.ResponseWriter, req *http.Request) {

}
