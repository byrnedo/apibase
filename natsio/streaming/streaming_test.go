package streaming

import (
	"testing"
	"github.com/byrnedo/apibase/testenvironment"
	"os"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nats"
	"time"
	"fmt"
	"sync"
)

var (
	tEnv *testenvironment.TestEnvironment
	stanUrl string
	clientIdInc int
)

func setup() {

	tEnv = testenvironment.NewTestEnvironment().WithNatsStreaming()
	tEnv.Launch()

	stanUrl, _ = tEnv.GetUrl(testenvironment.NatsStreamingContainer)


}

func TestMain(m *testing.M) {
	setup()

	retCode := m.Run()

	// call with result of m.Run()
	tEnv.Cleanup()
	os.Exit(retCode)
}

func stanConnect(t *testing.T) stan.Conn {
	clientIdInc ++
	stanOpts := StanOptions{}
	stanOpts.Options = nats.DefaultOptions
	stanOpts.ClusterId = "test-cluster"
	stanOpts.ClientId = fmt.Sprintf("test-client-%d", clientIdInc)
	stanOpts.Options.Url = stanUrl


	var (
		//stanCon stan.Conn
		natsCon *nats.Conn
		err error
	)


	natsCon, err = stanOpts.Connect()
	if err != nil {
		t.Fatal(err)
	}


	stanCon, err := stan.Connect(stanOpts.ClusterId, stanOpts.ClientId, stan.NatsConn(natsCon), stan.ConnectWait(15*time.Second))
	if err != nil {
		panic("Failed to get stan con:" + err.Error())
	}

	return stanCon
}

func TestStanConnect(t *testing.T) {
	stanConnect(t)

}

func TestStan_BSubscribe(t *testing.T) {
	sCon :=  stanConnect(t)

	for i := 0; i < 10; i ++ {
		if err := sCon.Publish("test1", []byte("data")); err != nil {
			t.Fatal(err)
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(10)

	_, err := sCon.Subscribe("test1", func(m *stan.Msg) {
		t.Log("got one")
		wg.Done()
		if err := m.Ack(); err != nil {
			t.Fatal(err)
		}
	}, stan.DeliverAllAvailable(), stan.SetManualAckMode(), stan.AckWait(1 * time.Second))
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()


}

