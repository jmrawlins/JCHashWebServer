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
	if len(args) != 2 {
		usage()
	}

	// Parse port as 16-bit uint (max 65535)
	port, err := strconv.ParseUint(args[0], 10, 16)
	if err != nil {
		usage()
	}

	// Open or create hashes file
	filename := args[1]
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		log.Fatalf("Error creating/opening hashes file (%s) for rw: %s", filename, err.Error())
	}
	mds := datastore.NewMemoryHashDataStore()

	fds, err := datastore.NewFileHashDataStore(file, mds)
	if err != nil {
		log.Fatalf("problem initializing file data store: %s", err.Error())
	}

	shutdownCalled := make(chan struct{})
	errorChannel := make(chan error)
	wg := &sync.WaitGroup{}
	server := server.NewServer(wg, fds, mds, shutdownCalled, errorChannel, uint16(port))
	if err := server.RunGraceful(); err != nil {
		log.Fatalf("%s\n", err)
	}

	// Close hashes file
	file.Close()
}

func usage() {
	fmt.Println("Usage:", os.Args[0], "<port> <hashfile>")
	os.Exit(1)
}
