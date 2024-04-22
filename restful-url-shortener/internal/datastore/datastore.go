package datastore

import "time"

// Link is the in memory representation of a link. Not all of thes
// fields should be exposed externally. The `json:"name"` on each
// line refers to the field name when this structure is serialzied
// into json.
type Link struct {
	Id        string    `json:"id"`
	Url       string    `json:"url"`
	Owner     string    `json:"owner"`
	Views     int64     `json:"views"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LinkStorer is an interface which defines the set of operations
// a struct must fulfill to be considered a LinkStorer
type LinkStorer interface {
	GetLink(id string) (*Link, error)
	GetUserLinks(user string) []Link
	CreateLink(url string, owner string) (*Link, error)
	DeleteLink(url string, user string) error
}
