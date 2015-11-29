package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gopkg.in/mgutz/dat.v1"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
	"time"
	"github.com/DavidHuie/gomigrate"
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
		MaxIdleCons: 4,
		MaxOpenCons: 16,
		EnableQueryInterp: true,
		LogQueriesThreshold: 2 * time.Second,
		ProdMode: true,
	}
	for _, f := range confFuncs {
		f(&config)
	}
	return &config
}

func Init(confFunc func(conf *Config)) {

	conf := newDefaultConfig(confFunc)

	if conf.Cache != nil {
		if store := conf.Cache.GetStore(); store != nil {
			runner.SetCache(store)
		}
	}

	// create a normal database connection through database/sql
	db, err := sql.Open("postgres", conf.ConnectString)
	if err != nil {
		panic(err)
	}

	// ensures the database can be pinged with an exponential backoff (15 min)
	runner.MustPing(db)

	// set to reasonable values for production
	db.SetMaxIdleConns(conf.MaxIdleCons)
	db.SetMaxOpenConns(conf.MaxOpenCons)

	// set this to enable interpolation
	dat.EnableInterpolation = conf.EnableQueryInterp

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = !conf.ProdMode

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = conf.LogQueriesThreshold

	DB = runner.NewDB(db, "postgres")
}

func Migrate(migrations string) error {
	migrator, _ := gomigrate.NewMigrator(DB.DB.DB, gomigrate.Postgres{}, migrations)
	return migrator.Migrate()
}