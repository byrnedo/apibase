package controllers
import (
	"github.com/byrnedo/apibase/routes"
	"github.com/byrnedo/apibase/natsio"
"github.com/gorilla/mux"
)


type NatsController interface {
	GetRoutes() []*routes.NatsRoute
}

func SubscribeNatsRoutes(natsCon *natsio.Nats, queueName string, controllers NatsController) {
	for _, route := range controllers.GetRoutes() {
		natsCon.QueueSubscribe(route.GetPath(), queueName, route.GetHandler())
	}
}

type WebController interface {
	GetRoutes() []*routes.WebRoute
}

func RegisterMuxRoutes(rtr *mux.Router, controller WebController){
	for _, route := range controller.GetRoutes() {
		rtr.
		Methods(route.GetMethod()).
		Path(route.GetPath()).
		Name(route.GetName()).
		Handler(route.GetHandler())
	}
}
