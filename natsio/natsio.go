package natsio

import (
	"github.com/apcera/nats"
	"time"
	"errors"
)

type Nats struct {
	*nats.Options
	routes []*Route
}

type Route struct {
	Route   string
	Handler nats.Handler
	Subsc   *nats.Subscription
}

// Initiating nats with default options
func NewNats(optionFuncs  ...func(*Nats)) (options *Nats) {
	options = &Nats{}
	options.setOptions(setDefaultOptions, optionFuncs)
	return
}

func (n *Nats) setOptions(optionFuncs ...func(*nats.Options)) error {
	for _, opt := range optionFuncs {
		if err := opt(n); err != nil {
			return err
		}
	}
	return nil
}
func setDefaultOptions(options *nats.Options) error {
	// Optionally set ReconnectWait and MaxReconnect attempts.
	// This example means 10 seconds total per backend.
	options = nats.DefaultOptions
	options.MaxReconnect = 5
	options.ReconnectWait = (2 * time.Second)
	options.Timeout = (10 * time.Second)
	// Optionally disable randomization of the server pool
	options.NoRandomize = true
	return nil
}

func (n *Nats) HandleFunc(route string, handler nats.Handler){
	n.routes = append(n.routes, &Route{route, handler})
}

func (n *Nats) ListenAndServe() error {
	con, err := n.Connect()
	if err != nil {
		return err
	}

	encCon, err := nats.NewEncodedConn(con, "gob")
	if err != nil {
		return err
	}

	for _, route := range n.routes {
		route.Route, err = encCon.Subscribe(route.Route, route.Handler)
		if err != nil {
			return errors.New("Failed to make subcriptions for " + route.Route + ": " + err.Error())
		}
	}
	return nil
}

func (n *Nats) GetRoutes() []*Route{
	return n.routes
}




