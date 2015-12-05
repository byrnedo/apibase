package natsio

import (
	"errors"
"github.com/apcera/nats"
	"time"
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

func (n *Nats) Publish(subject string, data interface{}) error {
	return n.EncCon.Publish(subject, NewNatsDTO(n.Opts.appName, Publish, 0*time.Second, data))
}

func (n *Nats) PublishRequest(subject string, reply string, data interface{}) error {
	return n.EncCon.PublishRequest(subject, reply, NewNatsDTO(n.Opts.appName, Publish, 0*time.Second, data))
}

func (n *Nats) Request(subject string, data interface{}, responsePtr interface{}, timeout time.Duration) error {
	return n.EncCon.Request(subject,NewNatsDTO(n.Opts.appName, Request, timeout, data),responsePtr,timeout)
}



// Unsubscribe all handlers
func (n *Nats) UnsubscribeAll() {
	for _, route := range n.Opts.routes {
		route.subsc.Unsubscribe()
	}
}





