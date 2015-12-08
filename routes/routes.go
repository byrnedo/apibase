package routes

// Base interface for generic route types.
type Route interface {
	GetPath() string
	GetHandler() interface{}
}
