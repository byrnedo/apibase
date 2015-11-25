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
	Route   string
	Handler nats.Handler
	Subsc   *nats.Subscription
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
func (n *NatsOptions) HandleFunc(route string, handler nats.Handler){
	n.routes = append(n.routes, &Route{Route: route,Handler: handler, Subsc: nil})
}

// waits 1 second before trying again <attempts> number of times
func (n *NatsOptions) ListenAndServeOrRetry(attempts int) (natsObj *Nats,err error) {
	natsObj,err = n.ListenAndServe()
	if err != nil {
		if attempts == 1 {
			return
		}
		time.Sleep(1 * time.Second)
		natsObj, err = n.ListenAndServeOrRetry(attempts - 1)
	}
	return
}

// Start subscribing to subjects/routes. This is non blocking.
func (natsOpts *NatsOptions) ListenAndServe() (natsObj *Nats, err error) {
	con, err := natsOpts.Connect()
	if err != nil {
		return
	}

	natsObj = &Nats{Opts : natsOpts}

	natsObj.EncCon, err = nats.NewEncodedConn(con, natsOpts.encoding)
	if err != nil {
		return
	}

	for _, route := range natsOpts.routes {
		route.Subsc, err = natsObj.EncCon.Subscribe(route.Route, route.Handler)
		if err != nil {
			return natsObj, errors.New("Failed to make subcriptions for " + route.Route + ": " + err.Error())
		}
	}
	return
}

// Get slice of Routes
func (n *NatsOptions) GetRoutes() []*Route{
	return n.routes
}

// Unsubscribe all handlers
func (n *Nats) UnsubscribeAll() {
	for _, route := range n.Opts.routes {
		route.Subsc.Unsubscribe()
	}
}




