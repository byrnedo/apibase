package streaming

import (
	. "github.com/byrnedo/apibase/natsio/protobuf"
	"testing"
	"github.com/byrnedo/apibase/testenvironment"
	"os"
	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/nats"
	"time"
	pth "github.com/byrnedo/apibase/helpers/pointerhelp"
)

var (
	tEnv *testenvironment.TestEnvironment
	stanUrl string
)

type Wrap struct {
	*TestMessage
}

func (w *Wrap) SetContext(ctx *NatsContext) {
	w.Context = ctx
}

func WrapMessage(msg *TestMessage) *Wrap {
	return &Wrap{msg}
}

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

func stanConnect(t *testing.T) *Stan {
	stanOpts := NewStanOptions(func(o *StanOptions)error{
		o.Options = nats.DefaultOptions
		o.ClusterId = "test-cluster"
		o.ClientId = "test-client"
		o.Options.Url = stanUrl
		return nil
	})


	stanCon, err := stanOpts.Connect()
	if err != nil {
		t.Fatal(err)
	}

	return stanCon
}

func TestStanConnect(t *testing.T) {
	sCon := stanConnect(t)
	sCon.Con.Close()

}

func TestStan_BSubscribe(t *testing.T) {
	sCon :=  stanConnect(t)

	defer sCon.Con.Close()


	for i := 0; i < 10; i ++ {
		if err := sCon.Con.Publish("test1", []byte("data")); err != nil {
			t.Fatal(err)
		}
	}

	handled := make(chan bool, 1)

	err := sCon.Subscribe("test1", func(m *stan.Msg) {
		t.Log("got one")
		handled <- true
	}, stan.DurableName("streaming-test"))
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i ++ {
		if err := sCon.Con.Publish("test1", []byte("data")); err != nil {
			t.Fatal(err)
		}
	}


	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	handledNum := 0
	WAIT_LOOP:
	for {
		select {
		case <-handled:
			handledNum ++
		case <-timeout:
			if handledNum != 10 {
				t.Fatalf("Unexpected number of messages %d", handledNum)
			} else {
				break WAIT_LOOP
			}
		}
	}


}


func TestStan_CSubscribeWithProto(t *testing.T) {
	sCon :=  stanConnect(t)
	defer sCon.Con.Close()

	testMsg := WrapMessage(&TestMessage{Data:pth.StringPtr("test")})

	for i := 0; i < 10; i ++ {
		if err := sCon.Publish("test1", testMsg); err != nil {
			t.Fatal(err)
		}
	}

	handled := make(chan bool, 1)

	err := sCon.Subscribe("test1", func(m *stan.Msg) {
		t.Log("got one")
		handled <- true
	}, stan.DurableName("streaming-test"))
	if err != nil {
		t.Fatal(err)
	}

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()

	handledNum := 0
	WAIT_LOOP:
	for {
		select {
		case <-handled:
			handledNum ++
		case <-timeout:
			if handledNum != 10 {
				t.Fatalf("Unexpected number of messages %d", handledNum)
			} else {
				break WAIT_LOOP
			}
		}
	}

}
