package httpclient

import (
	"net/http"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	client := Default()

	if client == nil {
		t.Fatal("Default() returned nil client")
	}

	if client.Timeout != DefaultTimeout {
		t.Errorf("Default client timeout = %v, want %v", client.Timeout, DefaultTimeout)
	}

	// Verify it returns the same instance (singleton)
	client2 := Default()
	if client != client2 {
		t.Error("Default() should return the same instance")
	}
}

func TestNew(t *testing.T) {
	customTimeout := 60 * time.Second
	client := New(customTimeout)

	if client == nil {
		t.Fatal("New() returned nil client")
	}

	if client.Timeout != customTimeout {
		t.Errorf("New client timeout = %v, want %v", client.Timeout, customTimeout)
	}

	// Verify each call creates a new instance
	client2 := New(customTimeout)
	if client == client2 {
		t.Error("New() should create a new instance each time")
	}
}

func TestDefaultClientConfiguration(t *testing.T) {
	client := Default()

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("Default client transport is not *http.Transport")
	}

	if transport.MaxIdleConns != MaxIdleConnections {
		t.Errorf("MaxIdleConns = %d, want %d", transport.MaxIdleConns, MaxIdleConnections)
	}

	if transport.MaxIdleConnsPerHost != MaxIdleConnectionsPerHost {
		t.Errorf("MaxIdleConnsPerHost = %d, want %d", transport.MaxIdleConnsPerHost, MaxIdleConnectionsPerHost)
	}

	if transport.IdleConnTimeout != IdleConnectionTimeout {
		t.Errorf("IdleConnTimeout = %v, want %v", transport.IdleConnTimeout, IdleConnectionTimeout)
	}
}

func TestNewClientConfiguration(t *testing.T) {
	customTimeout := 45 * time.Second
	client := New(customTimeout)

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatal("New client transport is not *http.Transport")
	}

	if transport.MaxIdleConns != MaxIdleConnections {
		t.Errorf("MaxIdleConns = %d, want %d", transport.MaxIdleConns, MaxIdleConnections)
	}

	if transport.MaxIdleConnsPerHost != MaxIdleConnectionsPerHost {
		t.Errorf("MaxIdleConnsPerHost = %d, want %d", transport.MaxIdleConnsPerHost, MaxIdleConnectionsPerHost)
	}

	if transport.IdleConnTimeout != IdleConnectionTimeout {
		t.Errorf("IdleConnTimeout = %v, want %v", transport.IdleConnTimeout, IdleConnectionTimeout)
	}
}
