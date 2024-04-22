package server

import (
	"net/http"

	"elahi-arman.github.com/example-http-server/internal/datastore"
	"github.com/julienschmidt/httprouter"
)

type Server interface {
	GetLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	CreateLink(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
}

type serverImpl struct {
	linkStore datastore.LinkStorer
}

var _ Server = (*serverImpl)(nil)

func NewServer(ls datastore.LinkStorer) *serverImpl {
	return &serverImpl{
		linkStore: ls,
	}
}
