package ygggo_log

import (
	"bytes"
	"strings"
	"testing"
)

func TestDebug(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)
	
	logger.Debug("This is a debug message")
	
	output := buf.String()
	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Expected DEBUG level in output, got: %s", output)
	}
	if !strings.Contains(output, "This is a debug message") {
		t.Errorf("Expected debug message in output, got: %s", output)
	}
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)
	
	logger.Info("This is an info message")
	
	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected INFO level in output, got: %s", output)
	}
	if !strings.Contains(output, "This is an info message") {
		t.Errorf("Expected info message in output, got: %s", output)
	}
}

func TestWarning(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)
	
	logger.Warning("This is a warning message")
	
	output := buf.String()
	if !strings.Contains(output, "WARNING") {
		t.Errorf("Expected WARNING level in output, got: %s", output)
	}
	if !strings.Contains(output, "This is a warning message") {
		t.Errorf("Expected warning message in output, got: %s", output)
	}
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)
	
	logger.Error("This is an error message")
	
	output := buf.String()
	if !strings.Contains(output, "ERROR") {
		t.Errorf("Expected ERROR level in output, got: %s", output)
	}
	if !strings.Contains(output, "This is an error message") {
		t.Errorf("Expected error message in output, got: %s", output)
	}
}

func TestPanic(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)
	
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but function did not panic")
		}
	}()
	
	logger.Panic("This is a panic message")
}

func TestPanicMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)
	
	defer func() {
		if r := recover(); r != nil {
			output := buf.String()
			if !strings.Contains(output, "PANIC") {
				t.Errorf("Expected PANIC level in output, got: %s", output)
			}
			if !strings.Contains(output, "This is a panic message") {
				t.Errorf("Expected panic message in output, got: %s", output)
			}
		}
	}()
	
	logger.Panic("This is a panic message")
}

func TestLogFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)
	
	logger.Info("Test message")
	
	output := buf.String()
	// 检查日志格式是否包含时间戳
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Expected [INFO] format in output, got: %s", output)
	}
	// 检查是否有换行符
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("Expected output to end with newline, got: %s", output)
	}
}
