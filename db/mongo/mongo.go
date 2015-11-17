package mongo

import (
	"errors"
	"log"
	"os"
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
func Init(url string, debug bool ) {
	envUrl := os.Getenv("MONGODB_URL")
	if len(envUrl) > 0 {
		Info.Println("Found env GOAX_MONGODB_URL, using that")
		url = envUrl
	}

	if debug == true {
		mgo.SetDebug(true)
		var aLogger *log.Logger
		aLogger = log.New(os.Stderr, "mgo: ", log.LstdFlags)
		mgo.SetLogger(aLogger)
	}

	Info.Println("Attempting to connect to [" + url + "]")

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
func GetAll(
table string,
query map[string]string, // TODO
fields []string,
sortby []string,
order []string,
offset int,
limit int,
result interface{},
) (err error) {

	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("result argument must be a slice address")
	}

	mConn := Conn()
	defer mConn.Close()
	c := mConn.DB("").C(table)

	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			//qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) > 0 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		}
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

	err = c.Find(moddedQuery).Select(mgoSelect(fields...)).Skip(offset).Limit(limit).Sort(sortFields...).All(result)

	return
}
