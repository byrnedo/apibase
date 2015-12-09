package controllers

import (
	"github.com/byrnedo/apibase/natsio"
	"github.com/byrnedo/apibase/routes"
)

type NatsController interface {
	GetRoutes() []*routes.NatsRoute
}

func SubscribeNatsRoutes(natsCon *natsio.Nats, queueName string, controllers NatsController) {
	for _, route := range controllers.GetRoutes() {
		natsCon.QueueSubscribe(route.GetPath(), queueName, route.GetHandler())
	}
}
