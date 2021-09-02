package main

import (
	"log"

	"github.com/jmrawlins/JCHashWebServer/controller"
)

func main() {
	controller := controller.NewController(nil, nil, nil, nil)
	if err := controller.Run(); err != nil {
		log.Fatalf("%s\n", err)
	}
}
