package main

import (
	"flag"
	"log"
	"net/http"

	"elahi-arman.github.com/example-http-server/internal/datastore"
	"elahi-arman.github.com/example-http-server/internal/server"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// specify command line options
	port := flag.String("port", ":8080", "port to serve file server on")
	directory := flag.String("directory", "./static", "the directory of static file to host")
	flag.Parse()

	// initialize server and all its required dependencies
	// in this case the only dependency is the data store
	linkStorer, err := datastore.NewJsonFileStore("links.json")
	if err != nil {
		panic(err)
	}
	server := server.NewServer(linkStorer)

	// create a new instance of our http router, for documentation
	// refer to https://github.com/julienschmidt/httprouter
	router := httprouter.New()
	router.GET("/l/:link", server.GetLink)
	router.POST("/api/links", server.CreateLink)

	// serve all files from the directory specified (from command line arguments)
	router.ServeFiles("/public/*filepath", http.Dir(*directory))

	// start running the server. any error is considered fatal
	log.Fatal(http.ListenAndServe(*port, router))
}
