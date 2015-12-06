package natsio

import (
	"errors"
	"github.com/apcera/nats"
	"time"
	"github.com/pborman/uuid"
)

// nats.Options wrapper.
type Nats struct {
	Opts   *NatsOptions
	EncCon *nats.EncodedConn
}


// Subscribe and record subscription to routes
func (n *Nats) Subscribe(route string, handler nats.Handler) error {
	subsc, err := n.EncCon.Subscribe(route, handler)
	if err != nil {
		return errors.New("Failed to make subcriptions for " + route + ": " + err.Error())
	}
	n.Opts.routes = append(n.Opts.routes, &Route{route: route, handler: handler, subsc: subsc})
	return nil
}

// Subscribe to queue group and record subscription to routes
func (n *Nats) QueueSubscribe(route string, group string, handler nats.Handler) error {
	subsc, err := n.EncCon.QueueSubscribe(route, group, handler)
	if err != nil {
		return errors.New("Failed to make subcriptions for " + route + ": " + err.Error())
	}
	n.Opts.routes = append(n.Opts.routes, &Route{route: route, handler: handler, subsc: subsc, queueGroup: group})
	return nil
}

type PayloadWithContext interface {
	Context() *NatsContext
}

func (n *Nats) updateContext(data PayloadWithContext, requestType RequestType) {
	if len(data.Context().TraceID) == 0 {
		data.Context().TraceID = uuid.NewUUID().String()
	}
	data.Context().appendTrail(n.Opts.Name, requestType)
}

func (n *Nats) Publish(subject string, data PayloadWithContext) error {
	n.updateContext(data, Publish)
	return n.EncCon.Publish(subject, data)
}

func (n *Nats) PublishRequest(subject string, reply string, data PayloadWithContext) error {
	n.updateContext(data, PublishRequest)
	return n.EncCon.PublishRequest(subject, reply, data)
}

func (n *Nats) Request(subject string, data PayloadWithContext, responsePtr interface{}, timeout time.Duration) error {
	n.updateContext(data, Request)
	return n.EncCon.Request(subject, data, responsePtr, timeout)
}

// Unsubscribe all handlers
func (n *Nats) UnsubscribeAll() {
	for _, route := range n.Opts.routes {
		route.subsc.Unsubscribe()
	}
}





