package streaming

import (
	"github.com/nats-io/nats"
	"github.com/nats-io/go-nats-streaming"
	"time"
	"fmt"
)

type StanOptions struct {
	nats.Options
	ClientId string
	ClusterId string
}

// Function for applying options to StanOptions in NewStanOptions
// Using a function allows for a chain or heirarchy when applying them
// eg func1 then func2 then func3
// Internally allows for default options to be applied first.
type OptionsFunc func(*StanOptions) error

func prepend(slice []OptionsFunc, item OptionsFunc) []OptionsFunc {
	slice = append(slice, nil)
	copy(slice[1:], slice)
	slice[0] = item
	return slice
}

// Initiating nats with default options and then applies each
// option func in order on top of that.
func NewStanOptions(optionFuncs ...OptionsFunc) (options *StanOptions) {
	options = &StanOptions{Options: nats.DefaultOptions}
	options.setOptions(prepend(optionFuncs, setDefaultOptions)...)
	return
}

func (n *StanOptions) setOptions(optionFuncs ...OptionsFunc) error {
	for _, opt := range optionFuncs {
		if err := opt(n); err != nil {
			return err
		}
	}
	return nil
}

// Start subscribing to subjects/routes.
func (stanOpts *StanOptions) Connect() (stan.Conn, error) {
	if len(stanOpts.ClientId) == 0 {
		panic("Must set client id in StanOptions")
	}
	if len(stanOpts.ClusterId) == 0 {
		panic("Must set cluster id in StanOptions")
	}

	con, err := stanOpts.Options.Connect()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to underlying nats: %s", err)
	}

	stanCon, err := stan.Connect(stanOpts.ClusterId, stanOpts.ClientId, stan.NatsConn(con))
	if err != nil {
		return nil, fmt.Errorf("Failed to get stan con: %s", err)
	}
	return stanCon, nil
}

func setDefaultOptions(options *StanOptions) error {
	options.MaxReconnect = 5
	options.ReconnectWait = (2 * time.Second)
	options.Timeout = (10 * time.Second)
	options.NoRandomize = true
	return nil
}
