package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/jmrawlins/JCHashWebServer/datastore"
	"github.com/jmrawlins/JCHashWebServer/handlers"
)

type Server struct {
	hs             http.Server
	hds            datastore.HashDataStore
	sds            datastore.StatsDataStore
	err            chan error
	shutdownCalled chan struct{}
	wg             *sync.WaitGroup
	port           uint16
}

func NewServer(
	wg *sync.WaitGroup,
	hds datastore.HashDataStore,
	sds datastore.StatsDataStore,
	shutdownCalled chan struct{},
	errorChannel chan error,
	port uint16,
) *Server {

	srv := &Server{wg: wg, hds: hds, sds: sds, err: errorChannel, shutdownCalled: shutdownCalled, port: port}
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

	for routeSpec, handler := range routes {
		superHandler := handlers.SuperHandler{ActualHandler: handler, Sds: srv.sds}
		http.Handle(routeSpec, superHandler)
	}
}

func InitRoutes(sds datastore.StatsDataStore, routes map[string]http.Handler) {
}

// Listen and serve at the requested address, overriding the default serve mux
func (srv *Server) ListenAndServe(addr string, handler http.Handler) error {
	srv.hs.Addr = addr
	srv.hs.Handler = handler
	err := srv.hs.ListenAndServe()
	return err
}

// Signals the http server to shut down, allowing existing request handlers
// to complete their work
func (srv *Server) Shutdown() {
	srv.hs.Shutdown(context.Background())
}

// Runs the server with graceful shutdown
func (srv *Server) Run() error {
	defer srv.wg.Wait()

	srv.wg.Add(1)
	// Run the web server
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		addr := fmt.Sprint(":", srv.port)
		if err := srv.ListenAndServe(addr, nil); err != http.ErrServerClosed {
			log.Fatalln("ListenAndServer error:", err)
		}
	}(srv.wg)

	err := srv.waitForShutdownCondition()

	srv.gracefulShutdown(srv.wg)

	// Return any error
	return err
}

// Blocks on incoming shutdown condition channels. If the
// shutdown trigger is an error from the controller,
// returns the error
func (srv *Server) waitForShutdownCondition() error {
	for {
		select {
		case <-srv.shutdownCalled:
			log.Println("Time to shut down!")
			return nil
		case err := <-srv.err:
			log.Println(err.Error())
			return err
		}
	}
}

// Signal all in WaitGroup to finish work and return
func (srv *Server) gracefulShutdown(wg *sync.WaitGroup) {
	// Signal server to stop servicing requests and shut down
	srv.Shutdown()

	// Wait for jobs to complete
	log.Println("Waiting for jobs to stop:", wg)
	wg.Wait()
	log.Println("=============")
	log.Println(srv.hds.GetAllHashes())
	log.Println("=============")

	log.Println("Graceful Shutdown Complete")
}
