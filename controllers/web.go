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

// Registers an array of route handlers to gorilla/mux
func RegisterMuxRoutes(rtr *mux.Router, controller WebController) {
	for _, route := range controller.GetRoutes() {
		rtr.
		Methods(route.GetMethod()).
		Path(route.GetPath()).
		Name(route.GetName()).
		Handler(route.GetHandler())
	}
}

// Custom handler to transform error into
// json
// TODO use func type instead
type JsonErrorHandler interface {
	ToJson(message string, status int) []byte
}

// Controller with json helpers
type JsonController struct {
	errorHandler JsonErrorHandler
}

// Creates new controller using supplied error handler
func NewJsonController(errorHandler JsonErrorHandler) *JsonController {
	return &JsonController{
		errorHandler: errorHandler,
	}
}


// Serve standard 200
func (jC *JsonController) ServeJson(w http.ResponseWriter, data interface{}) {
	jC.ServeJsonStatus(w, data, 200)
}

// Serve with custom status
func (jC *JsonController) ServeJsonStatus(w http.ResponseWriter, data interface{}, status int) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		Error.Println("Failed to marshal payload:" + err.Error())
		jC.ServeError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(dataBytes)
}

// Serve error
func (jC *JsonController) ServeError(w http.ResponseWriter, message string, status int) {
	jsonErr := jC.errorHandler.ToJson(message, status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	w.Write(jsonErr)
}