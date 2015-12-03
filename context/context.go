package context
import "time"

type RequestType string

const (
	Publish RequestType = "PUB"
	Request RequestType = "REQ"
)

type Context struct {
	PutAppName string
	PutTime time.Time
	PutType RequestType
	Timeout time.Time
	TraceID string
}