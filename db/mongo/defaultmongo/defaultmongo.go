// This package provides a default mongo init.
// Configuration is loaded from config file under 'mongo' section
//
// See MongoConf docs for example config
package defaultmongo

import (
	"github.com/byrnedo/apibase/config"
	"github.com/byrnedo/apibase/db/mongo"
	. "github.com/byrnedo/apibase/logger"
	"gopkg.in/mgo.v2"
	"github.com/byrnedo/typesafe-config/parse"
)

// Mongo session holder
var session *mgo.Session

// Get a connection from the session
func Conn() *mgo.Session {
	return session.Copy()
}

func init() {

	mConf := &mongo.MongoConf{}
	parse.Populate(mConf, config.Conf, "mongo")
	Info.Println("Mongo config:", mConf)

	Info.Println("Attempting to connect to [" + mConf.Url + "]")

	session = mongo.InitFromConfig(mConf, Trace)
}
