package server

import (
	"net/http"

	"github.com/jmrawlins/JCHashWebServer/handlers"
	"github.com/jmrawlins/JCHashWebServer/hash/datastore"
	"github.com/jmrawlins/JCHashWebServer/services"
)

type Server struct {
	ds           datastore.DataStore
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

func NewServer(ds datastore.DataStore, scheduler services.HashJobScheduler, shutdownChannel chan<- bool, errorChannel chan<- error) *Server {
	srv := &Server{ds, scheduler, errorChannel}
	srv.initRoutes(shutdownChannel)
	return srv
}

func (srv *Server) initRoutes(shutdownChannel chan<- bool) {
	hashGetHandler := handlers.HashGetHandler{Ds: srv.ds}
	hashCreateHandler := handlers.HashCreateHandler{Ds: srv.ds, Scheduler: srv.scheduler}
	shutdownHandler := handlers.ShutdownHandler{ShutdownChannel: shutdownChannel}

	http.HandleFunc("/", hashGetHandler.HandleHashGet)
	http.HandleFunc("/hash", hashCreateHandler.HandleHashCreate)
	http.HandleFunc("/shutdown", shutdownHandler.HandleShutdown)
	http.HandleFunc("/stats", handlers.HandleStats)
}
