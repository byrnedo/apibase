package config

import (
	"errors"
	"github.com/byrnedo/typesafe-config/parse"
	"io/ioutil"
)

func ParseFile(path string) (*parse.Tree,error) {
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
