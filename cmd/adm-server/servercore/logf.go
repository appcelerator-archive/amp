package servercore

import (
	"log"
	"strings"
)

const (
	logError = 0
	logWarn  = 1
	logInfo  = 2
	logDebug = 3
)

// Logf logger
type Logf struct {
	level int
}

var logf = Logf{level: 3}

func (l *Logf) setLevel(level string) {
	if strings.ToLower(level) == "error" {
		l.level = logError
	} else if strings.ToLower(level) == "warn" {
		l.level = logWarn
	} else if strings.ToLower(level) == "info" {
		l.level = logInfo
	} else if strings.ToLower(level) == "debug" {
		l.level = logDebug
	}
}

func (l *Logf) levelString() string {
	switch l.level {
	case logError:
		return "error"
	case logWarn:
		return "warn"
	case logInfo:
		return "info"
	case logDebug:
		return "debug"
	default:
		return "?"
	}
}

func (l *Logf) error(format string, args ...interface{}) {
	if l.level >= logError {
		log.Printf(format, args...)
	}
}

func (l *Logf) warn(format string, args ...interface{}) {
	if l.level >= logWarn {
		log.Printf(format, args...)
	}
}

func (l *Logf) info(format string, args ...interface{}) {
	if l.level >= logInfo {
		log.Printf(format, args...)
	}
}

func (l *Logf) debug(format string, args ...interface{}) {
	if l.level >= logDebug {
		log.Printf(format, args...)
	}
}

func (l *Logf) printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
