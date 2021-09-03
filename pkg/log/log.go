package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// logger NONE LEVEL
var logger *log.Logger

// 默认log级别ERROR
var loggerLevel int

const (
	loggerLevelDebug = iota
	loggerLevelInfo
	loggerLevelError
)

//SetLoggerConfig   ["INFO","DEBUG","ERROR"]
func SetLoggerConfig(level string, logfile string) {
	if logfile != "" {
		logFile, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		writers := []io.Writer{
			logFile,
			os.Stdout,
		}
		fileStdoutWriter := io.MultiWriter(writers...)
		logger.SetOutput(fileStdoutWriter)
	}
	switch strings.ToUpper(level) {
	case "INFO":
		loggerLevel = loggerLevelInfo
	case "DEBUG":
		loggerLevel = loggerLevelDebug
	case "ERROR":
		loggerLevel = loggerLevelError
	default:
		// default Logger
		loggerLevel = loggerLevelError
	}

}

//init  setDefalutLogger sets the logger for this package
func init() {
	logger = log.New(os.Stdout, "[ERROR ] ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}

// Print delegates to the Logger
func Print(v ...interface{}) {
	logger.SetPrefix("[DEFAULT] ")
	logger.Output(2, fmt.Sprint(v...))
}

// Printf delegates to the Logger
func Printf(format string, v ...interface{}) {
	logger.SetPrefix("[DEFAULT] ")
	logger.Output(2, fmt.Sprintf(format, v...))
}

// Fatal log.Fatal
func Fatal(v ...interface{}) {
	logger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Debug delegates to the Logger
func Debug(v ...interface{}) {
	if loggerLevel > loggerLevelDebug {
		return
	}
	logger.SetPrefix("[DEBUG] ")
	logger.Output(2, fmt.Sprint(v...))
}

// Debugf delegates to the DebugLogger
func Debugf(format string, v ...interface{}) {
	if loggerLevel > loggerLevelDebug {
		return
	}
	logger.SetPrefix("[DEBUG] ")
	logger.Output(2, fmt.Sprintf(format, v...))
}

// Info delegates to the InfoLogger
func Info(v ...interface{}) {
	if loggerLevel > loggerLevelInfo {
		return
	}
	logger.SetPrefix("[INFO] ")
	logger.Output(2, fmt.Sprint(v...))
}

// Infof delegates to the InfoLogger
func Infof(format string, v ...interface{}) {
	if loggerLevel > loggerLevelInfo {
		return
	}
	logger.SetPrefix("[INFO] ")
	logger.Output(2, fmt.Sprintf(format, v...))
}

// Error delegates to the ErrorLogger
func Error(v ...interface{}) {
	logger.SetPrefix("[ERROR] ")
	logger.Output(2, fmt.Sprint(v...))
}

// Errorf delegates to the ErrorLogger
func Errorf(format string, v ...interface{}) {
	logger.SetPrefix("[ERROR] ")
	logger.Output(2, fmt.Sprintf(format, v...))
}
