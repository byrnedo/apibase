package logger

import (
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

	logFormat int = log.Ldate | log.Ltime | log.LUTC | log.Llongfile
)

type LogLevel int

const (
	TraceLevel LogLevel = 0
	InfoLevel  LogLevel = 1
	WarnLevel  LogLevel = 2
	ErrorLevel LogLevel = 3
)

func init() {

	Trace = log.New(ioutil.Discard,
		"TRACE: ",
		logFormat)

	Info = log.New(os.Stdout,
		"INFO: ",
		logFormat)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		logFormat)

	Error = log.New(os.Stderr,
		"ERROR: ",
		logFormat)
}

func InitFileLog(logFilePath string, logLevel LogLevel) {
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

			Info.Println("Initializing file logger, path:" + logFilePath)
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
	Info.Println("Initializing logger")
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
