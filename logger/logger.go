package logger

import (
	"fmt"
	"github.com/byrnedo/apibase/config"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger

	baseOptions *LogOptions

	logFormat int = log.Ldate | log.Ltime | log.Lshortfile
)

type LogLevel int

const (
	TraceLevel LogLevel = 0
	InfoLevel  LogLevel = 1
	WarnLevel  LogLevel = 2
	ErrorLevel LogLevel = 3
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

	createLogger(logFilePath, logLevel)
}

func createLogger(logFilePath string, logLevel LogLevel) {
	var (
		logWriter    io.Writer
		errLogWriter io.Writer
	)

	InitLog(func(logOpts *LogOptions) {

		logOpts.Level = logLevel
		if len(logFilePath) > 0 {
			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

func GetLogOptions() LogOptions {
	return *baseOptions
}

type LogOptions struct {
	Level         LogLevel
	TraceHandle   io.Writer
	TraceFormat   int
	InfoHandle    io.Writer
	InfoFormat    int
	WarningHandle io.Writer
	WarningFormat int
	ErrorHandle   io.Writer
	ErrorFormat   int
}

func (options *LogOptions) seedDefaults() {
	options.Level = InfoLevel
	options.TraceFormat = logFormat
	options.InfoFormat = logFormat
	options.InfoHandle = os.Stdout
	options.WarningFormat = logFormat
	options.WarningHandle = os.Stdout
	options.ErrorFormat = logFormat
	options.ErrorHandle = os.Stderr
}

func InitLog(optFunc func(*LogOptions)) {
	var options = &LogOptions{}
	options.seedDefaults()
	optFunc(options)

	if options.Level > TraceLevel {
		options.TraceHandle = ioutil.Discard
	}

	if options.Level > InfoLevel {
		options.InfoHandle = ioutil.Discard
	}

	if options.Level > InfoLevel {
		options.WarningHandle = ioutil.Discard
	}

	Trace = log.New(options.TraceHandle,
		"TRACE: ",
		options.TraceFormat)

	Info = log.New(options.InfoHandle,
		"INFO: ",
		options.InfoFormat)

	Warning = log.New(options.WarningHandle,
		"WARNING: ",
		options.WarningFormat)

	Error = log.New(options.ErrorHandle,
		"ERROR: ",
		options.ErrorFormat)
}
