package routes
import (
	"github.com/apcera/nats"
	"net/http"
)

type Route interface {
	GetPath() string
	GetHandler() interface{}
}

type NatsRoute struct {
	path  string
	handler nats.Handler
}


func NewNatsRoute(pathStr string, handlerFun nats.Handler) *NatsRoute {
	return &NatsRoute{
		path: pathStr,
		handler: handlerFun,
	}
}

func (n *NatsRoute) GetPath() string {
	return n.path
}

func (n *NatsRoute) GetHandler() nats.Handler {
	return n.handler
}


type HttpMethod string

const(
	GET HttpMethod = "GET"
	POST HttpMethod = "POST"
	PUT HttpMethod = "PUT"
	PATCH HttpMethod = "PATCH"
	DELETE HttpMethod = "DELETE"
)


type WebRoute struct {
	name string
	path string
	method HttpMethod
	handler http.HandlerFunc
}

func (n *WebRoute) GetName() string {
	return n.name
}

func (n *WebRoute) GetPath() string {
	return n.path
}

func (n *WebRoute) GetMethod() string {
	return string(n.method)
}

func (n *WebRoute) GetHandler() http.HandlerFunc {
	return n.handler
}

func NewWebRoute(name string, pathStr string, method HttpMethod, handlerFun http.HandlerFunc) *WebRoute {
	return &WebRoute{
		name: name,
		path: pathStr,
		method: method,
		handler: handlerFun,
	}
}
