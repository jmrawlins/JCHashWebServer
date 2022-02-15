package main

import (
	"context"
	"log"
	"net/http"

	jchttp "github.com/jmrawlins/JCHashWebServer/http"
)

// UnaryServerInterceptor is called on every request received from a client to a
// unary server operation, here, we pull out the client operating system from
// the metadata, and inspect the context to receive the IP address that the
// request was received from. We then modify the EdgeLocation type to include
// this information for every request
func UnaryServerInterceptor() jchttp.UnaryServerInterceptor {
	return func(ctx context.Context, resp http.ResponseWriter, req *http.Request, info *jchttp.UnaryServerInfo, handler jchttp.UnaryHandler) {
		// Get the metadata from the incoming context
		// md, ok := metadata.FromIncomingContext(ctx)
		// if !ok {
		// 	return nil, fmt.Errorf("couldn't parse incoming context metadata")
		// }

		// // Retrieve the client OS, this will be empty if it does not exist
		// os := md.Get("client-os")
		// // Get the client IP Address
		// ip, err := getClientIP(ctx)
		// if err != nil {
		// 	return nil, err
		// }

		// Populate the EdgeLocation type with the IP and OS
		// req.(*api.EdgeLocation).IpAddress = ip
		// req.(*api.EdgeLocation).OperatingSystem = os[0]
		log.Printf("IN INTERCEPTOR: Calling through to handler")
		handler(ctx, resp, req)
		log.Printf("IN INTERCEPTOR: Handler returned")
	}
}
