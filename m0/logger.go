package main

import (
	"fmt"
)

type ILogger interface {
	Printf(level Level, format string, args ...interface{})
}

type KV struct {
	Key   string
	Value interface{}
}

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelPanic Level = "PANIC"
	LevelFatal Level = "FATAL"
)

type Level string

func NewLogger(kvs ...KV) ILogger {
	s := ""
	for _, v := range kvs {
		s += fmt.Sprintf(" %s=%s", v.Key, v.Value)
	}

	return &Logger{
		fields:      kvs,
		preparedBuf: []byte(s),
	}
}

type Logger struct {
	fields      []KV
	preparedBuf []byte
}

func (l *Logger) Printf(level Level, format string, args ...interface{}) {
	buf := make([]byte, 0, 8)

	buf = append(buf, []byte(level)...)
	buf = append(buf, []byte("\t")...)
	buf = append(buf, []byte(fmt.Sprintf(format, args...))...)
	buf = append(buf, []byte("\t")...)

	buf = append(buf, l.preparedBuf...)

	fmt.Println(string(buf))
}
