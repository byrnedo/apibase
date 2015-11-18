package natsio

import (
	"github.com/apcera/nats"
	"time"
)

type NatsOptions struct {
	*nats.Options
}

func (n *NatsOptions) setOptions(optionFuncs ...func(*nats.Options)) error {
	for _, opt := range optionFuncs {
		if err := opt(n); err != nil {
			return err
		}
	}
	return nil
}

func setDefaultOptions(options *nats.Options)error {
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
func InitOpts(optionFuncs  ...func(*NatsOptions)) (options *NatsOptions) {
	options = &NatsOptions{}
	options.setOptions(setDefaultOptions,optionFuncs)
	return

//	if envNatsUrl := os.Getenv("GOAX_GNATSD_URL"); envNatsUrl != "" {
//		options.Servers = []string{envNatsUrl}
//		options.Url = ""
//	} else if len(options.Servers) == 0 {
//		options.Servers = []string{options.Url}
//		options.Url = ""
//	}
}
