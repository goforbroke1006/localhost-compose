package pkg

import (
	"fmt"
	"time"
)

type Logger interface {
	Info(serviceName string, message string)
	Infof(serviceName string, format string, args ...interface{})
	Errorf(serviceName string, format string, args ...interface{})
}

func NewLogger() *consoleLogger {
	return &consoleLogger{}
}

var _ Logger = &consoleLogger{}

type consoleLogger struct {
}

func (c consoleLogger) Info(serviceName string, message string) {
	fmt.Println(serviceName, "|", c.time(), "|", "INFO", "|", message)
}

func (c consoleLogger) Infof(serviceName string, format string, args ...interface{}) {
	fmt.Println(serviceName, "|", c.time(), "|", "INFO", "|", fmt.Sprintf(format, args...))
}

func (c consoleLogger) Errorf(serviceName string, format string, args ...interface{}) {
	fmt.Println(serviceName, "|", c.time(), "|", "ERROR", "|", fmt.Sprintf(format, args...))
}

func (c consoleLogger) time() string {
	return time.Now().Format(RFC3339)
}

const (
	RFC3339 = "2006-01-02 15:04:05"
)
