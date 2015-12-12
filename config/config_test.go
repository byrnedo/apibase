package config

import (
	"github.com/byrnedo/typesafe-config/parse"
	"testing"
	"fmt"
)

func TestParse(t *testing.T) {
	var (
		tree *parse.Tree
		err  error
	)

	if tree, err = ParseFile("./test.conf"); err != nil {
		t.Error("Failed to read ./test.conf:" + err.Error())
	}

	if val, err := tree.GetConfig().GetInt("http.port"); err != nil || val != 1234 {
		t.Error("Incorrect Port value")
	}

	if val, err := tree.GetConfig().GetString("http.host"); err != nil || val != "abcdef" {
		t.Error("Incorrect Url value")
	}

	t.Logf("Conf object after: %v", tree.GetConfig())

}

type MyConfig struct {
	SectionA struct {
		Int int
		Uint uint
		Int32 int32
		Int64 int64
		Float32 float32
		Float64 float64
		String string
		SectionB  struct {
			Int int
			Uint uint
			Int32 int32
			Int64 int64
			Float32 float32
			Float64 float64
			String string `config:"strong-string"`
		}
	}
}

func TestPopulate(t *testing.T) {

	var (
		tree *parse.Tree
		err  error
	)

	if tree, err = ParseFile("./test.conf"); err != nil {
		t.Error("Failed to read ./test.conf:" + err.Error())
	}

	testStruct := &MyConfig{}

	Populate(testStruct, tree.GetConfig())

	t.Logf("After populate: %+v", testStruct)
	if fmt.Sprintf("%+v", testStruct) != "&{SectionA:{Int:-999 Uint:999 Int32:-999 Int64:999 Float32:999.999 Float64:999.999 String:lalala SectionB:{Int:-999 Uint:999 Int32:-999 Int64:999 Float32:999.999 Float64:999.999 String:lalala}}}" {
		t.Logf("Got: %+v", testStruct)
		t.Error("Not as expected.")
	}


}
