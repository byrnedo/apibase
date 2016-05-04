package mysql

import (
	"github.com/DavidHuie/gomigrate"
	_ "github.com/go-sql-driver/mysql"
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
	DB, err = sqlx.Open("mysql", conf.ConnectString)
	if err != nil {
		panic(err)
	}

	// set to reasonable values for production
	DB.SetMaxIdleConns(conf.MaxIdleCons)
	DB.SetMaxOpenConns(conf.MaxOpenCons)
}

func Migrate(migrations string) error {
	migrator, _ := gomigrate.NewMigrator(DB.DB, gomigrate.Mysql{}, migrations)
	return migrator.Migrate()
}
