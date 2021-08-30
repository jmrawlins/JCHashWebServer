package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type hashRequest struct {
	id       int
	password string
}

var (
	shutdownChannel chan bool
	hashChannel     chan hashRequest
	hashMap         map[int][]byte
	nextId          int
)

func handleStats(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Stats!\n"))
}

func handleHash(resp http.ResponseWriter, req *http.Request) {
	nextId += 1
	fmt.Fprint(resp, nextId)
	hashChannel <- hashRequest{nextId, req.FormValue("password")}
}

func handleShutdown(resp http.ResponseWriter, req *http.Request) {
	io.WriteString(resp, "Shutdown!\n")
	shutdownChannel <- true
}

func handleRetrieveHash(resp http.ResponseWriter, req *http.Request) {
	// Get the id from the uri
	strUri := strings.TrimLeft(req.URL.Path, "/")
	fmt.Println("Received request at:", strUri)

	// Get endpoint as an int, if possible
	if len(strUri) == 0 {
		fmt.Println("Up and running...")
	} else if hashId, err := strconv.ParseInt(strUri, 10, 64); err != nil {
		fmt.Println("404")
	} else {
		if hash, ok := hashMap[int(hashId)]; !ok {
			fmt.Fprint(resp, "404 hash not defined for ", hashId)
		} else {
			fmt.Printf("{ \"id\":%d, \"hash\":\"%s\"\n", hashId, base64.StdEncoding.EncodeToString(hash))
			fmt.Fprint(resp, base64.StdEncoding.EncodeToString(hash))
		}
	}
}

func init() {
	shutdownChannel = make(chan bool)
	hashChannel = make(chan hashRequest)
	hashMap = make(map[int][]byte)
}

func main() {
	fmt.Println("Serving...")
	defer func() { fmt.Println("Done") }()

	// Handle Different endpoints
	http.HandleFunc("/stats", handleStats)
	http.HandleFunc("/hash", handleHash)
	http.HandleFunc("/", handleRetrieveHash)
	http.HandleFunc("/shutdown", handleShutdown)

	fmt.Println("Calling listenAndServe")
	// Start up server
	go func() {
		defer fmt.Println("Exiting ListenAndServe goroutine")
		http.ListenAndServe(":8080", nil)
	}()
	fmt.Println("Moved beyond listenAndServe")
	// Wait for shutdown signal
ProcessChannelTraffic:
	for {
		select {
		case isShutdownTime := <-shutdownChannel:
			if isShutdownTime {
				fmt.Println("Time to shut down!")
				fmt.Println("=============")
				fmt.Println(hashMap)
				fmt.Println("=============")
				break ProcessChannelTraffic
			}

		case hashReq := <-hashChannel:
			go func() {
				time.Sleep((5 * time.Second))
				hash := sha512.Sum512([]byte(hashReq.password))
				hashMap[hashReq.id] = hash[:]

				hashStr := base64.StdEncoding.EncodeToString(hash[:])
				fmt.Println("Got password to hash: "+hashReq.password, "-- hashed is:", hashStr)
			}()
		}
	}
}
