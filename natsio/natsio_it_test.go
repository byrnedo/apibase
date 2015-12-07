package natsio
import (

	gDoc "github.com/fsouza/go-dockerclient"
	"os"
	"testing"
	"time"
"github.com/apcera/nats"
	"reflect"
)

const (
	NatsImage = "nats:latest"
	NatsPort = "4223"
	NatsLabel = "APIBASE_NATSIO_TEST"
)

var (
	dockCli *gDoc.Client
	natsContainer *gDoc.Container
)

type TestPayload struct {
	Data string
}

type TestRequest struct {
	NatsContext *NatsContext
	Error error
	Data *TestPayload
}

func (t *TestRequest) Context() *NatsContext {
	return t.NatsContext
}

func (t *TestRequest) NewContext(ctx *NatsContext) {
	t.NatsContext = ctx
}

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
				 "4222/tcp" : []gDoc.PortBinding{
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

	var natsCon *Nats


	time.Sleep(500 * time.Millisecond)

	natsCon, err := natsOpts.Connect()
	if err != nil {
		t.Error("Failed to connect:" + err.Error())
		return
	}

	var handler = func(subj string, reply string, testData *TestRequest) {
		t.Logf("Got message on nats: %+v", testData)
		//EncCon is nil at this point but that's ok
		//since it wont get called until after connecting
		//when it will then get a ping message.
		natsCon.Publish(reply, &TestRequest{
			Error: nil,
			Data: &TestPayload{"Pong"},
		})
	}

	natsCon.Subscribe("test.a", handler)

	response := TestRequest{}
	request := &TestRequest{
		Data: &TestPayload{"Ping"},
	}
	err = natsCon.Request("test.a", request, &response, 2 * time.Second)

	if err != nil {
		t.Error("Failed to get response:" + err.Error())
		return
	}

	natsCon.UnsubscribeAll()
	err = natsCon.EncCon.Request("test.a", request, &response, 2 * time.Second)
	if err == nil {
		t.Error("Should have failed to get response")
		return
	}
}

func TestHandleFunc(t *testing.T) {

	var handleFunc1 = func(n *nats.Msg) {}
	var handleFunc2 = func(n *nats.Msg) {}

	natsOpts := setupConnection()

	var natsCon *Nats

	time.Sleep(500 * time.Millisecond)

	natsCon, err := natsOpts.Connect()
	if err != nil {
		t.Error("Failed to connect:" + err.Error())
		return
	}

	natsCon.Subscribe("test.handle_func", handleFunc1)
	natsCon.Subscribe("test.handle_func_2", handleFunc2)

	natsCon.QueueSubscribe("test.handle_func", "group", handleFunc1)
	natsCon.QueueSubscribe("test.handle_func_2", "group_2", handleFunc2)

	expectedRoutes := []*Route{
		&Route{
			"test.handle_func",
			handleFunc1,
			nil,
			"",
		},
		&Route{
			"test.handle_func_2",
			handleFunc2,
			nil,
			"",
		},
		&Route{
			"test.handle_func",
			handleFunc1,
			nil,
			"group",
		},
		&Route{
			"test.handle_func_2",
			handleFunc2,
			nil,
			"group_2",
		},
	}

	routes := natsOpts.GetRoutes()
	if len(routes) != 4 {
		t.Error("Not 4 routes created")
	}

	for ind, route := range routes {
		if route.GetRoute() != expectedRoutes[ind].GetRoute() {
			t.Errorf("Routes not as expected:\nexpected %+v\nactual %+v", expectedRoutes[ind].GetRoute(), route.GetRoute())
		}
		if route.GetSubscription() == nil {
			t.Errorf("Subscr not as expected:\nexpected %+v\nactual %+v", expectedRoutes[ind].GetSubscription(), route.GetSubscription())
		}

		f1 := reflect.ValueOf(route.GetHandler())
		f2 := reflect.ValueOf(expectedRoutes[ind].GetHandler())

		if f1.Pointer() != f2.Pointer() {
			t.Errorf("Handlers not as expected:\nexpected %+v\nactual %+v", f2, f1)
		}
	}

}


func setupConnection() *NatsOptions{

	return NewNatsOptions(func(n *NatsOptions) error {
		n.Url = "nats://localhost:" + NatsPort
		n.Name = "it_test"
		n.SetEncoding(nats.JSON_ENCODER)
		n.Timeout = 10 * time.Second
		return nil
	})

}
