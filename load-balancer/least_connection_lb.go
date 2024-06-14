package main

import (
	"net/http"
)

type LeastConnectionLB struct {
	servers []*Server
}

func NewLeastConnectionLB(servers []*Server) *LeastConnectionLB {
	return &LeastConnectionLB{
		servers: servers,
	}
}

// ServeHTTP distributes the incoming request to the backend server with the fewest active connections.
func (lb *LeastConnectionLB) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.GetNextAvailableServer()
	logger.Print("server is %v", server)
	proxy := NewReverseProxy(server.URL)
	proxy.ServerHttp(w, r)
}

// GetNextAvailableServer returns the backend server with the fewest active connections.
func (lb *LeastConnectionLB) GetNextAvailableServer() *Server {
	minConn := -1
	var selectedServer *Server
	for _, server := range lb.servers {
		server.mutex.Lock()
		alive := server.Alive
		connections := server.Connections
		server.mutex.Unlock()

		if !alive {
			continue
		}
		if minConn == -1 || connections < minConn {
			minConn = connections
			selectedServer = server
		}
	}

	if selectedServer != nil {
		selectedServer.mutex.Lock()
		selectedServer.Connections++
		selectedServer.mutex.Unlock()
		return selectedServer
	}

	// No available servers found, return nil
	return nil
}
