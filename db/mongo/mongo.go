package mongo

import (
	"errors"
	"log"
	"reflect"
	"strings"

	. "github.com/byrnedo/apibase/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo session holder
var session *mgo.Session

// Get a connection from the session
func Conn() *mgo.Session {
	return session.Copy()
}

// Dial up to mongo using the "mongodb-url" from the app.conf
// First checks for environent variable GOAX_MONGODB_URL
func Init(url string, debugLog *log.Logger) {

	if debugLog != nil {
		mgo.SetDebug(true)
		mgo.SetLogger(debugLog)
	}

	Info.Println("Attempting to connect to [" + url + "]\n")

	sess, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	session = sess
	session.SetMode(mgo.Monotonic, true)
}

// Makes a bson map out of a list of fields
func mgoSelect(fields ...string) (selBson bson.M) {

	selBson = make(bson.M, len(fields))
	for _, s := range fields {
		selBson[s] = 1
	}
	return
}

// GetAll retrieves all records matches certain condition. Returns empty list if
// no records exist
func GetAll(c *mgo.Collection,
			query map[string]string, // TODO
			fields []string,
 			sortBy []string,
			offset int,
			limit int,
			result interface{}) (err error) {

	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("result argument must be a slice address")
	}

	moddedQuery := make(map[string]interface{})
	// change the ids to be object ids before we go
	if len(query) != 0 {
		for index, qVal := range query {
			// Chanage the json 'id' into mongo '_id'
			if index == "id" {
				index = "_id"
			}
			if strings.HasSuffix(index, "_id") && bson.IsObjectIdHex(qVal) {
				moddedQuery[index] = bson.ObjectIdHex(qVal)
			} else {
				moddedQuery[index] = qVal
			}

		}
	}

	err = c.Find(moddedQuery).Select(mgoSelect(fields...)).Skip(offset).Limit(limit).Sort(sortBy...).All(result)

	return
}
