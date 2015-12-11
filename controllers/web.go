package controllers

import (

	"encoding/json"
	"errors"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/apibase/routes"
	"github.com/byrnedo/mapcast"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
	"regexp"
)

var (
	acceptsHtmlRegex = regexp.MustCompile(`(text/html|application/xhtml\+xml)(?:,|$)`)
	acceptsXmlRegex  = regexp.MustCompile(`(application/xml|text/xml)(?:,|$)`)
	acceptsJsonRegex = regexp.MustCompile(`(application/json)(?:,|$)`)
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
type JsonController struct {
	BaseController
}

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

type BaseController struct{}

// Gets first query string value as int
func (bC *BaseController) QueryInt(r *http.Request, param string) (int, error) {
	if vals := r.URL.Query()[param]; len(vals) > 0 {
		return strconv.Atoi(vals[0])
	}
	return 0, errors.New("Not found")
}

// Makes a map from a query string parameter formatted so that
// a query string like ?query=field:val&query=field:val2 becomes
// map[string]string{"field": {
// 		"val",
//		"val2",
// }
func (bC *BaseController) QueryMap(r *http.Request, param string) map[string]string {
	mapData := make(map[string]string, 0)
	if vals := r.URL.Query()[param]; len(vals) > 0 {
		for _, unsplitKeyVal := range vals {
			keyVal := strings.SplitN(unsplitKeyVal, ":", 2)
			if len(keyVal) == 2 {
				mapData[keyVal[0]] = keyVal[1]
			} else {
				mapData[keyVal[0]] = ""
			}
		}
	}
	return mapData
}

// Uses byrnedo/mapcast to turn query string map into a typed map.
func (bC *BaseController) QueryInterfaceMap(r *http.Request, param string, target interface{}) (ifMap map[string]interface{}) {
	stringMap := bC.QueryMap(r, param)
	return mapcast.CastViaJsonToBson(stringMap, target)
}

func (bC *BaseController) AcceptsJson(r *http.Request) bool {
	return acceptsJsonRegex.MatchString(r.Header.Get("Accept"))
}

func (bC *BaseController) AcceptsHtml(r *http.Request) bool {
	return acceptsHtmlRegex.MatchString(r.Header.Get("Accept"))
}

func (bC *BaseController) AcceptsXml(r *http.Request) bool {
	return acceptsXmlRegex.MatchString(r.Header.Get("Accept"))
}
