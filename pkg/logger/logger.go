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

func GetInstance() *Logger {
	file := os.Getenv("LOGFILE")
	once.Do(func() {
		logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
		instance = &Logger{logger: logger, logFile: logFile}
	})
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
