package natsio

import (
	"errors"
"github.com/apcera/nats"
)

// nats.Options wrapper.
type Nats struct {
	Opts *NatsOptions
	EncCon *nats.EncodedConn
}


// Subscribe and record subscription to routes
func (n *Nats) Subscribe(route string, handler nats.Handler) error{
	subsc, err :=  n.EncCon.Subscribe(route, handler)
	if err != nil {
		return errors.New("Failed to make subcriptions for " + route + ": " + err.Error())
	}
	n.Opts.routes = append(n.Opts.routes, &Route{route: route, handler: handler, subsc: subsc})
	return nil
}

// Subscribe to queue group and record subscription to routes
func (n *Nats) QueueSubscribe(route string, group string, handler nats.Handler) error{
	subsc, err :=  n.EncCon.QueueSubscribe(route, group, handler)
	if err != nil {
		return errors.New("Failed to make subcriptions for " + route + ": " + err.Error())
	}
	n.Opts.routes = append(n.Opts.routes, &Route{route: route, handler: handler, subsc: subsc, queueGroup: group})
	return nil
}


// Unsubscribe all handlers
func (n *Nats) UnsubscribeAll() {
	for _, route := range n.Opts.routes {
		route.subsc.Unsubscribe()
	}
}





