package natsio

import (
	"github.com/apcera/nats"
	"time"
	"errors"
)

// nats.Options wrapper.
type Nats struct {
	Opts *nats.Options
	routes []*Route
	EncCon *nats.EncodedConn
}

type OptionsFunc func(*nats.Options) error

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
func NewNats(optionFuncs  ...OptionsFunc) (options *Nats) {
	options = &Nats{Opts: &nats.DefaultOptions}
	options.setOptions(prepend(optionFuncs, setDefaultOptions)...)
	return
}

func (n *Nats) setOptions(optionFuncs ...OptionsFunc) error {
	for _, opt := range optionFuncs {
		if err := opt(n.Opts); err != nil {
			return err
		}
	}
	return nil
}

func setDefaultOptions(options *nats.Options) error {
	options.MaxReconnect = 5
	options.ReconnectWait = (2 * time.Second)
	options.Timeout = (10 * time.Second)
	options.NoRandomize = true
	return nil
}

// Like http.HandleFunc, give it a route and a handler (same as the normal nats subscribe)
func (n *Nats) HandleFunc(route string, handler nats.Handler){
	n.routes = append(n.routes, &Route{Route: route,Handler: handler, Subsc: nil})
}

// waits 1 second before trying again <attempts> number of times
func (n *Nats) ListenAndServeOrRetry(attempts int) error {
	err := n.ListenAndServe()
	if err != nil {
		if attempts == 1 {
			return err
		}
		time.Sleep(1 * time.Second)
		err = n.ListenAndServeOrRetry(attempts - 1)
	}
	return err
}

// Start subscribing to subjects/routes. This is non blocking.
func (n *Nats) ListenAndServe() error {
	con, err := n.Opts.Connect()
	if err != nil {
		return err
	}

	n.EncCon, err = nats.NewEncodedConn(con, "gob")
	if err != nil {
		return err
	}

	for _, route := range n.routes {
		route.Subsc, err = n.EncCon.Subscribe(route.Route, route.Handler)
		if err != nil {
			return errors.New("Failed to make subcriptions for " + route.Route + ": " + err.Error())
		}
	}
	return nil
}

// Get slice of Routes
func (n *Nats) GetRoutes() []*Route{
	return n.routes
}

// Unsubscribe all handlers
func (n *Nats) UnsubscribeAll() {
	for _, route := range n.routes {
		route.Subsc.Unsubscribe()
	}
}




