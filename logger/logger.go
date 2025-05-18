package logger

import (
	"fmt"
	"log"
)

var logLevelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

func NewLogger(prefix string, level LogLevel) *Logger {
	log.SetPrefix("[" + prefix + "] ")
	log.SetFlags(log.Ldate | log.Ltime)

	return &Logger{
		Level: level,
	}
}

func (l *Logger) log(level LogLevel, message string) {
	if level >= l.Level {
		log.Printf("[%s] %s", logLevelNames[level], message)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log(DEBUG, message)
}

func (l *Logger) Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log(INFO, message)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log(WARN, message)
}

func (l *Logger) Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log(ERROR, message)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log(FATAL, message)
	log.Fatal(message)
}
