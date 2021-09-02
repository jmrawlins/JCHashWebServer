package controller

import (
	"fmt"
	"log"

	"github.com/jmrawlins/JCHashWebServer/hash/datastore"
	"github.com/jmrawlins/JCHashWebServer/webserver"
)

type Controller struct {
	shutdownChannel chan bool
	errorChannel    chan error
	ds              datastore.DataStore
	server          *webserver.Server
}

func NewController(
	shutdownChannel chan bool,
	errorChannel chan error,
	ds datastore.DataStore,
	server *webserver.Server,
) Controller {
	controller := Controller{shutdownChannel, errorChannel, ds, server}

	if shutdownChannel == nil {
		controller.shutdownChannel = make(chan bool)
	}
	if errorChannel == nil {
		controller.errorChannel = make(chan error)
	}
	if ds == nil {
		controller.ds = datastore.NewMemoryDataStore()
	}
	if server == nil {
		scheduler := webserver.NewHashJobScheduler(controller.ds)
		controller.server = webserver.NewServer(controller.ds, scheduler, controller.shutdownChannel, controller.errorChannel)
	}

	return controller
}

func (controller *Controller) Run() error {
	// Run the web server
	go controller.server.ListenAndServe(":8080", nil)

	// Wait for shutdown condition
	controller.waitForShutdown()

	// Return any error
	return nil
}

func (controller *Controller) waitForShutdown() {
WaitForErrorOrShutdown:
	for {
		select {
		case isShutdownTime := <-controller.shutdownChannel:
			if isShutdownTime {
				fmt.Println("Time to shut down!")
				fmt.Println("=============")
				fmt.Println(controller.ds.GetAllHashes())
				fmt.Println("=============")
				break WaitForErrorOrShutdown
			}
		case err := <-controller.errorChannel:
			log.Fatalln(err.Error())
			break WaitForErrorOrShutdown
		}
	}
}
