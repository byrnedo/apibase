package natsio
import (

	gDoc "github.com/fsouza/go-dockerclient"
	"os"
	"testing"
	"time"
)

type TestData struct {
	Message string
}

const (
	NatsImage = "nats"
	NatsPort = "4222"
)

var (
	dockCli *gDoc.Client
	natsContainer *gDoc.Container
)


func setup(dockCli *gDoc.Client) *gDoc.Container {
	if err := dockCli.PullImage(gDoc.PullImageOptions{Repository: NatsImage, OutputStream: os.Stdout}, gDoc.AuthConfiguration{}); err != nil {
		panic("Failed to pull nats image:" + err.Error())
	}

	con, err := dockCli.CreateContainer(gDoc.CreateContainerOptions{
		Config: &gDoc.Config{
			Image: NatsImage,
		},
		HostConfig: &gDoc.HostConfig{
			PortBindings: map[gDoc.Port][]gDoc.PortBinding{
				NatsPort + "/tcp" : []gDoc.PortBinding{
					gDoc.PortBinding{HostIP: "127.0.0.1", HostPort: NatsPort},
				},
			},
		},
	})
	if err != nil {
		panic("Failed to create nats container:" + err.Error())
	}

	if err := dockCli.StartContainer(con.ID, nil); err != nil {
		panic("Failed to start nats container:" + err.Error())
	}

	return con
}

func teardown() {
	err := dockCli.RemoveContainer(gDoc.RemoveContainerOptions{
		Force: true,
		ID: natsContainer.ID,
	})
	if err != nil {
		panic("Failed to remove nats container:" + err.Error())
	}
}

func TestMain(m *testing.M) {
	// your func

	var err error

	dockCli, err = gDoc.NewClientFromEnv()
	if err != nil {
		panic("Failed to connect to docker:" + err.Error())
	}

	natsContainer = setup(dockCli)

	retCode := m.Run()

	teardown()
	// call with result of m.Run()
	os.Exit(retCode)
}

func TestNewNatsConnect(t *testing.T) {
	natsOpts := NewNats(func(n *Nats) error {
		n.Url = "nats://localhost:" + NatsPort
		return nil
	})

	var handler = func(testData *TestData) {
		t.Logf("Got message on nats: %+v", testData)
	}

	natsOpts.HandleFunc("test.a", handler)

	err := natsOpts.ListenAndServe()
	if err != nil {
		teardown()
		t.Error("Failed to connect:" + err.Error())
	}

	response := TestData{}
	err = natsOpts.GetEncCon().Request("test.a", TestData{"Ping"}, response, 2 * time.Second)
	if err != nil {
		teardown()
		t.Error("Failed to get response:" + err.Error())
	}
}
