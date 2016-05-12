package natsio

import (
	"errors"
	. "github.com/byrnedo/apibase/natsio/protobuf"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats"
	"github.com/pborman/uuid"
	"time"
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

// For use when using nats request/publish/publishrequest wrapper functions
type PayloadWithContext interface {
	proto.Message
	GetContext() *NatsContext
	SetContext(*NatsContext)
}

// Adds a context if it doesn't exist. Otherwise appends which app and time
// that this message is being sent at.
// Adds a traceID if not already there
func (n *Nats) updateContext(data PayloadWithContext, requestType RequestType) {
	var ctx *NatsContext

	if data.GetContext() == nil {
		data.SetContext(&NatsContext{})
	}

	ctx = data.GetContext()

	if len(ctx.GetTraceId()) == 0 {
		newId := uuid.NewUUID().String()
		ctx.TraceId = &newId
	}
	timeNow := time.Now().Unix()
	ctx.Trail = append(ctx.Trail, &NatsContext_Trail{&(n.Opts.Name), &requestType, &timeNow, nil})

}

// Wrapper for nats Publish function. Needs a payload which has
// a context object (see PayloadWithContext type)
// Adds a context if it doesn't exist. Otherwise appends which app and time
// that this message is being sent at.
// Adds a traceID if not already there
func (n *Nats) Publish(subject string, data PayloadWithContext) error {
	n.updateContext(data, RequestType_PUB)
	return n.EncCon.Publish(subject, data)
}

// Wrapper for nats PublishRequest function with context.
// Adds a context if it doesn't exist. Otherwise appends which app and time
// that this message is being sent at.
// Adds a traceID if not already there
func (n *Nats) PublishRequest(subject string, reply string, data PayloadWithContext) error {
	n.updateContext(data, RequestType_PUBREQ)
	return n.EncCon.PublishRequest(subject, reply, data)
}

// Wrapper for nats Request function with context.
// Adds a context if it doesn't exist. Otherwise appends which app and time
// that this message is being sent at.
// Adds a traceID if not already there
func (n *Nats) Request(subject string, data PayloadWithContext, responsePtr interface{}, timeout time.Duration) error {
	n.updateContext(data, RequestType_REQ)
	return n.EncCon.Request(subject, data, responsePtr, timeout)
}

// Unsubscribe all handlers
func (n *Nats) UnsubscribeAll() {
	for _, route := range n.Opts.routes {
		route.subsc.Unsubscribe()
	}
}
