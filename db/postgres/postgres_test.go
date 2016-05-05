package postgres

import (
	"reflect"
	"testing"
	"time"
	"github.com/byrnedo/prefab"
)


var (
	psqlCntr string
	psqlUrl string
)

func TestDefaultConfig(t *testing.T) {
	conf := newDefaultConfig(func(c *Config) {

	})

	expectedConf := Config{
		ConnectString:       "",
		MaxIdleCons:         4,
		MaxOpenCons:         16,
		EnableQueryInterp:   true,
		LogQueriesThreshold: 2 * time.Second,
		ProdMode:            true,
	}

	if !reflect.DeepEqual(conf, &expectedConf) {
		t.Errorf("Default conf not as expected\nexpected %+v\n     got %+v\n", &expectedConf, conf)
	}
}

func TestConnectAndQuery(t *testing.T) {
	var err error

	psqlCntr, psqlUrl := prefab.StartPostgresContainer()

	defer func() {
		if r := recover(); r != nil {
			t.Error("Should not have panicked: ", r)
		}
	}()

	prefab.WaitForPostgres(psqlUrl, 10*time.Second)
	Init(func(c *Config) {
		c.ConnectString = psqlUrl + "?sslmode=disable"
	})

	if DB == nil {
		t.Error("DB was nil, not good.")
	}

	_, err = DB.SQL("select $1", 1).Exec()
	if err != nil {
		t.Error("Failed to do select: ", err.Error())
	}

	prefab.Remove(psqlCntr)

}

