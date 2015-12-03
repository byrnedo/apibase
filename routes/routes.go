package routes
import (
)

type Route interface {
	GetPath() string
	GetHandler() interface{}
}
