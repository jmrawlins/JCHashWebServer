package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/jmrawlins/JCHashWebServer/datastore"
	"github.com/jmrawlins/JCHashWebServer/server"
)

func main() {
	// Parse arguments -- only have 1, it's a port number
	args := os.Args[1:]
	if len(args) != 1 {
		usage()
	}
	port, err := strconv.ParseUint(args[0], 10, 16)
	if err != nil {
		usage()
	}

	shutdownCalled := make(chan struct{})
	errorChannel := make(chan error)
	ds := datastore.NewMemoryDataStore()
	wg := &sync.WaitGroup{}
	server := server.NewServer(wg, ds, ds, shutdownCalled, errorChannel, uint16(port))
	if err := server.Run(); err != nil {
		log.Fatalf("%s\n", err)
	}
}

func usage() {
	log.Fatalln("Usage:", os.Args[0], "<port>")
}
