package config

import (
	"github.com/byrnedo/typesafe-config/parse"
	"testing"
)

func TestParse(t *testing.T) {
	var (
		tree *parse.Tree
		err  error
	)

	if tree, err = ParseFile("./test.conf"); err != nil {
		t.Error("Failed to read ../test.conf:" + err.Error())
	}

	if val, err := tree.GetConfig().GetInt("http.port"); err != nil || val != 1234 {
		t.Error("Incorrect Port value")
	}

	if val, err := tree.GetConfig().GetString("http.host"); err != nil || val != "abcdef" {
		t.Error("Incorrect Url value")
	}

	t.Logf("Conf object after: %v", tree.GetConfig())

}
