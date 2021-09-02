package webserver

import (
	"fmt"
	"net/http"

	"github.com/jmrawlins/JCHashWebServer/hash"
	"github.com/jmrawlins/JCHashWebServer/hash/datastore"
)

type HashCreateHandler struct {
	ds        datastore.DataStore
	scheduler HashJobScheduler
}

func (handler HashCreateHandler) HandleHashCreate(resp http.ResponseWriter, req *http.Request) {
	id, _ := handler.ds.GetNextId() // TODO handle error
	handler.scheduler.Schedule(hash.HashCreateRequest{Id: id, Password: req.FormValue("password")})

	fmt.Fprintf(resp, "%v", id)
}
