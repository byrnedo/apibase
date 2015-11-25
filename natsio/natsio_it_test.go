package natsio
import (

	gDoc "github.com/fsouza/go-dockerclient"
	"os"
	"testing"
	"time"
"github.com/apcera/nats"
)

type TestData struct {
	Message string
}

const (
	NatsImage = "nats"
	NatsPort = "4222"
	NatsLabel = "APIBASE_NATSIO_TEST"
)

var (
	dockCli *gDoc.Client
	natsContainer *gDoc.Container
)
func startNatsContainer(dockCli *gDoc.Client) *gDoc.Container {

	if err := dockCli.PullImage(gDoc.PullImageOptions{Repository: NatsImage, OutputStream: os.Stdout}, gDoc.AuthConfiguration{}); err != nil {
		panic("Failed to pull nats image:" + err.Error())
	}

	con, err := dockCli.CreateContainer(gDoc.CreateContainerOptions{
		Config: &gDoc.Config{
			Cmd : []string{"--debug", "--logtime"},
			Labels: map[string]string{
				NatsLabel : "true",
			},
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

func runningNatsContainer(dockCli *gDoc.Client) *gDoc.Container {
	cons, err :=dockCli.ListContainers(gDoc.ListContainersOptions{
		Filters: map[string][]string{
			"label":[]string{NatsLabel},
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

func setup(dockCli *gDoc.Client) *gDoc.Container {
	var con *gDoc.Container

	if con = runningNatsContainer(dockCli); con ==  nil{
		con = startNatsContainer(dockCli)
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

	// call with result of m.Run()
	os.Exit(retCode)
}

func TestNewNatsConnect(t *testing.T) {
	natsOpts := setupConnection()

	var handler = func(subj string, reply string, testData *TestData) {
		t.Logf("Got message on nats: %+v", testData)
		//EncCon is nil at this point but that's ok
		//since it wont get called until after connecting
		//when it will then get a ping message.
		natsOpts.EncCon.Publish(reply, &TestData{Message: "Pong"})
	}

	natsOpts.HandleFunc("test.a", handler)

	err := natsOpts.ListenAndServeOrRetry(3)
	if err != nil {
		t.Error("Failed to connect:" + err.Error())
		return
	}

	response := TestData{}
	err = natsOpts.EncCon.Request("test.a", TestData{"Ping"}, &response, 2 * time.Second)
	if err != nil {
		t.Error("Failed to get response:" + err.Error())
		return
	}

	natsOpts.UnsubscribeAll()
	err = natsOpts.EncCon.Request("test.a", TestData{"Ping"}, &response, 2 * time.Second)
	if err == nil {
		t.Error("Should have failed to get response")
		return
	}
}


func setupConnection() *Nats{

	return NewNats(func(n *nats.Options) error {
		n.Url = "nats://localhost:" + NatsPort
		return nil
	})

}
