package logger

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	Level LogLevel
}
