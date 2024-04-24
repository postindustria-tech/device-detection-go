package onpremise

import (
	"log"
	"os"
)

type LogWriter interface {
	Printf(format string, v ...interface{})
}

var DefaultLogger LogWriter = log.New(os.Stdout, "\r\n", log.LstdFlags)

// log wrapper controls provided logger
type logWrapper struct {
	enabled bool
	logger  LogWriter
}

func (l logWrapper) Printf(format string, v ...interface{}) {
	if !l.enabled {
		return
	}
	l.logger.Printf(format, v...)
}
