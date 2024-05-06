package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"elahi-arman.github.com/example-http-server/internal/datastore"
	"github.com/julienschmidt/httprouter"
)

// GetLink is the function called when a user makes a request to retrieve a certain link
func (s *serverImpl) GetLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ps are the parameters attached to this route. the paramter to ByName()
	// must match the name of the link from main.go
	linkId := ps.ByName("link")

	// do some preemptive error checking
	if linkId == "" {
		fmt.Println("GetLink: no linkId provided")
		w.WriteHeader(400)
		return
	}

	// access the datastore attached to the server and try to fetch the link
	link, err := s.linkStore.GetLink(linkId)
	if errors.Is(err, &datastore.NotFoundError{}) {
		fmt.Printf("GetLink: no entry for linkId=%s\n", linkId)
		w.WriteHeader(404)
		return
	}

	// return a 302 to redirect users
	fmt.Printf("GetLink: found link for linkId=%s, redirecting to url=%s", link.Id, link.Url)
	w.Header().Add("Location", link.Url) // the location header is the destination URL
	w.WriteHeader(302)                   // 302 informs the client to read the Location header for a redirection
}

// createLinkParams represents the structure of the request body to
// a CreateLink function call
type createLinkParams struct {
	Url string `json:"url"`
	// temporary, eventually we'll replace this by retrieving from context
	Owner string `json:"owner"`
}

func (s *serverImpl) CreateLink(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// retrieve the value of the content-type header, if none is specified
	// the request should be rejected
	fmt.Printf("CREATE \n")
	contentType := r.Header.Get("content-type")
	if contentType == "" {
		fmt.Println("CreateLink: no content-type header is sent")
		w.WriteHeader(400) // the status message will automatically be filled in
		return
	}

	var url string
	var owner string
	if strings.Contains(contentType, "json") {
		// read the body of the request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("CreateLink: error while reading body of request %v\n", err)
			w.WriteHeader(400)
			return
		}

		// convert the request body into json
		lp := &createLinkParams{}
		err = json.Unmarshal(body, lp)
		if err != nil {
			fmt.Printf("CreateLink: error while unmarshalling err=%v. \n body=%s\n", err, body)
			w.WriteHeader(400)
			return
		}

		url = lp.Url
		owner = lp.Owner
	} else if strings.Contains(contentType, "form") {

		// when dealing with form data, call ParseForm to trigger parsing
		// then r.Form will have a map of the form values
		r.ParseForm()
		if formUrl, ok := r.Form["url"]; !ok || len(formUrl) == 0 || formUrl[0] == "" {
			fmt.Println("CreateLink: url key is not part of form data")
			w.Header().Add("Location", fmt.Sprintf("/public?error=%s", "cannot create a link without a url"))
			w.WriteHeader(303)
			return
		} else {
			url = formUrl[0]
		}

		if formOwner, ok := r.Form["owner"]; !ok || len(formOwner) == 0 || formOwner[0] == "" {
			fmt.Println("CreateLink: owner key is not part of form data")
			w.Header().Add("Location", fmt.Sprintf("/public?error=%s", "cannot create a link without an owner"))
			w.WriteHeader(303)
			return
		} else {
			owner = formOwner[0]
		}

	}

	// call the datastore function
	link, err := s.linkStore.CreateLink(url, owner)
	if err != nil {
		fmt.Printf("CreateLink: error while creating a link err=%v\n", err)
		w.WriteHeader(500)
		return
	}

	// redirect users
	w.Header().Add("Location", fmt.Sprintf("/public?link=%s", link.Id))
	w.WriteHeader(303)
}

// read a header / body to get a user
// return a list of links in json format where Owner == user passed in
func (s *serverImpl) GetUserLinks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ps are the parameters attached to this route. the paramter to ByName()
	// must match the name of the link from main.go
	// retrieve the value of the content-type header, if none is specified
	// the request should be rejected
	fmt.Printf("GET \n")
	contentType := r.Header.Get("content-type")
	fmt.Printf("contentType " + contentType + " \n")
	// do some preemptive error checking
	// must match the name of the link from main.go

	// access the datastore attached to the server and try to fetch the link
	userName := ps.ByName("user")
	links := s.linkStore.GetUserLinks(userName)
	if links == nil {
		fmt.Printf("GETLinks: no entry for current user \n")
		w.WriteHeader(404)
		return
	}

	// show final result
	fmt.Fprintf(w, "-------------------------------Results--------------------------------\n")
	for index, link := range links {
		// index is the index where we are
		// element is the element from someSlice for where we are
		json, err := json.Marshal(link)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintf(w, "%d   %s\n", index, string(json))
	}

	//w.Write([]byte(link))
}

func (s *serverImpl) DeleteLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Base Check
	fmt.Printf("DELETE \n")
	contentType := r.Header.Get("content-type")
	/*
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header field %q, Value %q\n", k, v)
		}
	*/

	var url string
	var owner string
	if strings.Contains(contentType, "json") {
		// read the body of the request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("GETLinks: error while reading body of request %v\n", err)
			w.WriteHeader(400)
			return
		}

		// convert the request body into json
		lp := &createLinkParams{}
		err = json.Unmarshal(body, lp)
		if err != nil {
			fmt.Printf("GETLinks: error while unmarshalling err=%v. \n body=%s\n", err, body)
			w.WriteHeader(400)
			return
		}

		url = lp.Url
		owner = lp.Owner
	} else if strings.Contains(contentType, "form") {

		// when dealing with form data, call ParseForm to trigger parsing
		// then r.Form will have a map of the form values
		r.ParseForm()
		if formUrl, ok := r.Form["url"]; !ok || len(formUrl) == 0 || formUrl[0] == "" {
			fmt.Println("GETLinks: url key is not part of form data")
			w.Header().Add("Location", fmt.Sprintf("/public?error=%s", "cannot create a link without a url"))
			w.WriteHeader(303)
			return
		} else {
			url = formUrl[0]
		}

		if formOwner, ok := r.Form["owner"]; !ok || len(formOwner) == 0 || formOwner[0] == "" {
			fmt.Println("GETLinks: owner key is not part of form data")
			w.Header().Add("Location", fmt.Sprintf("/public?error=%s", "cannot create a link without an owner"))
			w.WriteHeader(303)
			return
		} else {
			owner = formOwner[0]
		}

	}

	fmt.Printf("url: " + url + " owner: " + owner + "\n")

	//Init
	// ps are the parameters attached to this route. the paramter to ByName()
	// must match the name of the link from main.go
	user := ps.ByName("user")
	fmt.Println("\nuser:", user)

	// do some preemptive error checking
	if len(user) <= 0 {
		fmt.Println("DELETELink: no user provided")
		w.WriteHeader(400)
		return
	}

	// access the datastore attached to the server and try to fetch the link
	links, err := s.linkStore.DeleteLink(url, user)
	if errors.Is(err, &datastore.NotFoundError{}) {
		fmt.Printf("DeleteLink: no entry for user=%s\n", user)
		w.WriteHeader(404)
		return
	}

	// show final result
	fmt.Fprintf(w, "-------------------------------Results--------------------------------\n")
	for index, link := range links {
		// index is the index where we are
		// element is the element from someSlice for where we are
		json, err := json.Marshal(link)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintf(w, "%d   %s\n", index, string(json))
	}

	// return a 302 to redirect users
	w.Write([]byte("DeleteLink Success!"))

}
