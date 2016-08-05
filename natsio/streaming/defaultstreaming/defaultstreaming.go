package defaultstreaming

import (
	"github.com/byrnedo/apibase/config"
	"github.com/nats-io/nats"
	"github.com/byrnedo/typesafe-config/parse"
	. "github.com/byrnedo/apibase/logger"
	"time"
	"github.com/byrnedo/apibase/natsio/streaming"
)

var StanConn *streaming.Stan

func init() {

	stanOpts := streaming.StanOptions{}
	stanOpts.Options = nats.DefaultOptions

	parse.Populate(&stanOpts, config.Conf, "stan")

	natsOpts := stanOpts.Options

	Info.Printf("Stan underlying Nats options: %#v", natsOpts)

	encoding := config.Conf.GetDefaultString("stan.nats.encoding", "protobuf")

	Info.Println("Nats encoding:", encoding)

	natsOpts.AsyncErrorCB = func(c *nats.Conn, s *nats.Subscription, err error) {
		Error.Println("Got stan nats async error:", err)
	}

	natsOpts.DisconnectedCB = func(c *nats.Conn) {
		Warning.Println("Stan Nats disconnected")
	}

	natsOpts.ReconnectedCB = func(c *nats.Conn) {
		Info.Println("Stan Nats reconnected")
	}

	Info.Print("Connecting to stan nats")
	var (
		err error
	)
	attempts := 1
	for attempts <= 5 {
		attempts ++
		StanConn, err = stanOpts.Connect()
		if err == nil {
			Warning.Println(err)
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic("Failed to get stan con:" + err.Error())
	}
}
