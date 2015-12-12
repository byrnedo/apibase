package defaultconfig

import (
	"flag"
	"fmt"
	"github.com/byrnedo/apibase/config"
	"github.com/byrnedo/typesafe-config/parse"
	"os"
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

	tree, err := config.ParseFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config file:"+err.Error())
		os.Exit(1)
	}
	Conf = tree.GetConfig()
}
