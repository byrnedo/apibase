package apibase

import (
	"flag"
	"fmt"
	"github.com/byrnedo/apibase/config"
	. "github.com/byrnedo/apibase/logger"
	"io/ioutil"
	"os"
)

var (
	Conf config.Config
	configFile  string
	logFilePath string
	showUsage   bool
)

func createLogger(logFilePath string) {

	var (
		logOpts *LogOptions
	)

	if len(logFilePath) > 0 {
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("Failed to open log file " + logFilePath + ":" + err.Error())
		}
		defer file.Close()
		logOpts = &LogOptions{
			TraceHandle:   ioutil.Discard,
			InfoHandle:    file,
			WarningHandle: file,
			ErrorHandle:   file,
		}
	} else {
		logOpts = &LogOptions{
			TraceHandle:   ioutil.Discard,
			InfoHandle:    os.Stdout,
			WarningHandle: os.Stdout,
			ErrorHandle:   os.Stderr,
		}
	}
	InitLog(*logOpts)

}

func Init() {
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

	if logFilePath, err = Conf.GetString("log-file"); err != nil {
		fmt.Println("No log-file config var, logging to std out/err")
	}

	createLogger(logFilePath)
}