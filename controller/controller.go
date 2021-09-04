package controller

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/jmrawlins/JCHashWebServer/datastore"
	"github.com/jmrawlins/JCHashWebServer/server"
)

type Controller struct {
	wg              *sync.WaitGroup
	shutdownChannel chan struct{}
	errorChannel    chan error
	hds             datastore.HashDataStore
	sds             datastore.StatsDataStore
	server          *server.Server
	port            uint16
}

func NewController(
	wg *sync.WaitGroup,
	shutdownCalled chan struct{},
	errorChannel chan error,
	hds datastore.HashDataStore,
	sds datastore.StatsDataStore,
	srv *server.Server,
	port uint16,
) Controller {

	return Controller{
		wg:              wg,
		shutdownChannel: shutdownCalled,
		errorChannel:    errorChannel,
		hds:             hds,
		sds:             sds,
		server:          srv,
		port:            port,
	}

}

func (controller *Controller) waitForShutdownCommand() error {
	for {
		select {
		case <-controller.shutdownChannel:
			log.Println("Time to shut down!")
			return nil
		case err := <-controller.errorChannel:
			log.Println(err.Error())
			return err
		}
	}
}

func (controller *Controller) gracefulShutdown(wg *sync.WaitGroup) {
	// Signal server to stop servicing requests and shut down
	controller.server.Shutdown()

	// Wait for jobs to complete
	log.Println("Waiting for jobs to stop:", wg)
	wg.Wait()
	log.Println("=============")
	log.Println(controller.hds.GetAllHashes())
	log.Println("=============")

	log.Println("Graceful Shutdown Complete")
}

func (controller *Controller) Run() error {
	defer controller.wg.Wait()

	controller.wg.Add(1)
	// Run the web server
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		addr := fmt.Sprint(":", controller.port)
		if err := controller.server.ListenAndServe(addr, nil); err != http.ErrServerClosed {
			log.Fatalln("ListenAndServer error:", err)
		}
	}(controller.wg)

	// Wait for shutdown condition
	err := controller.waitForShutdownCommand()

	// Signal all in WaitGroup to finish work and return
	controller.gracefulShutdown(controller.wg)

	// Return any error
	return err
}
