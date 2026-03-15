// Package httpclient provides centralized HTTP client configuration for GitHub API requests.
package httpclient

import (
	"net/http"
	"time"
)

const (
	// DefaultTimeout is the default timeout for HTTP requests
	DefaultTimeout = 30 * time.Second

	// MaxIdleConnections is the maximum number of idle connections across all hosts
	MaxIdleConnections = 100

	// MaxIdleConnectionsPerHost is the maximum number of idle connections per host
	MaxIdleConnectionsPerHost = 10

	// IdleConnectionTimeout is how long an idle connection is kept alive
	IdleConnectionTimeout = 90 * time.Second
)

// defaultClient is a singleton HTTP client with optimized connection pooling
var defaultClient = &http.Client{
	Timeout: DefaultTimeout,
	Transport: &http.Transport{
		MaxIdleConns:        MaxIdleConnections,
		MaxIdleConnsPerHost: MaxIdleConnectionsPerHost,
		IdleConnTimeout:     IdleConnectionTimeout,
	},
}

// Default returns the default HTTP client with connection pooling and timeout configured.
// This client is safe for concurrent use and should be reused across requests.
func Default() *http.Client {
	return defaultClient
}

// New creates a new HTTP client with custom timeout while maintaining connection pooling.
// Use this when you need a different timeout than the default.
func New(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        MaxIdleConnections,
			MaxIdleConnsPerHost: MaxIdleConnectionsPerHost,
			IdleConnTimeout:     IdleConnectionTimeout,
		},
	}
}
