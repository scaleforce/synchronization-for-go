package stdpubsublog

import (
	"log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Infof(msg string, args ...any) {
	log.Printf("INFO: "+msg, args...)
}

func (logger *Logger) Warnf(msg string, args ...any) {
	log.Printf("WARN: "+msg, args...)
}

func (logger *Logger) Errorf(msg string, args ...any) {
	log.Printf("ERROR: "+msg, args...)
}

func (logger *Logger) Debugf(msg string, args ...any) {
	log.Printf("DEBUG: "+msg, args...)
}
