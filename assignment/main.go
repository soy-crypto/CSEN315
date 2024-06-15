package main

import (
	"fmt"
	"math/rand"
	"net/http"
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

/** ---------------------------------------------LRU-------------------------------------------*/
// Server represents a backend server in the pool
type ServerLRU struct {
	Addr string
}

// Node represents a node in the linked list for LRU balancing
type Node struct {
	Server *ServerLRU
	Next   *Node
	Prev   *Node
}

// LRU balancer maintains the server list and access history
type LRU struct {
	head  *Node
	tail  *Node
	cache map[string]*Node // Map server address to its corresponding node in the list
}

// NewLRU creates a new LRU balancer instance
func NewLRU() *LRU {
	return &LRU{
		cache: make(map[string]*Node),
	}
}

// AddServer adds a server to the LRU pool
func (lru *LRU) AddServer(server *ServerLRU) {
	newNode := &Node{Server: server}
	if lru.head == nil {
		lru.head = newNode
		lru.tail = newNode
	} else {
		newNode.Next = lru.head
		lru.head.Prev = newNode
		lru.head = newNode
	}
	lru.cache[server.Addr] = newNode
}

// GetServer retrieves the least recently contacted server and updates its position
func (lru *LRU) GetServer(addr string) *ServerLRU {
	node, ok := lru.cache[addr]
	if !ok {
		return nil
	}

	if node == lru.head {
		return node.Server
	}

	// Move the node to the head, indicating recent access
	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}
	if lru.tail == node {
		lru.tail = node.Prev
	}

	node.Next = lru.head
	node.Prev = nil
	lru.head.Prev = node
	lru.head = node
	return node.Server

} //LRU

// call RB algorithm
func callRB(requests int) [10]int {
	/** --------------------------------Stimulate RB algorithm */
	fmt.Printf("Simulating RB algorithm\n")
	// Define RB servers
	serversRB := []*ServerRB{
		{Addr: "server1"},
		{Addr: "server2"},
		{Addr: "server3"},
		{Addr: "server4"},
		{Addr: "server5"},
		{Addr: "server6"},
		{Addr: "server7"},
		{Addr: "server8"},
		{Addr: "server9"},
		{Addr: "server10"},
	}

	balancerRB := NewBalancerRB(serversRB)
	// Simulate handling requests
	var hits = [10]int{}
	for i := 0; i < requests; i++ {
		server := balancerRB.NextServer()
		fmt.Printf("Sending request to server: %s\n", server.Addr)
		hits[server.Addr[len(server.Addr)-1]] = hits[server.Addr[len(server.Addr)-1]] + 1
	}

	return hits
} //LRU

// call P2R algorithm
func callP2R(requests int) [10]int {
	/** --------------------------------Stimulate P2R algorithm */
	//Define P2R servers
	serversP2R := []*ServerP2R{
		{Addr: "server1"},
		{Addr: "server2"},
		{Addr: "server3"},
		{Addr: "server4"},
		{Addr: "server5"},
		{Addr: "server6"},
		{Addr: "server7"},
		{Addr: "server8"},
		{Addr: "server9"},
		{Addr: "server10"},
	}

	balancerP2R := NewBalancerP2R(serversP2R)
	fmt.Printf("\nSimulating P2B algorithm\n")
	// Simulate handling requests
	var hits = [10]int{}
	for i := 0; i < requests; i++ {
		server := balancerP2R.SelectServer()
		fmt.Printf("Sending request to server: %s\n", server.Addr)
		hits[server.Addr[len(server.Addr)-1]] = hits[server.Addr[len(server.Addr)-1]] + 1
	}

	return hits

} //LRU

// call WRB algorithm
func callWRB(requests int) [10]int {
	/** ---------------------------------Weighted Round Robin ALgorithm */
	fmt.Printf("\nSimulating weighted RB algorithm\n")
	serversWRR := []*Server{
		{Addr: "server1"},
		{Addr: "server2"},
		{Addr: "server3"},
		{Addr: "server4"},
		{Addr: "server5"},
		{Addr: "server6"},
		{Addr: "server7"},
		{Addr: "server8"},
		{Addr: "server9"},
		{Addr: "server10"},
	}

	balancerWRR := NewBalancer(serversWRR)

	// Simulate handling requests
	var hits = [10]int{}
	for i := 0; i < requests; i++ {
		server := balancerWRR.NextServer()
		fmt.Printf("Sending request to server: %s (weight: %d)\n", server.Addr, server.Weight)
		hits[server.Addr[len(server.Addr)-1]] = hits[server.Addr[len(server.Addr)-1]] + 1
	}

	return hits
} //

// call LRU algorithm
func callLRU() {
	/** --------------------------------Stimulate LRU algorithm */
	fmt.Printf("\n Simulating LRU algorithm\n")
	// Define some backend servers
	servers := []*ServerLRU{
		{Addr: "server1:8080"},
		{Addr: "server2:8080"},
		{Addr: "server3:8080"},
	}

	balancer := NewLRU()

	// Add servers to the LRU pool
	for _, server := range servers {
		balancer.AddServer(server)
	}

	// Simulate handling requests
	for i := 0; i < 5; i++ {
		// Access server2 twice to simulate recent contact
		server := balancer.GetServer("server2:8080")
		if server != nil {
			fmt.Printf("Sending request (access 1) to server: %s\n", server.Addr)
		}
		server = balancer.GetServer("server2:8080")
		if server != nil {
			fmt.Printf("Sending request (access 2) to server: %s\n", server.Addr)
		}

		// Access a different server
		server = balancer.GetServer("server1:8080")
		if server != nil {
			fmt.Printf("Sending request to server: %s\n", server.Addr)
		}

	}

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

func httpserverHome(w http.ResponseWriter, _ *http.Request) {

}

func httpserverLoadBalancing(w http.ResponseWriter, _ *http.Request) {
	// create a new line instance
	bar := charts.NewBar()

	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line example in Westeros theme",
			Subtitle: "Line chart rendered by the http server this time",
		}))

	// Put data into instance
	bar.SetXAxis([]string{"Server1", "Server2", "Server3", "Server4", "Server5", "Server6", "Server7", "Server8", "Server9", "Server10"}).
		AddSeries("Category Round Bobin", generateLoadBalancingItems(callRB(10000))).
		AddSeries("Category Power of 2 Random", generateLoadBalancingItems(callP2R(10000))).
		AddSeries("Category Weighted Round Robin", generateLoadBalancingItems(callWRB(10000)))

	bar.Render(w)

}

func main() {

	//create a new line instantce
	http.HandleFunc("/", httpserverHome)
	http.HandleFunc("/lb", httpserverLoadBalancing)

	//Call http services
	http.ListenAndServe(":8081", nil)

} //
