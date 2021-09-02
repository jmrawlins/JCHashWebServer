package router

import (
	"net/http"

	"github.com/jmrawlins/JCHashWebServer/handlers"
)

func InitRoutes(routes map[string]http.Handler) {
	for routeSpec, handler := range routes {
		superHandler := handlers.SuperHandler{ActualHandler: handler}
		http.Handle(routeSpec, superHandler)
	}
}
