package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var consoleLogger *logrus.Logger

// InitConsoleLogger sets up a console logger
func InitConsoleLogger() {
	consoleLogger = logrus.New()
	consoleLogger.Out = os.Stdout
	consoleLogger.SetFormatter(&logrus.JSONFormatter{})
}

// Logger returns the default logger
func Logger() *logrus.Logger {
	return consoleLogger
}
