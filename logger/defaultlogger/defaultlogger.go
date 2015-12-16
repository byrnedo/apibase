package defaultlogger

import (
	"fmt"
	"github.com/byrnedo/apibase/config"
	. "github.com/byrnedo/apibase/logger"
)

func init() {
	var (
		logLevel    LogLevel
		logFilePath string
		err         error
	)

	if logFilePath, err = config.Conf.GetString("log.file"); err != nil {
		fmt.Println("No log-file config var, logging to std out/err")
	}

	lvl := config.Conf.GetDefaultString("log.level", "info")
	switch lvl {
	case "trace":
		logLevel = TraceLevel
	case "info":
		logLevel = InfoLevel
	case "warn":
		logLevel = WarnLevel
	case "error":
		logLevel = ErrorLevel
	default:
		panic("Unknown log level given:" + lvl)
	}

	InitFileLog(logFilePath, logLevel)
}
