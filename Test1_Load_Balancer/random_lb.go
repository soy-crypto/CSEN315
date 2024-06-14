package main

import (
	"math/rand"
	"net/http"
)

type RandomLB struct {
	servers []*Server
}

func NewRandomLB(servers []*Server) *RandomLB {
	return &RandomLB{
		servers: servers,
	}
}

// ServeHTTP distributes the incoming request to a random backend server.
func (lb *RandomLB) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.GetNextAvailableServer()
	if server != nil {
		proxy := NewReverseProxy(server.URL)
		logger.Print("server is ", server)
		proxy.ServerHttp(w, r)

	} else {
		logger.Print("server is ", server)
		//TODO:
		// Handle the case when no available server is found
	}
}

// GetNextAvailableServer returns a random backend server.
func (lb *RandomLB) GetNextAvailableServer() *Server {
	var availableServers []*Server

	for _, server := range lb.servers {
		server.mutex.Lock()

		if server.Alive {
			availableServers = append(availableServers, server)
		}
		server.mutex.Unlock()
	}
	if len(availableServers) == 0 {
		// No available servers found, return nil
		return nil
	}
	// return some random server from available servers

	return availableServers[rand.Intn(len(availableServers))]
}
