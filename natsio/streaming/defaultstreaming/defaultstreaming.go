package defaultstreaming

import (
	"github.com/byrnedo/apibase/config"
	"github.com/nats-io/nats"
	"github.com/byrnedo/typesafe-config/parse"
	. "github.com/byrnedo/apibase/logger"
	"time"
	"github.com/byrnedo/apibase/natsio/streaming"
	"github.com/byrnedo/capitan/helpers"
)

var StanConn *streaming.Stan

func init() {

	stanOpts := streaming.StanOptions{}
	stanOpts.Options = nats.DefaultOptions

	parse.Populate(&stanOpts, config.Conf, "stan")
	parse.Populate(&stanOpts.Options, config.Conf, "stan.nats")

	// make client id unique
	stanOpts.ClientId = stanOpts.ClientId + "-" + helpers.RandStringBytesMaskImprSrc(5)

	Info.Printf("Stan underlying Nats options: %#v", stanOpts.Options)

	stanOpts.AsyncErrorCB = func(c *nats.Conn, s *nats.Subscription, err error) {
		Error.Println("Got stan nats async error:", err)
	}

	stanOpts.DisconnectedCB = func(c *nats.Conn) {
		Warning.Println("Stan Nats disconnected")
	}

	stanOpts.ReconnectedCB = func(c *nats.Conn) {
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
			break
		} else {
			Warning.Println(err)
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic("Failed to get stan con:" + err.Error())
	}
	Info.Println("Connected to nats streaming server")
}
