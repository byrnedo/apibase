package streaming

import (
	"errors"
	. "github.com/byrnedo/apibase/natsio/protobuf"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats-streaming"
	"github.com/pborman/uuid"
	"time"
)

// stan.Options wrapper.
type Stan struct {
	Opts *StanOptions
	Con  stan.Conn
}

// Subscribe and record subscription to routes
func (n *Stan) Subscribe(route string, handler stan.MsgHandler, opts ...stan.SubscriptionOption) error {
	subsc, err := n.Con.Subscribe(route, handler, opts...)
	if err != nil {
		return errors.New("Failed to make subcriptions for " + route + ": " + err.Error())
	}
	n.Opts.routes = append(n.Opts.routes, &Route{route: route, handler: handler, subsc: subsc})
	return nil
}

// Subscribe to queue group and record subscription to routes
func (n *Stan) QueueSubscribe(route string, group string, handler stan.MsgHandler, opts ...stan.SubscriptionOption) error {
	subsc, err := n.Con.QueueSubscribe(route, group, handler, opts...)
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
func (n *Stan) updateContext(data PayloadWithContext, requestType RequestType) {
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

// Wrapper for stan Publish function. Needs a payload which has
// a context object (see PayloadWithContext type)
// Adds a context if it doesn't exist. Otherwise appends which app and time
// that this message is being sent at.
// Adds a traceID if not already there
func (n *Stan) Publish(subject string, data PayloadWithContext) error {
	n.updateContext(data, RequestType_PUB)
	bData, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	return n.Con.Publish(subject, bData)
}

// Wrapper for stan PublishAsync function with context.
// Adds a context if it doesn't exist. Otherwise appends which app and time
// that this message is being sent at.
// Adds a traceID if not already there
func (n *Stan) PublishAsync(subject string, data PayloadWithContext, ackH stan.AckHandler) (string, error) {
	n.updateContext(data, RequestType_PUB)
	bData, err := proto.Marshal(data)
	if err != nil {
		return "", err
	}
	return n.Con.PublishAsync(subject, bData, ackH)
}

// Unsubscribe all handlers
func (n *Stan) UnsubscribeAll() {
	for _, route := range n.Opts.routes {
		route.subsc.Unsubscribe()
	}
}
