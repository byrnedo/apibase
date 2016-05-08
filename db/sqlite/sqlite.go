package sqlite

import (
	"github.com/DavidHuie/gomigrate"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
)

// global database (pooling provided by SQL driver)
var DB *sqlx.DB

// ConnectString should look something like "dbname=dat_test user=dat password=!test host=localhost sslmode=disable"
type Config struct {
	ConnectString string
	MaxIdleCons   int
	MaxOpenCons   int
}

func newDefaultConfig(confFuncs ...func(*Config)) *Config {
	config := Config{
		ConnectString: ":memory:",
		MaxIdleCons: 4,
		MaxOpenCons: 16,
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
	DB, err = sqlx.Open("sqlite3", conf.ConnectString)
	if err != nil {
		panic(err)
	}

	// set to reasonable values for production
	DB.SetMaxIdleConns(conf.MaxIdleCons)
	DB.SetMaxOpenConns(conf.MaxOpenCons)
}

func Migrate(migrations string) error {
	migrator, _ := gomigrate.NewMigrator(DB.DB, gomigrate.Sqlite3{}, migrations)
	return migrator.Migrate()
}
