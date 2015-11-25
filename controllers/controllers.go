package controllers
import "github.com/byrnedo/apibase/routes"


type NatsController interface {
	GetRoutes() []*routes.NatsRoute
}