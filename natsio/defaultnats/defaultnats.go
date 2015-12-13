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
	"github.com/byrnedo/apibase/config/defaultconfig"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/natsio"
	"github.com/byrnedo/apibase/config"
)

var Conn *natsio.Nats

func init() {

	natsOpts := natsio.NewNatsOptions(func(n *natsio.NatsOptions) error { return nil; })

	config.Populate(&natsOpts.Options, defaultconfig.Conf, "nats")

	Info.Printf("Nats options: %#v", natsOpts)

	Info.Println("Nats encoding:", natsOpts.GetEncoding())

	var err error
	Conn, err = natsOpts.Connect()
	if err != nil {
		panic("Failed to connect to nats:" + err.Error())
	}
}
