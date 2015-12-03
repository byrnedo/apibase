package routes
import "net/http"

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
