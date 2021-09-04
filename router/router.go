package router

import (
	"net/http"

	"github.com/jmrawlins/JCHashWebServer/datastore"
	"github.com/jmrawlins/JCHashWebServer/handlers"
)

func InitRoutes(sds datastore.StatsDataStore, routes map[string]http.Handler) {
	for routeSpec, handler := range routes {
		superHandler := handlers.SuperHandler{ActualHandler: handler, Sds: sds}
		http.Handle(routeSpec, superHandler)
	}
}
