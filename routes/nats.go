package routes
import "github.com/apcera/nats"

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

