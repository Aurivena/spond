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
	formatTime = "20060102-150405"
)

var (
	_        io.WriteCloser = (*Logger)(nil)
	megabyte int64          = 1024 * 1024
)

// Usage example
// logger:=NewLog("log/io.log",50)
// defer logger.Close()
type Logger struct {
	Filename string // path storage
	Size     int64  // max size file - megabytes
	base     string
	mu       sync.Mutex
	file     *os.File
}

func NewLog(filename string, size int64) *Logger {
	filename, base := setFilename(filename)
	size = setSize(size)

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		log.Fatalf("can't create log dir: %v", err)
	}

	f := createFile(filename)

	return &Logger{
		Filename: filename,
		base:     base,
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
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fileLimit()
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
	resLog := format(logType, logType, msg, args...)
	writeToConsole(resLog)
}

func (l *Logger) fileLimit() {
	fileInfo, err := os.Stat(l.Filename)
	if err != nil {
		log.Fatalf("file not open")
	}
	size := fileInfo.Size()

	if size >= l.Size*megabyte {
		_ = l.file.Close()
		l.rotation()
		l.file = createFile(l.Filename)
	}
}

func (l *Logger) rotation() {
	baseDir := filepath.Dir(l.Filename)
	baseName := filepath.Base(l.base)
	ext := filepath.Ext(l.Filename)
	nameOnly := baseName[:len(baseName)-len(ext)]
	newName := fmt.Sprintf("%s-%s%s", nameOnly, time.Now().UTC().Format(formatTime), ext)
	l.Filename = filepath.Join(baseDir, newName)
}

func setFilename(filename string) (string, string) {
	if filename == "" {
		return filepath.Join("log", defaultLogFilename()), "spond.log"
	}
	if len(filename) > 0 && (filename[len(filename)-1] == '/' || filename[len(filename)-1] == '\\') {
		return filepath.Join(filename, defaultLogFilename()), "spond.log"
	}
	if stat, err := os.Stat(filename); err == nil && stat.IsDir() {
		return filepath.Join(filename, defaultLogFilename()), "spond.log"
	}
	return filename, filepath.Base(filename)
}

func defaultLogFilename() string {
	return time.Now().UTC().Format("20060102-150405") + "-" + "spond" + ".log"
}

func setSize(size int64) int64 {
	if size <= 0 {
		return defaultMB * megabyte
	}
	return size
}

func createFile(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("can't open log file: %v", err)
	}
	return f
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
	currentTime := time.Now().UTC().Format(formatTime)
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
