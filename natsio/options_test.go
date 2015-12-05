package natsio
import (
	"testing"
)

func Test_NewNatsOptionsAppliedInOrder(t *testing.T) {
	var natsOpts = NewNatsOptions(func(n *NatsOptions) error {
		n.Url = "test"
		n.MaxPingsOut = 4
		return nil
	}, func(n *NatsOptions) error {
		n.Url = "test2"
		n.MaxPingsOut = 5
		return nil
	})

	if natsOpts.Url != "test2" || natsOpts.MaxPingsOut != 5 {
		t.Error("Options not applied in correct order")
	}
}

func Test_NewNatsDefaultOptionsApplied(t *testing.T) {
	var natsOpts = NewNatsOptions()
	if natsOpts.NoRandomize != true {
		t.Error("Default options were not applied")
	}
}

