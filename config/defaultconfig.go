// This package provides a default config init.
// Config file is loaded according to -conf flag.
// Falls back to conf/app.conf location.
//
package config

import (
	"flag"
	"github.com/byrnedo/typesafe-config/parse"
)

var (
	configPath string
	showUsage  bool
	Conf       *parse.Config
)

func init() {
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.StringVar(&configPath, "conf", "", "Configuration file path")

	if len(configPath) == 0 {
		configPath = "conf/app.conf"
	}

	tree, err := parse.ParseFile(configPath)
	if err != nil {

		panic("Error parsing config file: " + err.Error())
	}
	Conf = tree.GetConfig()
}
