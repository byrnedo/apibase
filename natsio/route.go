package natsio
import "github.com/apcera/nats"

// Holds route info including subscription
type Route struct {
	route      string
	handler    nats.Handler
	subsc      *nats.Subscription
	queueGroup string
}

func (r *Route) GetRoute() string {
	return r.route
}

func (r *Route) GetHandler() nats.Handler {
	return r.handler
}

func (r *Route) GetSubscription() nats.Handler {
	return r.subsc
}

func (r *Route) GetQueueGroup() string {
	return r.queueGroup
}

