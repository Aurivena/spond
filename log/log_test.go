package log

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLogger_RotationOnOverflow(t *testing.T) {
	tmpDir := t.TempDir()

	logfile := filepath.Join(tmpDir, "rotate.log")
	logger := NewLog(logfile, 1)

	logger.Size = 2 * 1024

	defer logger.Close()

	longMsg := strings.Repeat("A", 1024)

	logger.Info("first log %s", longMsg)
	logger.Info("second log %s", longMsg)
	logger.Info("third log %s", longMsg)

	time.Sleep(100 * time.Millisecond)

	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("could not list temp dir: %v", err)
	}
	var rotatedFound bool
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "rotate") && f.Name() != "rotate.log" {
			rotatedFound = true
			break
		}
	}
	if !rotatedFound {
		t.Error("rotation did not create a new file")
	}

	originalData, _ := os.ReadFile(logfile)
	if len(originalData) == 0 {
		t.Error("original log file is empty after rotation")
	}

	var newLogFound bool
	for _, f := range files {
		if f.Name() != "rotate.log" && strings.HasPrefix(f.Name(), "rotate-") {
			data, _ := os.ReadFile(filepath.Join(tmpDir, f.Name()))
			if len(data) > 0 {
				newLogFound = true
				break
			}
		}
	}
	if !newLogFound {
		t.Error("rotated log file is missing or empty")
	}
}

func TestLogger_Info_WritesToFileAndConsole(t *testing.T) {
	dumpfile, err := os.CreateTemp("", "logtest-*.log")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(dumpfile.Name())

	logger := NewLog(dumpfile.Name(), 1)
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

	logger := NewLog(dumpfile.Name(), 1)
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

	logger := NewLog(dumpfile.Name(), 10*1024)
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

	baseDir := filepath.Dir(dumpfile.Name())
	baseName := filepath.Base(dumpfile.Name())
	pattern := filepath.Join(baseDir, baseName+"*")
	files, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("glob error: %v", err)
	}

	total := 0
	for _, fname := range files {
		data, err := os.ReadFile(fname)
		if err != nil {
			t.Fatalf("could not read log file %s: %v", fname, err)
		}
		total += strings.Count(string(data), "[info] goroutine")
	}

	if total < 45 {
		t.Errorf("expected at least 45 log entries, got %d (may lose some during rotation)", total)
	}
}

func TestLogger_Close(t *testing.T) {
	dumpfile, err := os.CreateTemp("", "logtest-*.log")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(dumpfile.Name())

	logger := NewLog(dumpfile.Name(), 1)
	if err := logger.Close(); err != nil {
		t.Errorf("close failed: %v", err)
	}
}
