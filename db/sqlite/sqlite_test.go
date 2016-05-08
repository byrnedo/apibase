package sqlite

import (
	"reflect"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	conf := newDefaultConfig(func(c *Config) {
	})

	expectedConf := Config{
		ConnectString: ":memory:",
		MaxIdleCons:   4,
		MaxOpenCons:   16,
	}

	if !reflect.DeepEqual(conf, &expectedConf) {
		t.Errorf("Default conf not as expected\nexpected %+v\n     got %+v\n", &expectedConf, conf)
	}
}

func TestConnectAndQuery(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Error("Should not have panicked: ", r)
		}
	}()

	Init(func(c *Config) {
		c.ConnectString = ":memory:"
		t.Log("["+c.ConnectString+"]")
	})

	if DB == nil {
		t.Error("DB was nil, not good.")
	}

	_ = DB.MustExec("select 1")

}

