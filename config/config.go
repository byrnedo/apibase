package config

import (
	"errors"
	"github.com/liyinhgqw/typesafe-config/parse"
	"io/ioutil"
	"regexp"
)

type Config struct {
	*parse.Tree
}

func (c *Config) ParseFile(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("Failed to read config file")
	}
	return c.Parse(bytes)
}

func (c *Config) Parse(configFileData []byte) (err error) {
	c.Tree, err = parse.Parse("config", string(configFileData))
	return
}

func unquoteString(value string) string {
	re := regexp.MustCompile("^\"(.*)\"$")
	if strippedVal := re.FindStringSubmatch(value); strippedVal != nil {
		return strippedVal[1]
	} else {
		return value
	}
}

func (c *Config) GetBool(key string) (bool, error) {
	typesafeConf := c.GetConfig()
	return typesafeConf.GetBool(key)
}

func (c *Config) GetDefaultBool(key string, defaultVal bool) bool {
	val, err := c.GetBool(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func (c *Config) GetFloat(key string) (float64, error) {
	typesafeConf := c.GetConfig()
	return typesafeConf.GetFloat(key)
}

func (c *Config) GetDefaultFloat(key string, defaultVal float64) float64 {
	val, err := c.GetFloat(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func (c *Config) GetInt(key string) (int, error) {
	typesafeConf := c.GetConfig()
	val64, err := typesafeConf.GetInt(key)
	return int(val64), err
}

func (c *Config) GetDefaultInt(key string, defaultVal int) int {
	val, err := c.GetInt(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func (c *Config) GetInt64(key string) (int64, error) {
	typesafeConf := c.GetConfig()
	return typesafeConf.GetInt(key)
}

func (c *Config) GetDefaultInt64(key string, defaultVal int64) int64 {
	val, err := c.GetInt64(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func (c *Config) GetString(key string) (val string, err error) {
	typesafeConf := c.GetConfig()
	if val, err = typesafeConf.GetString(key); err != nil {
		return
	}
	return unquoteString(val), nil
}

func (c *Config) GetDefaultString(key string, defaultVal string) string {
	val, err := c.GetString(key)
	if err != nil {
		return defaultVal
	}
	return val
}
