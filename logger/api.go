package logger

import (
	"io"
	"log"
)

type InfoLogger interface {
	Info(f string, m ...interface{})
}

type WarnLogger interface {
	Warn(f string, m ...interface{})
}

type ErrorLogger interface {
	Error(f string, m ...interface{})
}

type ServerLogger interface {
	Server(f string, m ...interface{})
}

type KillLogger interface {
	Kill(f string, m ...interface{})
}

type Logger interface {
	InfoLogger
	WarnLogger
	ErrorLogger
	ServerLogger
	KillLogger
}

type PrettyLogger struct {
	info,
	warning,
	error,
	server *log.Logger
}

func New(out io.Writer) *PrettyLogger {
	return &PrettyLogger{
		info:    log.New(out, "\033[0;36m[INFO]\033m", log.Ldate|log.Ltime),
		warning: log.New(out, "\033[0;33m[WARN]\033m", log.Ldate|log.Ltime),
		error:   log.New(out, "\033[0;31m[ERROR]\033m", log.Ldate|log.Ltime),
		server:  log.New(out, "\033[1;36m[SERVER]\033m", log.Ldate|log.Ltime),
	}
}

func (l *PrettyLogger) Info(f string, m ...interface{}) {
	l.info.Printf(f, m...)
}

func (l *PrettyLogger) Warn(f string, m ...interface{}) {
	l.warning.Printf(f, m...)
}

func (l *PrettyLogger) Error(f string, m ...interface{}) {
	l.error.Printf(f, m...)
}

func (l *PrettyLogger) Server(f string, m ...interface{}) {
	l.server.Printf(f, m...)
}

func (l *PrettyLogger) Kill(f string, m ...interface{}) {
	l.error.Fatalf(f, m...)
}
