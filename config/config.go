package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/byrnedo/typesafe-config/parse"
	"io/ioutil"
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
	fs.BoolVar(&showUsage, "help", false, "Show usage information")

	if len(configPath) == 0 {
		configPath = "conf/app.conf"
	}

	if showUsage {
		flag.Usage()
		os.Exit(1)
	}

	tree, err := ParseFile(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config file:"+err.Error())
		os.Exit(1)
	}
	Conf = tree.GetConfig()
}

func ParseFile(path string) (*parse.Tree, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Failed to read config file")
	}
	tree, err := Parse(bytes)
	return tree, err
}

func Parse(configFileData []byte) (tree *parse.Tree, err error) {
	tree, err = parse.Parse("config", string(configFileData))
	return
}
