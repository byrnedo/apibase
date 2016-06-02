// A default nats initialization. Reads config from
// the default config file, thus depends on config
// being ready at this point ( initialized in it's own
// init function )
//
// This uses config.Populate to populate the options struct
// directly from the config.
// The config section would look as follows:
// nats {
//
// }
//
//
package defaultnats

import (
	"github.com/byrnedo/apibase/config"
	. "github.com/byrnedo/apibase/logger"
	_ "github.com/byrnedo/apibase/logger/defaultlogger"
	"github.com/byrnedo/apibase/natsio"
	"github.com/byrnedo/typesafe-config/parse"
	"time"
	"github.com/nats-io/nats"
)

var Conn *natsio.Nats

func init() {

	natsOpts := natsio.NewNatsOptions(func(n *natsio.NatsOptions) error { return nil })

	natsOpts.Options.MaxReconnect = 15

	parse.Populate(&natsOpts.Options, config.Conf, "nats")

	Info.Printf("Nats options: %#v", natsOpts)

	encoding := config.Conf.GetDefaultString("nats.encoding", natsOpts.GetEncoding())

	Info.Println("Nats encoding:", encoding)
	natsOpts.SetEncoding(encoding)

	natsOpts.Options.AsyncErrorCB = func(c *nats.Conn, s *nats.Subscription, err error) {
		Error.Println("Got nats async error:", err)
	}

	natsOpts.Options.DisconnectedCB = func(c *nats.Conn) {
		Warning.Println("Nats disconnected")
	}

	natsOpts.Options.ReconnectedCB = func(c *nats.Conn) {
		Info.Println("Nats reconnected")
	}

	var err error

	Info.Print("Connecting to nats")
	attempts := 1
	for attempts <= 5 {
		attempts ++
		Conn, err = natsOpts.Connect()
		if err == nil {
			Warning.Println(err)
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic("Failed to connect to nats:" + err.Error())
	}
	Info.Print("Connected to nats")
}
