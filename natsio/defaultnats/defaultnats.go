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
	_ "github.com/byrnedo/apibase/logger/defaultlogger"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/natsio"
	"github.com/byrnedo/typesafe-config/parse"
	"time"
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

	var err error

	attempts := 1
	for attempts <= 5 {
		Conn, err = natsOpts.Connect()
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic("Failed to connect to nats:" + err.Error())
	}
}
