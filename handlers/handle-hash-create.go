package handlers

import (
	"fmt"
	"net/http"

	"github.com/jmrawlins/JCHashWebServer/hash"
	"github.com/jmrawlins/JCHashWebServer/hash/datastore"
	"github.com/jmrawlins/JCHashWebServer/services"
)

type HashCreateHandler struct {
	Ds        datastore.DataStore
	Scheduler services.HashJobScheduler
}

func (handler HashCreateHandler) HandleHashCreate(resp http.ResponseWriter, req *http.Request) {
	id, _ := handler.Ds.GetNextId() // TODO handle error
	handler.Scheduler.Schedule(hash.HashCreateRequest{Id: id, Password: req.FormValue("password")})

	fmt.Fprintf(resp, "%v", id)
}
