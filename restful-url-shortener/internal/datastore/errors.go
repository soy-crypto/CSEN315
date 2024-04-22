package datastore

// A custom NotFoundError returned by datastore methods
type NotFoundError struct{}

var _ error = (*NotFoundError)(nil)

func (e *NotFoundError) Error() string {
	return "not found"
}

// A custon NotAllowedUser Error returned by datastore methods
type UnauthorizedError struct{}

var _ error = (*UnauthorizedError)(nil)

func (e *UnauthorizedError) Error() string {
	return "user cannot complete that action"
}
