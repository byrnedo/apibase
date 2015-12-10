package natsio

import (
	"time"
)

type RequestTypeOld string

const (
	Publish        RequestTypeOld = "PUB"
	PublishRequest RequestTypeOld = "PUBREQ"
	Request        RequestTypeOld = "REQ"
)

// Holder for trail of requests between services
type TrailOld struct {
	AppName string
	PutType RequestTypeOld
	Time    time.Time
}

// Holds info about a chain of nats requests.
type NatsContextOld struct {
	AppTrail []TrailOld
	TraceID  string
}

// Adds another link in the trail.
func (n *NatsContextOld) appendTrail(appName string, requestType RequestTypeOld) {
	n.AppTrail = append(n.AppTrail, TrailOld{appName, requestType, time.Now()})
}

// Struct used in natsio calls for request/publish. Holds context
// and an error field
type NatsDTO struct {
	NatsCtx NatsContextOld
	Error   error
}

// Retrieve a pointer to the context object
func (n *NatsDTO) Context() *NatsContextOld {
	return &n.NatsCtx
}

// Replace the context on a message
func (n *NatsDTO) NewContext(nC *NatsContextOld) {
	n.NatsCtx = *nC
}
