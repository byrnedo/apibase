package config

import (
	"testing"
)

func TestParse(t *testing.T) {
	conf := Config{}

	t.Logf("Conf object before: %v", conf)

	if err := conf.ParseFile("./test.conf"); err != nil {
		t.Error("Failed to read ../test.conf:" + err.Error())
	}

	if conf.Http.Port != 1234 {
		t.Error("Incorrect Port value")
	}

	if conf.Http.Host != "abcdef" {
		t.Error("Incorrect Url value")
	}

	t.Logf("Conf object after: %v", conf)

}
