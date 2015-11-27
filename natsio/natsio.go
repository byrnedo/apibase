package natsio

import (
	"github.com/apcera/nats"
	"time"
	"errors"
)

type NatsOptions struct {
	nats.Options
	routes []*Route
	encoding string
}

// nats.Options wrapper.
type Nats struct {
	Opts *NatsOptions
	EncCon *nats.EncodedConn
}

type OptionsFunc func(*NatsOptions) error

// Holds route info
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

func prepend(slice []OptionsFunc, item OptionsFunc) []OptionsFunc{
	slice = append(slice, nil)
	copy(slice[1:], slice)
	slice[0] = item
	return slice
}

// Initiating nats with default options and then applies each
// option func in order on top of that.
func NewNatsOptions(optionFuncs  ...OptionsFunc) (options *NatsOptions) {
	options = &NatsOptions{Options: nats.DefaultOptions}
	options.setOptions(prepend(optionFuncs, setDefaultOptions)...)
	return
}

func (n *NatsOptions) SetEncoding(enc string){
	n.encoding = enc
}


func (n *NatsOptions) setOptions(optionFuncs ...OptionsFunc) error {
	for _, opt := range optionFuncs {
		if err := opt(n); err != nil {
			return err
		}
	}
	return nil
}

func setDefaultOptions(options *NatsOptions) error {
	options.encoding = nats.DEFAULT_ENCODER
	options.MaxReconnect = 5
	options.ReconnectWait = (2 * time.Second)
	options.Timeout = (10 * time.Second)
	options.NoRandomize = true
	return nil
}

// Like http.HandleFunc, give it a route and a handler (same as the normal nats subscribe)
func (n *Nats) Subscribe(route string, handler nats.Handler) error{
	subsc, err :=  n.EncCon.Subscribe(route, handler)
	if err != nil {
		return errors.New("Failed to make subcriptions for " + route + ": " + err.Error())
	}
	n.Opts.routes = append(n.Opts.routes, &Route{route: route, handler: handler, subsc: subsc})
	return nil
}

// Like http.HandleFunc, give it a route and a handler (same as the normal nats subscribe)
func (n *Nats) QueueSubscribe(route string, group string, handler nats.Handler) error{
	subsc, err :=  n.EncCon.QueueSubscribe(route, group, handler)
	if err != nil {
		return errors.New("Failed to make subcriptions for " + route + ": " + err.Error())
	}
	n.Opts.routes = append(n.Opts.routes, &Route{route: route, handler: handler, subsc: subsc, queueGroup: group})
	return nil
}

// waits 1 second before trying again <attempts> number of times
func (n *NatsOptions) ConnectOrRetry(attempts int) (natsObj *Nats,err error) {
	natsObj,err = n.Connect()
	if err != nil {
		if attempts == 1 {
			return
		}
		time.Sleep(1 * time.Second)
		natsObj, err = n.ConnectOrRetry(attempts - 1)
	}
	return
}

// Start subscribing to subjects/routes. This is non blocking.
func (natsOpts *NatsOptions) Connect() (natsObj *Nats, err error) {
	con, err := natsOpts.Options.Connect()
	if err != nil {
		return
	}

	natsObj = &Nats{Opts : natsOpts}

	natsObj.EncCon, err = nats.NewEncodedConn(con, natsOpts.encoding)
	return
}

// Get slice of Routes
func (n *NatsOptions) GetRoutes() []*Route{
	return n.routes
}

// Unsubscribe all handlers
func (n *Nats) UnsubscribeAll() {
	for _, route := range n.Opts.routes {
		route.subsc.Unsubscribe()
	}
}




