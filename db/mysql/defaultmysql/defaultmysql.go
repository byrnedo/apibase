// This package provides a default mysql init.
// Configuration is loaded from config file under 'mysql' section
//
package defaultmysql

import (
	"github.com/byrnedo/apibase/config"
	"github.com/byrnedo/apibase/db/mysql"
	. "github.com/byrnedo/apibase/logger"
	"github.com/byrnedo/typesafe-config/parse"
	mysqlReal "github.com/go-sql-driver/mysql"
)

func init() {

	mysql.Init(func(c *mysql.Config) {
		parse.Populate(c, config.Conf, "mysql")

		c2, err := mysqlReal.ParseDSN(c.ConnectString)
		if err == nil {
			Info.Printf("Attempting to connect to %s@%s\n", c2.User, c2.Addr)
		}
	})
}
