package mongo

import (
	"log"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
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

	sess, err := mgo.DialWithTimeout(url, 15*time.Second)
	if err != nil {
		panic(err)
	}

	session = sess
	session.SetMode(mgo.Monotonic, true)
}

// Makes a bson map out of a list of fields
func ToBsonMap(fields ...string) (selBson bson.M) {
	selBson = make(bson.M, len(fields))
	for _, s := range fields {
		selBson[s] = 1
	}
	return
}

func ConvertObjectIds(query map[string]interface{}) {
	for index, qVal := range query {
		// Chanage the json 'id' into mongo '_id'
		if index == "id" {
			index = "_id"
		}
		if strings.HasSuffix(index, "_id") {
			switch query[index].(type){
			case string:
				var strVal = qVal.(string)
				if bson.IsObjectIdHex(strVal) {
					query[index] = bson.ObjectIdHex(strVal)
				}
			}
		}
	}
}

