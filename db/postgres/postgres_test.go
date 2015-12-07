package postgres
import (
	"testing"
	"reflect"
	"time"
	gDoc "github.com/fsouza/go-dockerclient"
	"github.com/byrnedo/apibase/dockertest"
	"fmt"
)

const (
	PostgresImage = "postgres:9.4"
	PostgresPort = "5532"
	PostgresPassword = "mysecretpassword"

)

var (
	dckrCli *gDoc.Client
	psqlCntr *gDoc.Container
)


func TestDefaultConfig(t *testing.T) {
	conf := newDefaultConfig(func(c *Config) {

	})

	expectedConf := Config{
		ConnectString: "",
		MaxIdleCons: 4,
		MaxOpenCons: 16,
		EnableQueryInterp: true,
		LogQueriesThreshold:2 * time.Second,
		ProdMode: true,
	}

	if ! reflect.DeepEqual(conf, &expectedConf) {
		t.Errorf("Default conf not as expected\nexpected %+v\n     got %+v\n", &expectedConf, conf)
	}
}

func TestConnectAndQuery(t *testing.T) {
	var err error
	dckrCli, err = gDoc.NewClientFromEnv()
	if err != nil {
		panic("Failed to connect to docker:" + err.Error())
	}

	setupContainer()
	t.Log("Container set up.")

	defer func() {
		if r := recover(); r != nil {
			t.Error("Should not have panicked: ", r)
		}

	}()
	Init(func(c *Config) {
		c.ConnectString = fmt.Sprintf("dbname=postgres user=postgres password=%s host=localhost port=%d sslmode=disable connect_timeout=5", PostgresPassword, PostgresPort)
	})

	if DB == nil {
		t.Error("DB was nil, not good.")
	}

	_, err = DB.SQL("select $1", 1).Exec()
	if err != nil {
		t.Error("Failed to do select: ", err.Error())
	}


}

func setupContainer() {

	if id, err := dockertest.Running(PostgresImage); err != nil || len(id) < 1 {
		if _, err := dockertest.Start(PostgresImage, map[gDoc.Port][]gDoc.PortBinding{
			"5432/tcp" : []gDoc.PortBinding{gDoc.PortBinding{
				HostIP: "127.0.0.1",
				HostPort: PostgresPort,
			}},
		}); err != nil {
			panic("Error starting postgres:" + err.Error())
		}

	}
}


