package log

import (
	"os"
	"strings"
	"sync"
	"testing"
)

func TestLogger_Info_WritesToFileAndConsole(t *testing.T) {
	dumpfile, err := os.CreateTemp("", "logtest-*.log")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(dumpfile.Name())

	logger := NewLog(dumpfile.Name(), 1, false)
	defer logger.Close()

	logger.Info("hello %s", "world")

	data, err := os.ReadFile(dumpfile.Name())
	if err != nil {
		t.Fatalf("could not read temp file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "[info] hello world") {
		t.Errorf("log file does not contain expected info message: %s", content)
	}
}

func TestLogger_Error(t *testing.T) {
	dumpfile, err := os.CreateTemp("", "logtest-*.log")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(dumpfile.Name())

	logger := NewLog(dumpfile.Name(), 1, false)
	defer logger.Close()

	logger.Error("some error %d", 42)

	data, err := os.ReadFile(dumpfile.Name())
	if err != nil {
		t.Fatalf("could not read temp file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "[error] some error 42") {
		t.Errorf("log file does not contain expected error message: %s", content)
	}
}

func TestLogger_Concurrency(t *testing.T) {
	dumpfile, err := os.CreateTemp("", "logtest-*.log")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(dumpfile.Name())

	logger := NewLog(dumpfile.Name(), 1, false)
	defer logger.Close()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			logger.Info("goroutine %d", n)
		}(i)
	}
	wg.Wait()

	data, err := os.ReadFile(dumpfile.Name())
	if err != nil {
		t.Fatalf("could not read temp file: %v", err)
	}
	content := string(data)
	count := strings.Count(content, "[info] goroutine")
	if count != 50 {
		t.Errorf("expected 50 log entries, got %d", count)
	}
}

func TestLogger_Close(t *testing.T) {
	dumpfile, err := os.CreateTemp("", "logtest-*.log")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(dumpfile.Name())

	logger := NewLog(dumpfile.Name(), 1, false)
	if err := logger.Close(); err != nil {
		t.Errorf("close failed: %v", err)
	}
}
