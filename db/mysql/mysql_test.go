package mysql

import (
	"reflect"
	"testing"
	"github.com/byrnedo/prefab"
	"time"
)

var (
	mysqlCtnr string
	mysqlUrl string
)

func TestDefaultConfig(t *testing.T) {
	conf := newDefaultConfig(func(c *Config) {

	})

	expectedConf := Config{
		ConnectString: "",
		MaxIdleCons:   4,
		MaxOpenCons:   16,
	}

	if !reflect.DeepEqual(conf, &expectedConf) {
		t.Errorf("Default conf not as expected\nexpected %+v\n     got %+v\n", &expectedConf, conf)
	}
}

func TestConnectAndQuery(t *testing.T) {

	mysqlCtnr, mysqlUrl = prefab.StartMysqlContainer()

	defer func() {
		if r := recover(); r != nil {
			t.Error("Should not have panicked: ", r)
		}
	}()

	prefab.WaitForMysql(mysqlUrl, 30 *time.Second)

	Init(func(c *Config) {
		c.ConnectString = mysqlUrl
		t.Log("["+c.ConnectString+"]")
	})

	if DB == nil {
		t.Error("DB was nil, not good.")
	}

	_ = DB.MustExec("select 1")

	prefab.Remove(mysqlCtnr)

}

