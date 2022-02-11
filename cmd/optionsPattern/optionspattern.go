package main

import "log"

// the options type that holds all the options
type thingOptions struct {
	port int
	host string
}

// ThingOption is just a func that mutates the options in some way
// We could define it simply as a function -- however, that lacks flexibility
// -- we can't add state or contextf or example.

// an object that possibly holds state
// type ThingOption func(*thingOptions)

// Create a type with an apply function -- then the option can hold state
type ThingOption interface {
	apply(*thingOptions)
}

// funcThingOption wraps a function that modifies serverOptions into an
// implementation of the ThingOption interface.
type funcThingOption struct {
	f func(*thingOptions)
}

func (fto *funcThingOption) apply(to *thingOptions) {
	fto.f(to)
}

func newFuncServerOption(f func(*thingOptions)) *funcThingOption {
	return &funcThingOption{
		f: f,
	}
}

// here's a ThingOption -- it sets the port
func WithPort(port int) ThingOption {
	return newFuncServerOption(
		func(to *thingOptions) {
			to.port = port
		})
}

// here's a ThingOption -- it sets the port
func WithLocalhost() ThingOption {
	return newFuncServerOption(
		func(to *thingOptions) {
			to.host = "localhost"
		})
}

// This is the thing that has options
type Thing struct {
	opts thingOptions
}

// Thing Factory calls all the options functions and returns a thing with those options set.
func NewThing(options ...ThingOption) *Thing {
	opts := &thingOptions{}

	for _, o := range options {
		o.apply(opts)
	}
	return &Thing{
		opts: *opts,
	}
}

func main() {
	thing1 := NewThing(WithLocalhost())
	thing2 := NewThing(WithPort(1337))

	log.Printf("Thing1: %v\n", thing1)
	log.Printf("Thing2: %v\n", thing2)
}
