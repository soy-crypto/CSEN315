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

// Server represents a backend server with its weight (for WRR)
type Server struct {
	Addr   string
	Weight int
	Hits   int // Added to track access count
}

// Balancer interface defines the common logic for selecting a server
type Balancer interface {
	SelectServer() *Server
}

// RBBalancer implements Round Robin balancing
type RBBalancer struct {
	servers []*Server
	index   int
}

func (b *RBBalancer) SelectServer() *Server {
	server := b.servers[b.index]
	b.index = (b.index + 1) % len(b.servers)
	server.Hits++ // Track hits for RB
	return server
}

// P2RBalancer implements Power of Two Random Choices balancing
type P2RBalancer struct {
	servers []*Server
}

func (b *P2RBalancer) SelectServer() *Server {
	rand.Seed(time.Now().UnixNano())
	n := len(b.servers)
	first := b.servers[rand.Intn(n)]
	second := b.servers[rand.Intn(n)]
	if rand.Intn(2) == 0 {
		first.Hits++ // Track hits for P2R
		return first
	}
	second.Hits++ // Track hits for P2R
	return second
}

// LRUBalancer implements Least Recently Used balancing
type LRUBalancer struct {
	head  *Node
	tail  *Node
	cache map[string]*Node
}

type Node struct {
	Server *Server
	Next   *Node
	Prev   *Node
}

func (b *LRUBalancer) SelectServer(addr string) *Server {
	if b.head == nil {
		return nil
	}

	node, ok := b.cache[addr]
	if !ok {
		return nil
	}

	if node == b.head {
		return node.Server
	}

	// Move the node to the head
	b.moveToHead(node)
	return node.Server
}

func (b *LRUBalancer) moveToHead(node *Node) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}
	if b.tail == node {
		b.tail = node.Prev
	}

	node.Next = b.head
	node.Prev = nil
	b.head.Prev = node
	b.head = node
}

// WRRBalancer implements Weighted Round Robin balancing
type WRRBalancer struct {
	servers     []*Server
	current     int
	totalWeight int
}

func (b *WRRBalancer) SelectServer() *Server {
	rand.Seed(time.Now().UnixNano())
	target := rand.Intn(b.totalWeight)

	var weightSum int
	for i, server := range b.servers {
		weightSum += server.Weight
		if weightSum >= target {
			b.current = i
			server.Hits++ // Track hits for WRR
			return server
		}
	}

	return b.servers[0] // Handle uneven weight distribution (can be improved)
}

func selectBalancer(choice int) Balancer {
	servers := []*Server{
		{Addr: "server1", Weight: 2, Hits: 0},
		{Addr: "server2", Weight: 1, Hits: 0},
		{Addr: "server3", Weight: 3, Hits: 0},
	}

	switch choice {
	case 0:
		return &RBBalancer{servers: servers}

	case 1:
		return &P2RBalancer{servers: servers}

	case 2:
		return &WRRBalancer{servers: servers}

	default:
		fmt.Println("Invalid choice. Please select between 0 (RB), 1 (P2R),or 3 (WRR).")
		return nil

	}

}

/*
func main() {
	// Simulate requests (replace with your actual logic)
	requests := 100

	// Track server hits for all algorithms
	var rbHits, p2rHits, lruHits, wrrHits int

	for i := 0; i < requests; i++ {
		choice := rand.Intn(4) // Choose a random balancing algorithm
		balancer := selectBalancer(choice)

		if balancer != nil {
			server := balancer.SelectServer()
			switch choice {
			case 0:
				rbHits += server.Hits
			case 1:
				p2rHits += server.Hits
			case 2:
				lruHits += 1 // LRU doesn't directly track hits, increment for each selection
			case 3:
				wrrHits += server.Hits

			}

		}

	}

	// Generate charts using go-echarts
	e := render.NewEngine()
	opt := opts.NewBarOpts().
		SetLegend(opts.LegendOpts{Data: []string{"RB", "P2R", "WRR"}}).
		SetXAxis(opts.XAxisOpts{Data: []string{"Server Hits"}})
	opt.AddYAxis("Hits", opts.YAxisOpts{Name: "Number of Requests"})
	opt.AddSeries("RB", opts.BarOpts{Data: []int{rbHits}})
	opt.AddSeries("P2R", opts.BarOpts{Data: []int{p2rHits}})
	opt.AddSeries("WRR", opts.BarOpts{Data: []int{wrrHits}})
	err := e.Render(opt)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Open the rendered chart in the browser (replace with your preferred method for visualization)
	fmt.Println("Chart generated, open http://localhost:8080 to view")
	err = e.Serve(fmt.Sprintf(":%d", 8080))
	if err != nil {
		fmt.Println(err)
	}

}

*/

// generate random data for bar chart
func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}

// generate random data for line chart
func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

func generateLoadBalancingItems(rb int, p2r int, wrb int) []opts.BarData {
	items := make([]opts.BarData, 0)
	items = append(items, opts.BarData{Value: rb})
	items = append(items, opts.BarData{Value: p2r})
	items = append(items, opts.BarData{Value: wrb})
	return items
}

func httpserverHome(w http.ResponseWriter, _ *http.Request) {

}

func httpserverLine(w http.ResponseWriter, _ *http.Request) {
	// create a new line instance
	line := charts.NewLine()

	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line example in Westeros theme",
			Subtitle: "Line chart rendered by the http server this time",
		}))

	// Put data into instance
	line.SetXAxis([]string{"RB", "P2R", "WRR"}).
		AddSeries("Category Round Robin ", generateLineItems()).
		AddSeries("Category Power of 2 Random ", generateLineItems()).
		AddSeries("Category Weighted Round Robin ", generateLineItems()).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	line.Render(w)
}

func httpserverBar(w http.ResponseWriter, _ *http.Request) {
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
	bar.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Category A", generateBarItems()).
		AddSeries("Category B", generateBarItems())

	bar.Render(w)
}

func httpserverLoadBalancing(w http.ResponseWriter, _ *http.Request) {
	/** Create all the data */
	// Simulate requests (replace with your actual logic)
	requests := 100

	// Track server hits for all algorithms
	var rbHits, p2rHits, wrrHits int

	for i := 0; i < requests; i++ {
		choice := rand.Intn(3) // Choose a random balancing algorithm
		balancer := selectBalancer(choice)

		if balancer != nil {
			server := balancer.SelectServer()
			switch choice {
			case 0:
				rbHits += server.Hits
			case 1:
				p2rHits += server.Hits
			case 2:
				wrrHits += server.Hits
			}

		} // if

	} // for

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
	bar.SetXAxis([]string{"RB", "P2R", "WRR"}).
		AddSeries("Category A", generateLoadBalancingItems(rbHits, p2rHits, wrrHits)).
		AddSeries("Category B", generateLoadBalancingItems(rbHits, p2rHits, wrrHits)).
		AddSeries("Category C", generateLoadBalancingItems(rbHits, p2rHits, wrrHits))

	bar.Render(w)

}

func main() {

	//create a new line instantce
	http.HandleFunc("/", httpserverHome)
	http.HandleFunc("/bar", httpserverBar)
	http.HandleFunc("/line", httpserverLine)
	http.HandleFunc("/lb", httpserverLoadBalancing)

	//Call http services
	http.ListenAndServe(":8081", nil)
}
