package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/jmrawlins/JCHashWebServer/datastore"
	"github.com/jmrawlins/JCHashWebServer/server"
)

func main() {
	// Parse arguments
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
	ds := datastore.NewMemoryHashDataStore()
	wg := &sync.WaitGroup{}
	server := server.NewServer(wg, ds, ds, shutdownCalled, errorChannel, uint16(port))
	if err := server.RunGraceful(); err != nil {
		log.Fatalf("%s\n", err)
	}
}

func usage() {
	fmt.Println("Usage:", os.Args[0], "<port>")
	os.Exit(1)
}
