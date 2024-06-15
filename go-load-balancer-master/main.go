package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

/*----------------------------------------RB balancing-------------------------------------------- */
type ServerRB struct {
	Addr string
}

// Balancer handles load balancing between servers
type BalancerRB struct {
	servers []*ServerRB
	index   int
}

// NewBalancer creates a new Balancer instance
func NewBalancerRB(servers []*ServerRB) *BalancerRB {
	return &BalancerRB{
		servers: servers,
		index:   0,
	}
}

// NextServer selects the next server in the pool using round robin approach
func (b *BalancerRB) NextServer() *ServerRB {
	server := b.servers[b.index]
	b.index = (b.index + 1) % len(b.servers)
	return server
}

/* --------------------------------------P2R balancing------------------------------------------- */
type ServerP2R struct {
	Addr string
}

// Balancer handles load balancing between servers using power of two choices
type BalancerP2R struct {
	servers []*ServerP2R
}

// NewBalancer creates a new Balancer instance
func NewBalancerP2R(servers []*ServerP2R) *BalancerP2R {
	return &BalancerP2R{
		servers: servers,
	}
}

// SelectServer randomly selects a server using power of two choices approach
func (b *BalancerP2R) SelectServer() *ServerP2R {
	rand.Seed(time.Now().UnixNano())

	// Choose two random servers
	n := len(b.servers)
	first := b.servers[rand.Intn(n)]
	second := b.servers[rand.Intn(n)]

	// Return the server with potentially less load (simulated)
	if rand.Intn(2) == 0 {
		return first
	}

	return second
}

/* ------------------------------Weighted Round Robin balancing------------------------------------*/

// Server represents a backend server in the pool with its weight
type Server struct {
	Addr   string
	Weight int
}

// Balancer handles load balancing between servers using weighted round robin
type Balancer struct {
	servers     []*Server
	current     int
	totalWeight int
}

// NewBalancer creates a new Balancer instance
func NewBalancer(servers []*Server) *Balancer {
	totalWeight := 0
	for _, server := range servers {
		totalWeight += server.Weight
	}
	return &Balancer{
		servers:     servers,
		current:     0,
		totalWeight: totalWeight,
	}

}

// NextServer selects a server using weighted round robin approach
func (b *Balancer) NextServer() *Server {
	rand.Seed(time.Now().UnixNano())
	target := rand.Intn(b.totalWeight)

	var weightSum int
	for i, server := range b.servers {
		weightSum += server.Weight
		if weightSum >= target {
			b.current = i
			return server
		}
	}

	// In case of uneven weight distribution, return the first server (can be improved)
	return b.servers[0]
}

/** --------------------------------------------Least Connections-------------------------------------------*/
// Server represents a backend server in the pool
type ServerLC struct {
	Addr   string
	mu     sync.Mutex
	Active int // Tracks the number of active connections
}

// Balancer interface defines the common logic for selecting a server
type BalancerLC interface {
	SelectServer() *Server
}

type LeastConnections struct {
	servers []*ServerLC
	mu      sync.Mutex
}

func NewLeastConnections(addrs []string) *LeastConnections {
	servers := make([]*ServerLC, len(addrs))
	for i, addr := range addrs {
		servers[i] = &ServerLC{Addr: addr}
	}
	return &LeastConnections{servers: servers}
}

func (lc *LeastConnections) SelectServer() *ServerLC {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// Find server with the least number of active connections
	var minActive int = math.MaxInt
	var selectedServer *ServerLC
	for _, server := range lc.servers {
		server.mu.Lock()
		active := server.Active
		server.mu.Unlock()
		if active < minActive {
			minActive = active
			selectedServer = server
		}
	}
	if selectedServer != nil {
		selectedServer.mu.Lock()
		selectedServer.Active++
		selectedServer.mu.Unlock()
	}
	return selectedServer
}

// Call RB algorithm
func callRB(requests int) [10]int {
	/** --------------------------------Stimulate RB algorithm */
	//fmt.Printf("Simulating RB algorithm\n")
	// Define RB servers
	serversRB := []*ServerRB{
		{Addr: "server0"},
		{Addr: "server1"},
		{Addr: "server2"},
		{Addr: "server3"},
		{Addr: "server4"},
		{Addr: "server5"},
		{Addr: "server6"},
		{Addr: "server7"},
		{Addr: "server8"},
		{Addr: "server9"},
	}

	balancerRB := NewBalancerRB(serversRB)
	// Simulate handling requests
	var hits = [10]int{}
	for i := 0; i < requests; i++ {
		server := balancerRB.NextServer()
		n, _ := strconv.Atoi(string(server.Addr[len(server.Addr)-1]))
		//fmt.Printf("Sending request to server: %d\n", n)
		hits[n] = hits[n] + 1
	}
	/*
		for i := 0; i < 10; i++ {
			fmt.Printf("%d", hits[i])
		}
	*/

	return hits
} //LRU

// Call P2R algorithm
func callP2R(requests int) [10]int {
	/** --------------------------------Stimulate P2R algorithm */
	//Define P2R servers
	serversP2R := []*ServerP2R{
		{Addr: "server0"},
		{Addr: "server1"},
		{Addr: "server2"},
		{Addr: "server3"},
		{Addr: "server4"},
		{Addr: "server5"},
		{Addr: "server6"},
		{Addr: "server7"},
		{Addr: "server8"},
		{Addr: "server9"},
	}

	balancerP2R := NewBalancerP2R(serversP2R)
	//fmt.Printf("\nSimulating P2B algorithm\n")
	// Simulate handling requests
	var hits = [10]int{}
	for i := 0; i < requests; i++ {
		server := balancerP2R.SelectServer()
		n, _ := strconv.Atoi(string(server.Addr[len(server.Addr)-1]))
		//fmt.Printf("Sending request to server: %d \n", server.Addr)
		hits[n] = hits[n] + 1
	}

	/*
		for i := 0; i < 10; i++ {
			fmt.Printf("%d", hits[i])
		}
	*/

	return hits

} //P2R

// Call WRB algorithm
func callWRB(requests int) [10]int {
	/** ---------------------------------Weighted Round Robin ALgorithm */
	fmt.Printf("\nSimulating weighted RB algorithm\n")
	serversWRR := []*Server{
		{Addr: "server0", Weight: 9},
		{Addr: "server1", Weight: 8},
		{Addr: "server2", Weight: 7},
		{Addr: "server3", Weight: 6},
		{Addr: "server4", Weight: 5},
		{Addr: "server5", Weight: 4},
		{Addr: "server6", Weight: 3},
		{Addr: "server7", Weight: 2},
		{Addr: "server8", Weight: 1},
		{Addr: "server9", Weight: 0},
	}

	balancerWRR := NewBalancer(serversWRR)

	// Simulate handling requests
	var hits = [10]int{}
	for i := 0; i < requests; i++ {
		server := balancerWRR.NextServer()
		n, _ := strconv.Atoi(string(server.Addr[len(server.Addr)-1]))
		//fmt.Printf("Sending request to server: %s (weight: %d)\n", server.Addr, server.Weight)
		hits[n] = hits[n] + 1
	}

	return hits
} //

// Call Least connections algorithm
func callLC(requests int) [10]int {
	/** --------------------------------Stimulate LRU algorithm */
	serversLC := []string{
		"server1",
		"server2",
		"server3",
		"server4",
		"server5",
		"server6",
		"server7",
		"server8",
		"server9",
		"server10",
	}

	// Create LeastConnections balancer
	balancer := NewLeastConnections(serversLC)

	// Simulate requests (replace with your actual logic)
	var hits = [10]int{}
	for i := 0; i < requests; i++ {
		server := balancer.SelectServer()
		// Simulate handling the request
		//fmt.Printf("Request %d sent to server %s (Active connections: %d)\n", i+1, server.Addr, server.Active)

		//Update hits
		n, _ := strconv.Atoi(string(server.Addr[len(server.Addr)-1]))
		hits[n] = hits[n] + 1

		// Simulate connection closing (replace with actual logic)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // Simulate random processing time
		server.mu.Lock()
		server.Active--
		server.mu.Unlock()
	}

	//Return
	return hits
}

func generateLoadBalancingItems(hits [10]int) []opts.BarData {
	//make room for return array
	items := make([]opts.BarData, 0)

	//update the return array
	for _, hit := range hits {
		items = append(items, opts.BarData{Value: hit})
	}

	//Return
	return items
}

func httpserverLoadBalancing(w http.ResponseWriter, _ *http.Request) {
	// create a new line instance
	bar := charts.NewBar()

	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "",
			Subtitle: "",
		}))

	// Put data into instance
	hitsRB := callRB(100000)
	hitsP2R := callP2R(100000)
	hitsWRB := callWRB(100000)
	//hitsLC := callLC(5000)
	bar.SetXAxis([]string{"Server1", "Server2", "Server3", "Server4", "Server5", "Server6", "Server7", "Server8", "Server9", "Server10"}).
		AddSeries("Round Bobin", generateLoadBalancingItems(hitsRB)).
		AddSeries("Power of 2 Random", generateLoadBalancingItems(hitsP2R)).
		AddSeries("Weighted Round Robin", generateLoadBalancingItems(hitsWRB))
	//AddSeries("Least Connections", generateLoadBalancingItems(hitsLC))

	bar.Render(w)
}

func main() {

	//create a new line instantce
	http.HandleFunc("/", httpserverLoadBalancing)

	//Call http services
	http.ListenAndServe(":8081", nil)

} //
