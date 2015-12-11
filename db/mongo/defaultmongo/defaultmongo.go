package defaultmongo

import (
	"github.com/byrnedo/apibase/config"
	"github.com/byrnedo/apibase/db/mongo"
	"github.com/byrnedo/apibase/helpers/env"
	. "github.com/byrnedo/apibase/logger"
	"gopkg.in/mgo.v2"
)

// Mongo session holder
var session *mgo.Session

// Get a connection from the session
func Conn() *mgo.Session {
	return session.Copy()
}

func init() {

	mongoUrl := env.GetOr("MONGO_URL", config.Conf.GetDefaultString("mongo.url", ""))
	Info.Println("Attempting to connect to [" + mongoUrl + "]")

	session = mongo.Init(mongoUrl, Trace)
}
