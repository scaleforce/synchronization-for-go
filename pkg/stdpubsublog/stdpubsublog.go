package stdpubsublog

import (
	"log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Infof(msg string, args ...any) {
	log.Printf("INFO: "+msg+"\n", args...)
}

func (logger *Logger) Warnf(msg string, args ...any) {
	log.Printf("WARN: "+msg+"\n", args...)
}

func (logger *Logger) Errorf(msg string, args ...any) {
	log.Printf("ERROR: "+msg+"\n", args...)
}

func (logger *Logger) Debugf(msg string, args ...any) {
	log.Printf("DEBUG: "+msg+"\n", args...)
}
