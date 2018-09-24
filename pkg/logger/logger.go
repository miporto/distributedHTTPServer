package logger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	logger  *log.Logger
	logFile *os.File
}

var instance *Logger
var once sync.Once

func newLogger(file string) {
	once.Do(func() {
		logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
		instance = &Logger{logger: logger, logFile: logFile}
	})
}

func NewLogger(file string) *Logger {
	newLogger(file)
	return instance
}

func GetInstance() *Logger {
	newLogger("log.txt")
	return instance
}

func (l *Logger) Info(msg string) {
	l.logger.Println("INFO: " + msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Println("ERROR: " + msg)
}

func (l *Logger) Close() {
	l.logFile.Close()
}
