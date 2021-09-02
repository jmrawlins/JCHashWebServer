package controller

import (
	"log"

	"github.com/jmrawlins/JCHashWebServer/datastore/hashdatastore"
	"github.com/jmrawlins/JCHashWebServer/server"
	"github.com/jmrawlins/JCHashWebServer/services"
)

type Controller struct {
	shutdownChannel chan bool
	errorChannel    chan error
	ds              hashdatastore.HashDataStore
	server          *server.Server
}

func NewController(
	shutdownChannel chan bool,
	errorChannel chan error,
	ds hashdatastore.HashDataStore,
	srv *server.Server,
) Controller {
	controller := Controller{shutdownChannel, errorChannel, ds, srv}

	if shutdownChannel == nil {
		controller.shutdownChannel = make(chan bool)
	}
	if errorChannel == nil {
		controller.errorChannel = make(chan error)
	}
	if ds == nil {
		controller.ds = hashdatastore.NewMemoryDataStore()
	}

	if srv == nil {
		scheduler := services.NewHashJobScheduler(controller.ds)
		controller.server = server.NewServer(controller.ds, scheduler, controller.shutdownChannel, controller.errorChannel)
	}

	return controller
}

func (controller *Controller) Run() error {
	// Run the web server
	go controller.server.ListenAndServe(":8080", nil)

	// Wait for shutdown condition
	err := controller.waitForShutdown()

	// Signal all in WaitGroup to finish work and return
	controller.gracefulShutdown()

	// Return any error
	return err
}

func (controller *Controller) waitForShutdown() error {
	for {
		select {
		case isShutdownTime := <-controller.shutdownChannel:
			if isShutdownTime {
				log.Println("Time to shut down!")
				log.Println("=============")
				log.Println(controller.ds.GetAllHashes())
				log.Println("=============")
				return nil
			}
		case err := <-controller.errorChannel:
			log.Println(err.Error())
			return err
		}
	}
}

func (controller *Controller) gracefulShutdown() {
	// Wait for jobs to complete

}
