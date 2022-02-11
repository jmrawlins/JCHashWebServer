package http

// dialOptions configure a Dial call. dialOptions are set by the DialOption
// values passed to Dial.
type dialOptions struct {
	unaryInt  UnaryClientInterceptor
	streamInt StreamClientInterceptor

	chainUnaryInts  []UnaryClientInterceptor
	chainStreamInts []StreamClientInterceptor

	// cp              Compressor
	// dc              Decompressor
	// bs              internalbackoff.Strategy
	// block           bool
	// returnLastError bool
	// insecure        bool
	// timeout         time.Duration
	// scChan          <-chan ServiceConfig
	// authority       string
	// copts           transport.ConnectOptions
	// callOptions     []CallOption
	// // This is used by WithBalancerName dial option.
	// balancerBuilder             balancer.Builder
	// channelzParentID            int64
	// disableServiceConfig        bool
	// disableRetry                bool
	// disableHealthCheck          bool
	// healthCheckFunc             internal.HealthChecker
	// minConnectTimeout           func() time.Duration
	// defaultServiceConfig        *ServiceConfig // defaultServiceConfig is parsed from defaultServiceConfigRawJSON.
	// defaultServiceConfigRawJSON *string
	// // This is used by ccResolverWrapper to backoff between successive calls to
	// // resolver.ResolveNow(). The user will have no need to configure this, but
	// // we need to be able to configure this in tests.
	// resolveNowBackoff func(int) time.Duration
	// resolvers         []resolver.Builder
}

// funcDialOption wraps a function that modifies dialOptions into an
// implementation of the DialOption interface.
type funcDialOption struct {
	f func(*dialOptions)
}

func (fdo *funcDialOption) apply(do *dialOptions) {
	fdo.f(do)
}

func newFuncDialOption(f func(*dialOptions)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}

// DialOption configures how we set up the connection.
type DialOption interface {
	apply(*dialOptions)
}

// WithUnaryInterceptor returns a DialOption that specifies the interceptor for
// unary RPCs.
func WithUnaryInterceptor(f UnaryClientInterceptor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.unaryInt = f
	})
}

// WithChainUnaryInterceptor returns a DialOption that specifies the chained
// interceptor for unary RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All interceptors added by this method will be chained, and the interceptor
// defined by WithUnaryInterceptor will always be prepended to the chain.
func WithChainUnaryInterceptor(interceptors ...UnaryClientInterceptor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.chainUnaryInts = append(o.chainUnaryInts, interceptors...)
	})
}

// WithStreamInterceptor returns a DialOption that specifies the interceptor for
// streaming RPCs.
func WithStreamInterceptor(f StreamClientInterceptor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.streamInt = f
	})
}

// WithChainStreamInterceptor returns a DialOption that specifies the chained
// interceptor for streaming RPCs. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All interceptors added by this method will be chained, and the interceptor
// defined by WithStreamInterceptor will always be prepended to the chain.
func WithChainStreamInterceptor(interceptors ...StreamClientInterceptor) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.chainStreamInts = append(o.chainStreamInts, interceptors...)
	})
}
