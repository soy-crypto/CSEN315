package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

// LoadBalancer defines the interface for a load balancer.
type LoadBalancer interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
	GetNextAvailableServer() *Server
}

// Server represents a backend server.
type Server struct {
	URL         string
	Alive       bool
	Weight      int
	Connections int
	mutex       sync.Mutex // using it to protect concurrent access to alive and connections field
}

type ReverseProxy struct {
	backendURL string
	proxy      *httputil.ReverseProxy
}

func NewReverseProxy(backendURL string) *ReverseProxy {
	backend, _ := url.Parse(backendURL)

	return &ReverseProxy{
		backendURL: backendURL,
		proxy:      httputil.NewSingleHostReverseProxy(backend),
	}

}

// Forwards the incoming request to backend server
func (rp *ReverseProxy) ServerHttp(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Forwarding request to %s : %s\n", rp.backendURL, r.URL.Path)
	rp.proxy.ServeHTTP(w, r)
}
