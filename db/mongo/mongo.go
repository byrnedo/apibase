package mongo

import (
	"log"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Structure for mongo config section
// would look like this in the conf file:
//
//mongo {
//    url = "mongodb://localhost:27017"
//    connect-timeout = 10
//    debug = false
//}
type MongoConf struct {
	Url            string `config:"url"`
	ConnectTimeout int    `config:"connect-timeout,10"`
	Debug          bool
}

// Dial up to mongo using the "mongodb-url" from the app.conf
// First checks for environent variable GOAX_MONGODB_URL
func InitFromConfig(conf *MongoConf, debugLog *log.Logger) *mgo.Session {

	if debugLog != nil {
		mgo.SetDebug(conf.Debug)
		mgo.SetLogger(debugLog)
	}

	sess, err := mgo.DialWithTimeout(conf.Url, time.Duration(conf.ConnectTimeout)*time.Second)
	if err != nil {
		panic(err)
	}

	sess.SetMode(mgo.Monotonic, true)
	return sess
}

// Dial up to mongo using the "mongodb-url" from the app.conf
// First checks for environent variable GOAX_MONGODB_URL
func Init(url string, debugLog *log.Logger) *mgo.Session {

	if debugLog != nil {
		mgo.SetDebug(true)
		mgo.SetLogger(debugLog)
	}

	sess, err := mgo.DialWithTimeout(url, 15*time.Second)
	if err != nil {
		panic(err)
	}

	sess.SetMode(mgo.Monotonic, true)
	return sess
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
			switch query[index].(type) {
			case string:
				var strVal = qVal.(string)
				if bson.IsObjectIdHex(strVal) {
					query[index] = bson.ObjectIdHex(strVal)
				}
			}
		}
	}
}
