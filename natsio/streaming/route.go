package streaming

import (
	"github.com/nats-io/go-nats-streaming"
)

// Holds route info including subscription
// holds route,handler,nats subscripiton and queue group.
type Route struct {
	route      string
	handler    stan.MsgHandler
	subsc      stan.Subscription
	queueGroup string
}

func (r *Route) GetRoute() string {
	return r.route
}

func (r *Route) GetHandler() stan.MsgHandler {
	return r.handler
}

func (r *Route) GetSubscription() stan.Subscription {
	return r.subsc
}

func (r *Route) GetQueueGroup() string {
	return r.queueGroup
}
