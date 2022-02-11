package main

import (
	"context"
	"log"
	"reflect"

	"github.com/jmrawlins/JCHashWebServer/http"
)

// UnaryClientInterceptor is called on every request from a client to a unary
// server operation, here, we grab the operating system of the client and add it
// to the metadata within the context of the request so that it can be received
// by the server
func UnaryClientInterceptor() http.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *http.ClientConn, invoker http.UnaryInvoker, opts ...http.CallOption) error {
		// Get the operating system the client is running on
		// cos := runtime.GOOS

		// Append the OS info to the outgoing request
		// ctx = metadata.AppendToOutgoingContext(ctx, "client-os", cos)
		log.Printf("client interceptor doing upstream stuff, then calling through to invoker")

		// Invoke the original method call
		err := invoker(ctx, method, req, reply, cc, opts...)

		log.Printf("client interceptor doing downstream stuff")

		return err
	}
}

// Embedded EdgeServerStream to allow us to access the RecvMsg function on
// intercept
type EdgeServerStream struct {
	http.ServerStream
}

// RecvMsg receives messages from a stream
func (e *EdgeServerStream) RecvMsg(m interface{}) error {
	// Here we can perform additional logic on the received message, such as
	// validation
	log.Printf("intercepted server stream message, type: %s", reflect.TypeOf(m).String())
	if err := e.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}

// Set up a wrapper to allow us to access the RecvMsg function
func StreamServerInterceptor() http.StreamServerInterceptor {
	return func(srv interface{}, ss http.ServerStream, info *http.StreamServerInfo, handler http.StreamHandler) error {
		wrapper := &EdgeServerStream{
			ServerStream: ss,
		}
		return handler(srv, wrapper)
	}
}

// StreamClientInterceptor allows us to log on each client stream opening
func StreamClientInterceptor() http.StreamClientInterceptor {
	return func(ctx context.Context, desc *http.StreamDesc, cc *http.ClientConn, method string, streamer http.Streamer, opts ...http.CallOption) (http.ClientStream, error) {
		log.Printf("opening client streaming to the server method: %v", method)

		return streamer(ctx, desc, cc, method)
	}
}

// UnaryServerInterceptor is called on every request received from a client to a
// unary server operation, here, we pull out the client operating system from
// the metadata, and inspect the context to receive the IP address that the
// request was received from. We then modify the EdgeLocation type to include
// this information for every request
func UnaryServerInterceptor() http.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *http.UnaryServerInfo, handler http.UnaryHandler) (interface{}, error) {
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
		log.Printf("Calling through to handler")
		h, err := handler(ctx, req)
		// log.Printf("server interceptor hit: hydrating type with OS: '%v' and IP: '%v'", os[0], ip)

		return h, err
	}
}
