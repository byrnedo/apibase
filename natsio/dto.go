package natsio
import (
	"time"
	"github.com/pborman/uuid"
)

type RequestType string

const (
	Publish RequestType = "PUB"
	PubSub RequestType = "PUBSUB"
	Request RequestType = "REQ"
)

type Context struct {
	PutAppName string
	PutTime time.Time
	PutType RequestType
	Timeout time.Duration
	TraceID string
}

type NatsDTO struct {
	Context Context
	Payload interface{}
}

func NewNatsDTO(caller string, reqType RequestType, timeout time.Duration, payload interface{}) *NatsDTO {
	return &NatsDTO{
		Context: Context{
			PutAppName: caller,
			PutTime: time.Now(),
			PutType: reqType,
			Timeout: timeout,
			TraceID: uuid.NewUUID().String(),
		},
		Payload: payload,
	}
}