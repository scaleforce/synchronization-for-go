package stdpubsublog

import (
	"log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Info(msg string, args ...any) {
	log.Printf("INFO: "+msg+"\n", args...)
}

func (logger *Logger) Warn(msg string, args ...any) {
	log.Printf("WARN: "+msg+"\n", args...)
}

func (logger *Logger) Error(msg string, args ...any) {
	log.Printf("ERROR: "+msg+"\n", args...)
}

func (logger *Logger) Debug(msg string, args ...any) {
	log.Printf("DEBUG: "+msg+"\n", args...)
}
