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

//func Test_HandleFunc(t *testing.T) {
//	var natsOpts = NewNatsOptions()
//
//	var handleFunc1 = func(n *nats.Msg) {}
//	var handleFunc2 = func(n *nats.Msg) {}
//
//	natsOpts.HandleFunc("test.handle_func", handleFunc1)
//	natsOpts.HandleFunc("test.handle_func_2", handleFunc2)
//
//	expectedRoutes := []*Route{
//		&Route{
//			"test.handle_func",
//			handleFunc1,
//			nil,
//		},
//		&Route{
//			"test.handle_func_2",
//			handleFunc2,
//			nil,
//		},
//	}
//
//	routes := natsOpts.GetRoutes()
//	if len(routes) != 2 {
//		t.Error("Not 2 routes created")
//	}
//
//	for ind, route := range routes {
//		if route.Route != expectedRoutes[ind].Route {
//			t.Errorf("Routes not as expected:\nexpected %+v\nactual %+v", expectedRoutes[ind].Route, route.Route)
//		}
//		if route.Subsc != expectedRoutes[ind].Subsc {
//			t.Errorf("Subscr not as expected:\nexpected %+v\nactual %+v", expectedRoutes[ind].Subsc, route.Subsc)
//		}
//
//		f1 := reflect.ValueOf(route.Handler)
//		f2 := reflect.ValueOf(expectedRoutes[ind].Handler)
//
//		if f1.Pointer() != f2.Pointer() {
//			t.Errorf("Handlers not as expected:\nexpected %+v\nactual %+v", f2, f1)
//		}
//	}
//
//}
