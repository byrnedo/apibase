package natsio

import (
	"github.com/apcera/nats"
	"time"
)

type Nats struct {
	*nats.Options
	routes []Route
}

type Route struct {
	route   string
	handler func(*nats.Conn, *nats.EncodedConn, nats.Msg)
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

// InitOpts : Initiating nats with default options
func NewNats(routes []Route, optionFuncs  ...func(*Nats)) (options *Nats) {
	options = &Nats{}
	options.setOptions(setDefaultOptions, optionFuncs)
	options.routes = routes
	return

	//	if envNatsUrl := os.Getenv("GOAX_GNATSD_URL"); envNatsUrl != "" {
	//		options.Servers = []string{envNatsUrl}
	//		options.Url = ""
	//	} else if len(options.Servers) == 0 {
	//		options.Servers = []string{options.Url}
	//		options.Url = ""
	//	}
}


