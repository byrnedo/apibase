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

type Trail struct {
	AppName string
	Time time.Time
}

type NatsContext struct {
	PutAppName string
	Trail []Trail
	PutTime time.Time
	PutType RequestType
	Timeout time.Duration
	TraceID string
}

func (n *NatsContext) AppendTrail(appName string){
	n.Trail = append(n.Trail, Trail{appName, time.Now()})
}


func NewNatsContext(caller string, reqType RequestType, timeout time.Duration) *NatsContext {
	return &NatsContext{
			PutAppName: caller,
			PutTime: time.Now(),
			PutType: reqType,
			Timeout: timeout,
			TraceID: uuid.NewUUID().String(),
		}

}