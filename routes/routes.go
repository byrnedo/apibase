package routes
import (
	"github.com/apcera/nats"
	"net/http"
)

type Route interface {
	GetPath() string
	GetHandler() interface{}
}
