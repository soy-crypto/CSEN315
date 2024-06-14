package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger

func init() {
	// Create the logger with desired settings
	logger = log.New(os.Stdout, "", log.LstdFlags)
}

func main() {
	// Define the backend servers
	servers := []*Server{{
		URL: "https://jsonplaceholder.typicode.com", Weight: 1, Alive: true},
		{URL: "https://httpbin.org", Weight: 2, Alive: true},
		{URL: "https://reqres.in", Weight: 3, Alive: true},
	}

	// Create the load balancers
	roundRobinLB := NewRoundRobinLB(servers)
	leastConnectionLB := NewLeastConnectionLB(servers)
	randomLB := NewRandomLB(servers)

	// Register the load balancers as HTTP handlers
	http.Handle("/round-robin", roundRobinLB)
	http.Handle("/least-connections", leastConnectionLB)
	http.Handle("/random", randomLB)

	// Start the server
	fmt.Println("Load balancers started.")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err.Error())
	}
}
