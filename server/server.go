package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/jmrawlins/JCHashWebServer/datastore"
	"github.com/jmrawlins/JCHashWebServer/router"

	"github.com/jmrawlins/JCHashWebServer/handlers"
)

type Server struct {
	hs           http.Server
	hds          datastore.HashDataStore
	sds          datastore.StatsDataStore
	errorChannel chan<- error
	wg           *sync.WaitGroup
}

func (srv *Server) ListenAndServe(addr string, handler http.Handler) error {
	srv.hs.Addr = addr
	err := srv.hs.ListenAndServe()
	return err
}

func (srv *Server) Shutdown() {
	srv.hs.Shutdown(context.Background())
}

func NewServer(
	wg *sync.WaitGroup,
	ds datastore.HashDataStore,
	sds datastore.StatsDataStore,
	shutdownCalled chan<- struct{},
	errorChannel chan<- error) *Server {

	srv := &Server{wg: wg, hds: ds, sds: sds, errorChannel: errorChannel}
	srv.initRoutes(shutdownCalled)
	return srv
}

func (srv *Server) initRoutes(shutdownChannel chan<- struct{}) {
	routes := make(map[string]http.Handler)

	hashGetHandler := handlers.HashGetHandler{Ds: srv.hds}
	hashCreateHandler := handlers.HashCreateHandler{Wg: srv.wg, Ds: srv.hds}
	shutdownHandler := handlers.ShutdownHandler{ShutdownChannel: shutdownChannel}
	statsHandler := handlers.StatsHandler{Ds: srv.sds}

	routes["/"] = hashGetHandler
	routes["/hash"] = hashCreateHandler
	routes["/shutdown"] = shutdownHandler
	routes["/stats"] = statsHandler

	router.InitRoutes(srv.sds, routes)
}
