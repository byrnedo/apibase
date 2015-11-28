package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gopkg.in/mgutz/dat.v1"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
	"time"
	"gopkg.in/mgutz/dat.v1/kvs"
)

// global database (pooling provided by SQL driver)
var DB *runner.DB

// If GetStore is nil, no caching is done
type KeyStore interface {
	GetStore() kvs.KeyValueStore
}

// ConnectString should look something like "dbname=dat_test user=dat password=!test host=localhost sslmode=disable"
type Config struct {
	ConnectString string
	Cache KeyStore
}

func Init(conf Config){

	if store := conf.Cache.GetStore(); store != nil {
		runner.SetCache(store)
	}

	// create a normal database connection through database/sql
	db, err := sql.Open("postgres", conf.ConnectString)
	if err != nil {
		panic(err)
	}

	// ensures the database can be pinged with an exponential backoff (15 min)
	runner.MustPing(db)

	// set to reasonable values for production
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(16)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = 10 * time.Millisecond

	DB = runner.NewDB(db, "postgres")
}