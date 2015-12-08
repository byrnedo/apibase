package controllers
import (
	"github.com/byrnedo/apibase/routes"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	. "github.com/byrnedo/apibase/logger"
)


type WebController interface {
	GetRoutes() []*routes.WebRoute
}

func RegisterMuxRoutes(rtr *mux.Router, controller WebController) {
	for _, route := range controller.GetRoutes() {
		rtr.
		Methods(route.GetMethod()).
		Path(route.GetPath()).
		Name(route.GetName()).
		Handler(route.GetHandler())
	}
}

type JsonErrorHandler interface {
	ToJson(message string, status int) []byte
}

type JsonController struct {
	errorHandler JsonErrorHandler
}

func NewJsonController(errorHandler JsonErrorHandler) *JsonController {
	return &JsonController{
		errorHandler: errorHandler,
	}
}

func (jC *JsonController) ServeJson(w http.ResponseWriter, data interface{}) {
	jC.ServeJsonStatus(w, data, 200)
}

func (jC *JsonController) ServeJsonStatus(w http.ResponseWriter, data interface{}, status int) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		Error.Println("Failed to marshal payload:" + err.Error())
		jC.JsonError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(dataBytes)
}

func (jC *JsonController) JsonError(w http.ResponseWriter, message string, status int) {
	jsonErr := jC.errorHandler.ToJson(message, status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	w.Write(jsonErr)
}