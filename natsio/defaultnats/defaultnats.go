package defaultnats

import (
	"github.com/byrnedo/apibase/config"
	"github.com/byrnedo/apibase/helpers/env"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/natsio"
	"time"
)

var Conn *natsio.Nats

func init() {

	natsOpts := natsio.NewNatsOptions(func(n *natsio.NatsOptions) error {
		n.Url = env.GetOr("NATS_URL", config.Conf.GetDefaultString("nats.url", "nats://localhost:4222"))
		Info.Println("Attempting to connect to [" + n.Url + "]")
		n.Timeout = 10 * time.Second
		if appName, err := config.Conf.GetString("app-name"); err == nil && len(appName) > 0 {
			n.Name = appName
		} else {
			panic("must set app-name in conf.")
		}

		Trace.Printf("Nats Opts: %+v", n)

		return nil
	})

	Info.Println("Nats encoding:", natsOpts.GetEncoding())

	var err error
	Conn, err = natsOpts.Connect()
	if err != nil {
		panic("Failed to connect to nats:" + err.Error())
	}
}
