package natsio

import (
	"github.com/apcera/nats"
	"time"
	"errors"
)

type Nats struct {
	*nats.Options
	routes []*Route
	encCon *nats.EncodedConn
}

type OptionsFunc func(*Nats) error

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

// Initiating nats with default options
func NewNats(optionFuncs  ...OptionsFunc) (options *Nats) {
	options = &Nats{}
	options.setOptions(prepend(optionFuncs, setDefaultOptions)...)
	return
}

func (n *Nats) setOptions(optionFuncs ...OptionsFunc) error {
	for _, opt := range optionFuncs {
		if err := opt(n); err != nil {
			return err
		}
	}
	return nil
}
func setDefaultOptions(options *Nats) error {
	options.Options = &nats.DefaultOptions
	options.MaxReconnect = 5
	options.ReconnectWait = (2 * time.Second)
	options.Timeout = (10 * time.Second)
	options.NoRandomize = true
	return nil
}

func (n *Nats) HandleFunc(route string, handler nats.Handler){
	n.routes = append(n.routes, &Route{Route: route,Handler: handler, Subsc: nil})
}

// Start subscribing to subjects/routes. This is non blocking.
func (n *Nats) ListenAndServe() error {
	con, err := n.Connect()
	if err != nil {
		return err
	}

	n.encCon, err = nats.NewEncodedConn(con, "gob")
	if err != nil {
		return err
	}

	for _, route := range n.routes {
		route.Subsc, err = n.encCon.Subscribe(route.Route, route.Handler)
		if err != nil {
			return errors.New("Failed to make subcriptions for " + route.Route + ": " + err.Error())
		}
	}
	return nil
}

func (n *Nats) GetRoutes() []*Route{
	return n.routes
}

func (n *Nats) UnsubscribeAll() {
	for _, route := range n.routes {
		route.Subsc.Unsubscribe()
	}
}

func (n *Nats) GetEncCon() *nats.EncodedConn{
	return n.encCon
}




