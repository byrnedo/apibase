package logger

import (
	"io"
	"log"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger

	logFormat int = log.Ldate | log.Ltime | log.Lshortfile
)

type LogOptions struct {
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
	if options.TraceFormat == 0 {
		options.TraceFormat = logFormat
	}

	if options.InfoFormat == 0 {
		options.InfoFormat = logFormat
	}

	if options.WarningFormat == 0 {
		options.WarningFormat = logFormat
	}

	if options.ErrorFormat == 0 {
		options.ErrorFormat = logFormat
	}
}

func InitLog(options LogOptions){
	options.seedDefaults()

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
