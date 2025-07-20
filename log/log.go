package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
)

const (
	info       = "info"
	warn       = "warn"
	test       = "test"
	err        = "error"
	debug      = "debug"
	fileType   = "file"
	loggerType = "logger"
	defaultMB  = 10
)

var (
	_        io.WriteCloser = (*Logger)(nil)
	megabyte                = 1024 * 1024
)

// Usage example
// logger:=NewLog("log/io.log",50)
// defer logger.Close()
type Logger struct {
	Filename string // path storage
	Size     int    // max size file - megabytes
	my       sync.Mutex
	file     *os.File
}

func NewLog(filename string, size int) *Logger {
	filename = setFilename(filename)
	size = setSize(size)

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		log.Fatalf("can't create log dir: %v", err)
	}
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("can't open log file: %v", err)
	}
	return &Logger{
		Filename: filename,
		Size:     size,
		file:     f,
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(info, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(warn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(err, msg, args...)
}

func (l *Logger) Test(msg string, args ...any) {
	l.log(test, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(debug, msg, args...)
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.my.Lock()
	defer l.my.Unlock()
	n, err = l.file.Write(p)
	return n, err
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func (l *Logger) log(logType, msg string, args ...any) {
	resFile := format(logType, fileType, msg, args...)
	_, err := l.Write([]byte(resFile))
	if err != nil {
		return
	}
	l.my.Lock()
	defer l.my.Unlock()
	resLog := format(logType, logType, msg, args...)
	writeToConsole(resLog)
}

func setFilename(filename string) string {
	if filename == "" {
		return fmt.Sprintf("log/%s-spond.log", time.Now().UTC().Format(time.RFC3339))
	}

	return filename
}

func setSize(size int) int {
	if size <= 0 {
		return defaultMB
	}
	return size
}

func colored(logType string) string {
	switch logType {
	case info:
		return color.New(color.FgGreen, color.Bold).Sprint(logType)
	case warn, test:
		return color.New(color.FgYellow, color.Bold).Sprint(logType)
	case err:
		return color.New(color.FgRed, color.Bold).Sprint(logType)
	case debug:
		return color.New(color.FgCyan, color.Bold).Sprint(logType)
	default:
		return logType
	}
}

func writeToConsole(msg string) {
	os.Stdout.WriteString(msg)
}

func format(logType, formatType, msg string, args ...any) string {
	currentTime := time.Now().UTC().Format(time.RFC3339)
	notification := fmt.Sprintf(msg, args...)
	switch formatType {
	case fileType:
		return fmt.Sprintf("%s [%s] %s\n", currentTime, logType, notification)
	case loggerType:
		return fmt.Sprintf("%s [%s] %s\n", currentTime, colored(logType), notification)
	default:
		return fmt.Sprintf("%s [error] formatType is invalid\n", currentTime)
	}

}
