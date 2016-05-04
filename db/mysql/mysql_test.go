package mysql

import (
	"fmt"
	"github.com/byrnedo/apibase/dockertest"
	gDoc "github.com/fsouza/go-dockerclient"
	"reflect"
	"testing"
)

const (
	MysqlImage    = "mysql:5.7"
	MysqlDatabase = "test"
	MysqlPort     = "4306"
	MysqlPassword = "mysecretpassword"
)

var (
	dckrCli   *gDoc.Client
	mysqlCtnr *gDoc.Container
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

	fmt.Println("Setting up container")
	setupContainer(t)
	fmt.Println("Container set up.")

	defer func() {
		if r := recover(); r != nil {
			t.Error("Should not have panicked: ", r)
		}

	}()

	Init(func(c *Config) {
		c.ConnectString = fmt.Sprintf("root:%s@tcp(127.0.0.1:%s)/%s?timeout=90s", MysqlPassword, MysqlPort, MysqlDatabase)
		t.Log(c.ConnectString)
	})

	if DB == nil {
		t.Error("DB was nil, not good.")
	}

	_ = DB.MustExec("select 1")

}

func setupContainer(t *testing.T) {

	if id, err := dockertest.Running(MysqlImage); err != nil || len(id) < 1 {
		t.Log("Starting container")
		if _, err := dockertest.Start(MysqlImage, map[gDoc.Port][]gDoc.PortBinding{
			"3306/tcp": []gDoc.PortBinding{gDoc.PortBinding{
				HostIP:   "127.0.0.1",
				HostPort: MysqlPort,
			}},
		}, []string{"MYSQL_ROOT_PASSWORD=" + MysqlPassword, "MYSQL_DATABASE=" + MysqlDatabase}); err != nil {
			panic("Error starting mysql:" + err.Error())
		}

	}
}
