package testenvironment

import "testing"

func TestTestEnvironment_Launch(t *testing.T) {
	tEnv := NewTestEnvironment().WithMongo().WithNats().WithMysql().WithSftp()
	tEnv.Launch()
	tEnv.Cleanup()
}

