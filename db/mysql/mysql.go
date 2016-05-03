package mysql

import (
	"github.com/DavidHuie/gomigrate"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/mgutz/dat.v1"
	"gopkg.in/mgutz/dat.v1/kvs"
	"time"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
	"github.com/jmoiron/sqlx"
)

// global database (pooling provided by SQL driver)
var DB *sqlx.DB

// If GetStore is nil, no caching is done
type KeyStore interface {
	GetStore() kvs.KeyValueStore
}

// ConnectString should look something like "dbname=dat_test user=dat password=!test host=localhost sslmode=disable"
type Config struct {
	ConnectString       string
	Cache               KeyStore
	MaxIdleCons         int
	MaxOpenCons         int
	EnableQueryInterp   bool
	ProdMode            bool
	LogQueriesThreshold time.Duration
}

func newDefaultConfig(confFuncs ...func(*Config)) *Config {
	config := Config{
		MaxIdleCons:         4,
		MaxOpenCons:         16,
		EnableQueryInterp:   true,
		LogQueriesThreshold: 2 * time.Second,
		ProdMode:            true,
	}
	for _, f := range confFuncs {
		f(&config)
	}
	return &config
}

func Init(confFunc func(conf *Config)) {

	var err error

	conf := newDefaultConfig(confFunc)

	// create a normal database connection through database/sql
	DB, err = sqlx.Open("mysql", conf.ConnectString)
	if err != nil {
		panic(err)
	}

	if conf.Cache != nil {
		if store := conf.Cache.GetStore(); store != nil {
			runner.SetCache(store)
		}
	}

	// set to reasonable values for production
	DB.SetMaxIdleConns(conf.MaxIdleCons)
	DB.SetMaxOpenConns(conf.MaxOpenCons)

	// set this to enable interpolation
	dat.EnableInterpolation = conf.EnableQueryInterp

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = !conf.ProdMode

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = conf.LogQueriesThreshold
}

func Migrate(migrations string) error {
	migrator, _ := gomigrate.NewMigrator(DB.DB, gomigrate.Mysql{}, migrations)
	return migrator.Migrate()
}
