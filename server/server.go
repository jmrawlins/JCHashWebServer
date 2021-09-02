package server

import (
	"net/http"

	"github.com/jmrawlins/JCHashWebServer/datastore/hashdatastore"

	"github.com/jmrawlins/JCHashWebServer/handlers"
	"github.com/jmrawlins/JCHashWebServer/services"
)

type Server struct {
	ds           hashdatastore.HashDataStore
	scheduler    services.HashJobScheduler
	errorChannel chan<- error
}

func (srv *Server) ListenAndServe(addr string, handler http.Handler) error {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		srv.errorChannel <- err
	}
	return err
}

func NewServer(ds hashdatastore.HashDataStore, scheduler services.HashJobScheduler, shutdownChannel chan<- bool, errorChannel chan<- error) *Server {
	srv := &Server{ds, scheduler, errorChannel}
	srv.initRoutes(shutdownChannel)
	return srv
}

func (srv *Server) initRoutes(shutdownChannel chan<- bool) {
	hashGetHandler := handlers.HashGetHandler{Ds: srv.ds}
	hashCreateHandler := handlers.HashCreateHandler{Ds: srv.ds, Scheduler: srv.scheduler}
	shutdownHandler := handlers.ShutdownHandler{ShutdownChannel: shutdownChannel}
	statsHandler := handlers.StatsHandler{}

	http.Handle("/", hashGetHandler)
	http.Handle("/hash", hashCreateHandler)
	http.Handle("/shutdown", shutdownHandler)
	http.Handle("/stats", statsHandler)
}
