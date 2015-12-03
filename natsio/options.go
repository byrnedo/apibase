package natsio
import (
	"github.com/apcera/nats"
	"time"
)

type NatsOptions struct {
	nats.Options
	routes []*Route
	encoding string
}


// Function for applying options to NatsOptions in NewNatsOptions
// Using a function allows for a chain or heirarchy when applying them
// eg func1 then func2 then func3
// Internally allows for default options to be applied first.
type OptionsFunc func(*NatsOptions) error

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

func setDefaultOptions(options *NatsOptions) error {
	options.encoding = nats.DEFAULT_ENCODER
	options.MaxReconnect = 5
	options.ReconnectWait = (2 * time.Second)
	options.Timeout = (10 * time.Second)
	options.NoRandomize = true
	return nil
}
