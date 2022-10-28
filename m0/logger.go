package main

type ILogger interface {
	Printf(format string, args ...interface{})
}

type Logger struct {
}

func (l *Logger) Printf(format string, args ...interface{}) {

}
