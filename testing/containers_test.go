package testing

import "testing"

func TestTestEnvironment_Launch(t *testing.T) {
	tEnv := NewTestEnvironment().WithMongo().WithNats().WithMysql()
	tEnv.Launch()
	tEnv.Cleanup()
}

