package controllers
import (
	"github.com/byrnedo/apibase/routes"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	. "github.com/byrnedo/apibase/logger"
	"net/http"
)


type WebController interface {
	GetRoutes() []*routes.WebRoute
}

// Registers an array of route handlers to gorilla/mux
func RegisterRoutes(rtr *httprouter.Router, controller WebController) {
	for _, route := range controller.GetRoutes() {
		rtr.Handle(route.GetMethod(), route.GetPath(), route.GetHandler())
	}
}


// Controller with json helpers
type JsonController struct {}



// Serve standard 200
func (jC *JsonController) Serve(w http.ResponseWriter, data interface{}) {
	jC.ServeWithStatus(w, data, 200)
}

// Serve with custom status
func (jC *JsonController) ServeWithStatus(w http.ResponseWriter, data interface{}, status int) {
	bytes, err := json.Marshal(data)
	if err != nil {
		Error.Println("Failed to encode json:" + err.Error())
		panic("Failed to encode payload:" + err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(bytes)
}
