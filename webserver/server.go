package webserver

import (
	"net/http"

	"github.com/jmrawlins/JCHashWebServer/hash/datastore"
)

type Server struct {
	ds           datastore.DataStore
	scheduler    HashJobScheduler
	errorChannel chan<- error
	// Router *HashRouter
}

// See routes.go -- func (*Server) InitRoutes() error

func (server *Server) ListenAndServe(addr string, handler http.Handler) error {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		server.errorChannel <- err
	}
	return err
}

func NewServer(ds datastore.DataStore, scheduler HashJobScheduler, shutdownChannel chan<- bool, errorChannel chan<- error) *Server {
	server := &Server{ds, scheduler, errorChannel}
	server.initRoutes(shutdownChannel)
	return server
}

func (server *Server) initRoutes(shutdownChannel chan<- bool) {
	hashGetHandler := HashGetHandler{server.ds}
	hashCreateHandler := HashCreateHandler{server.ds, server.scheduler}
	shutdownHandler := ShutdownHandler{shutdownChannel}

	http.HandleFunc("/", hashGetHandler.HandleHashGet)
	http.HandleFunc("/hash", hashCreateHandler.HandleHashCreate)
	http.HandleFunc("/shutdown", shutdownHandler.handleShutdown)
	http.HandleFunc("/stats", handleStats)
}
