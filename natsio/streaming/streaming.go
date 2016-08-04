package streaming

import "github.com/nats-io/nats"

type StanOptions struct {
	nats.Options
	ClientId string
	ClusterId string
}
