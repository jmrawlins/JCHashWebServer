package main

import (
	"fmt"
	"io"
	"net/http"
)

var (
	shutdownChannel chan bool
	hashChannel     chan string
)

func handleStats(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Stats!\n"))
}

func handleHash(resp http.ResponseWriter, req *http.Request) {
	io.WriteString(resp, "Hash!\n")
	hashChannel <- req.FormValue("password")
}

func handleShutdown(resp http.ResponseWriter, req *http.Request) {
	io.WriteString(resp, "Shutdown!\n")
	shutdownChannel <- true
}

func init() {
	shutdownChannel = make(chan bool)
	hashChannel = make(chan string)
}

func main() {
	fmt.Println("Serving...")
	defer func() { fmt.Println("Done") }()

	// Handle Different endpoints
	http.HandleFunc("/stats", handleStats)
	http.HandleFunc("/hash", handleHash)
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("I'm alive!\n"))
	})
	http.HandleFunc("/shutdown", handleShutdown)

	fmt.Println("Calling listenAndServe")
	// Start up server
	go func() {
		defer fmt.Println("Exiting ListenAndServe goroutine")
		http.ListenAndServe(":8080", nil)
	}()
	fmt.Println("Moved beyond listenAndServe")
	// Wait for shutdown signal
	for {
		select {
		case isShutdownTime := <-shutdownChannel:
			if isShutdownTime {
				fmt.Println("Time to shut down!")
				break
			}

		case password := <-hashChannel:
			fmt.Println("Got password to hash: " + password)
		}
	}
}
