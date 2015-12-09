package natsio

import (
	"time"
)

type RequestType string

const (
	Publish        RequestType = "PUB"
	PublishRequest RequestType = "PUBREQ"
	Request        RequestType = "REQ"
)

// Holder for trail of requests between services
type Trail struct {
	AppName string
	PutType RequestType
	Time    time.Time
}

// Holds info about a chain of nats requests.
type NatsContext struct {
	AppTrail []Trail
	TraceID  string
}

// Adds another link in the trail.
func (n *NatsContext) appendTrail(appName string, requestType RequestType) {
	n.AppTrail = append(n.AppTrail, Trail{appName, requestType, time.Now()})
}

// Struct used in natsio calls for request/publish. Holds context
// and an error field
type NatsDTO struct {
	NatsCtx NatsContext
	Error   error
}

// Retrieve a pointer to the context object
func (n *NatsDTO) Context() *NatsContext {
	return &n.NatsCtx
}

// Replace the context on a message
func (n *NatsDTO) NewContext(nC *NatsContext) {
	n.NatsCtx = *nC
}
