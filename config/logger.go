package config

import (
	"io"
	"log"
	"os"
	"strings"
)

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
	writer  io.Writer
}

func NewLogger(prefix string) *Logger {
	writer := io.Writer(os.Stdout)
	flags := log.Ldate | log.Ltime | log.Lmsgprefix
	prefix = strings.ToUpper(prefix)

	return &Logger{
		debug:   log.New(writer, prefix+" DEBUG: ", flags),
		info:    log.New(writer, prefix+" INFO: ", flags),
		warning: log.New(writer, prefix+" WARNING: ", flags),
		error:   log.New(writer, prefix+" ERROR: ", flags),
		writer:  writer,
	}
}

func (l *Logger) Debug(v ...any) {
	l.debug.Println(v...)
}

func (l *Logger) Info(v ...any) {
	l.info.Println(v...)
}

func (l *Logger) Warn(v ...any) {
	l.warning.Println(v...)
}

func (l *Logger) Error(v ...any) {
	l.error.Println(v...)
}

func (l *Logger) Debugf(format string, v ...any) {
	l.debug.Printf(format, v...)
}

func (l *Logger) Infof(format string, v ...any) {
	l.info.Printf(format, v...)
}

func (l *Logger) Warnf(format string, v ...any) {
	l.warning.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...any) {
	l.error.Printf(format, v...)
}