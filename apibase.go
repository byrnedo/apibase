package apibase

import (
	"flag"
	"fmt"
	"github.com/byrnedo/apibase/config"
	. "github.com/byrnedo/apibase/logger"
	"os"
	"io"
)

var (
	Conf config.Config
)

func createLogger(logFilePath string, logLevel LogLevel ) {
	var (
		logWriter io.Writer
		errLogWriter io.Writer
	)

	InitLog(func(logOpts *LogOptions){

		logOpts.Level = logLevel
		if len(logFilePath) > 0 {
			file, err := os.OpenFile(logFilePath, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
			if err != nil {
				panic("Failed to open log file " + logFilePath + ":" + err.Error())
			}
			logWriter = file
			errLogWriter = file
			defer file.Close()
		} else {
			logWriter = os.Stdout
			errLogWriter = os.Stderr
		}


		logOpts.TraceHandle = logWriter
		logOpts.InfoHandle = logWriter
		logOpts.WarningHandle = logWriter
		logOpts.ErrorHandle = errLogWriter

	})
}

func Init() {

	var (
		configFile  string
		logFilePath string
		logLevel LogLevel
		showUsage   bool
	)

	flag.StringVar(&configFile, "conf", "", "Configuration file path")
	flag.BoolVar(&showUsage, "help", false, "Show usage information")
	flag.Parse()

	if len(configFile) == 0 {
		showUsage = true
	}

	if showUsage {
		flag.Usage()
		os.Exit(1)
	}

	Conf = config.Config{}
	err := Conf.ParseFile(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config file:"+err.Error())
		os.Exit(1)
	}

	if logFilePath, err = Conf.GetString("log.file"); err != nil {
		fmt.Println("No log-file config var, logging to std out/err")
	}

	switch Conf.GetDefaultString("log.level", "info") {
	case "trace":
		logLevel = TraceLevel
	case "info":
		logLevel = InfoLevel
	case "warn":
		logLevel = WarnLevel
	case "error":
		logLevel = ErrorLevel
	}

	createLogger(logFilePath, logLevel)
}
