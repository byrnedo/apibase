package postgres
import (
	"testing"
	"reflect"
	"time"
	gDoc "github.com/fsouza/go-dockerclient"
	"os"
)

const (
	PostgresImage = "postgres"
	PostgresPort = "5442"
	PostgresLabel = "APIBASE_POSTGRES_TEST"
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

	psqlCntr = setupContainer(dckrCli)

	defer func() {
		if r := recover(); r != nil {
			t.Error("Should not have panicked: ", r)
		}

	}()
	Init(func(c *Config) {
		c.ConnectString = "dbname=postgres user=postgres password=" + PostgresPassword + " host=localhost port=" + PostgresPort + " sslmode=disable"
	})

	if DB == nil {
		t.Error("DB was nil, not good.")
	}


}

func setupContainer(dockCli *gDoc.Client) *gDoc.Container {
	var con *gDoc.Container

	if con = runningPostgresContainer(dockCli); con == nil {
		con = startPsqlContainer(dockCli)
	}

	return con
}

func startPsqlContainer(dockCli *gDoc.Client) *gDoc.Container {

	if err := dockCli.PullImage(gDoc.PullImageOptions{Repository: PostgresImage, OutputStream: os.Stdout}, gDoc.AuthConfiguration{}); err != nil {
		panic("Failed to pull postgres image:" + err.Error())
	}

	con, err := dockCli.CreateContainer(gDoc.CreateContainerOptions{
		Config: &gDoc.Config{
			Labels: map[string]string{
				PostgresLabel : "true",
			},
			Image: PostgresImage,
			Env: []string{"POSTGRES_PASSWORD=" + PostgresPassword},
		},
		HostConfig: &gDoc.HostConfig{
			PortBindings: map[gDoc.Port][]gDoc.PortBinding{
				"5432/tcp" : []gDoc.PortBinding{
					gDoc.PortBinding{HostIP: "127.0.0.1", HostPort: PostgresPort},
				},
			},
		},
	})
	if err != nil {
		panic("Failed to create postgres container:" + err.Error())
	}

	if err := dockCli.StartContainer(con.ID, nil); err != nil {
		panic("Failed to start postgres container:" + err.Error())
	}
	return con
}

func runningPostgresContainer(dockCli *gDoc.Client) *gDoc.Container {
	cons, err := dockCli.ListContainers(gDoc.ListContainersOptions{
		Filters: map[string][]string{
			"label":[]string{PostgresLabel},
		},
	})
	if err != nil {
		panic("Error getting container:" + err.Error())
	}

	if len(cons) == 0 {
		return nil
	}
	return &gDoc.Container{
		ID : cons[0].ID,
	}
}
