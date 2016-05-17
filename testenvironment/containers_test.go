package testenvironment

import "testing"

func TestTestEnvironment_Launch(t *testing.T) {
	tEnv := NewTestEnvironment().WithMongo().WithNats().WithMysql().WithMysqlTmpfs().WithMongoTmpfs()
	tEnv.Launch()
	tEnv.Cleanup()
}

