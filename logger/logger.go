package logger

import (
	"io"
	"log"
	"io/ioutil"
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
	InfoLevel LogLevel = 1
	WarnLevel LogLevel = 2
	ErrorLevel LogLevel = 3
)

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
