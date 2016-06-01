package routes

import "github.com/nats-io/nats"

// Holder for a nats route's info.
type NatsRoute struct {
	path    string
	handler nats.MsgHandler
}

// Create a new route with it's subscription path and corresponding handler.
func NewNatsRoute(pathStr string, handlerFun nats.MsgHandler) *NatsRoute {
	return &NatsRoute{
		path:    pathStr,
		handler: handlerFun,
	}
}

// Get subscription path
func (n *NatsRoute) GetPath() string {
	return n.path
}

// Get Handler func
func (n *NatsRoute) GetHandler() nats.MsgHandler {
	return n.handler
}
