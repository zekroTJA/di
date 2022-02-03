package di

import "errors"

var (
	ErrNoInterface           = errors.New("the passed type parameter does not represent an interface")
	ErrDoesNotImplInterface  = errors.New("the passed service instance does not implement the specified interface")
	ErrNotRegistered         = errors.New("no service is registered for the interface requested")
	ErrInvalidImplementation = errors.New("the retrieved service does not implement the requested interface")
)
