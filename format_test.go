package ygggo_log

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

// LogEntry 表示JSON格式的日志条目
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

func TestLoadConfigFromEnv_DefaultFormat(t *testing.T) {
	// 清理环境变量
	os.Unsetenv("YGGGO_LOG_FORMAT")

	config := LoadConfigFromEnv()

	if config.Format != TextFormat {
		t.Errorf("Expected default format to be TextFormat, got: %v", config.Format)
	}
}

func TestLoadConfigFromEnv_JsonFormat(t *testing.T) {
	os.Setenv("YGGGO_LOG_FORMAT", "json")
	defer os.Unsetenv("YGGGO_LOG_FORMAT")

	config := LoadConfigFromEnv()

	if config.Format != JsonFormat {
		t.Errorf("Expected format to be JsonFormat, got: %v", config.Format)
	}
}

func TestLoadConfigFromEnv_TextFormat(t *testing.T) {
	os.Setenv("YGGGO_LOG_FORMAT", "text")
	defer os.Unsetenv("YGGGO_LOG_FORMAT")

	config := LoadConfigFromEnv()

	if config.Format != TextFormat {
		t.Errorf("Expected format to be TextFormat, got: %v", config.Format)
	}
}

func TestLoadConfigFromEnv_InvalidFormat(t *testing.T) {
	os.Setenv("YGGGO_LOG_FORMAT", "invalid")
	defer os.Unsetenv("YGGGO_LOG_FORMAT")

	config := LoadConfigFromEnv()

	// 无效格式应该使用默认的文本格式
	if config.Format != TextFormat {
		t.Errorf("Expected invalid format to default to TextFormat, got: %v", config.Format)
	}
}

func TestJsonFormatter_Info(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewJsonFormatter()

	formatter.Format(&buf, InfoLevel, "Test message")

	output := buf.String()

	// 验证输出是有效的JSON
	var entry LogEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}

	// 验证字段内容
	if entry.Level != "INFO" {
		t.Errorf("Expected level to be INFO, got: %s", entry.Level)
	}

	if entry.Message != "Test message" {
		t.Errorf("Expected message to be 'Test message', got: %s", entry.Message)
	}

	// 验证时间戳格式
	if _, err := time.Parse("2006-01-02T15:04:05Z07:00", entry.Timestamp); err != nil {
		t.Errorf("Invalid timestamp format: %s", entry.Timestamp)
	}
}

func TestJsonFormatter_AllLevels(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewJsonFormatter()

	testCases := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarningLevel, "WARNING"},
		{ErrorLevel, "ERROR"},
		{PanicLevel, "PANIC"},
	}

	for _, tc := range testCases {
		buf.Reset()
		formatter.Format(&buf, tc.level, "Test message")

		var entry LogEntry
		if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
			t.Errorf("Failed to parse JSON for level %v: %v", tc.level, err)
			continue
		}

		if entry.Level != tc.expected {
			t.Errorf("Expected level %s, got: %s", tc.expected, entry.Level)
		}
	}
}

func TestTextFormatter_Info(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewTextFormatter()

	formatter.Format(&buf, InfoLevel, "Test message")

	output := buf.String()

	// 验证文本格式
	if !strings.Contains(output, "[INFO]") {
		t.Error("Expected [INFO] in text format output")
	}

	if !strings.Contains(output, "Test message") {
		t.Error("Expected 'Test message' in text format output")
	}

	if !strings.HasSuffix(output, "\n") {
		t.Error("Expected text format output to end with newline")
	}
}

func TestLoggerWithJsonFormat(t *testing.T) {
	// 重置单例
	ResetSingleton()

	os.Setenv("YGGGO_LOG_FORMAT", "json")
	defer os.Unsetenv("YGGGO_LOG_FORMAT")

	var buf bytes.Buffer
	logger := NewLoggerFromEnvWithOutput(&buf)

	logger.Info("Test JSON message")

	output := buf.String()
	t.Logf("JSON output: %q", output)

	// 验证输出是JSON格式
	var entry LogEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Errorf("Logger with JSON format should output valid JSON: %v\nOutput: %s", err, output)
		return
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got: %s", entry.Level)
	}

	if entry.Message != "Test JSON message" {
		t.Errorf("Expected message 'Test JSON message', got: %s", entry.Message)
	}
}

func TestLoggerWithTextFormat(t *testing.T) {
	// 重置单例
	ResetSingleton()

	os.Setenv("YGGGO_LOG_FORMAT", "text")
	defer os.Unsetenv("YGGGO_LOG_FORMAT")

	var buf bytes.Buffer
	logger := NewLoggerFromEnvWithOutput(&buf)

	logger.Info("Test text message")

	output := buf.String()

	// 验证输出是文本格式
	if !strings.Contains(output, "[INFO]") {
		t.Error("Expected text format with [INFO]")
	}

	if !strings.Contains(output, "Test text message") {
		t.Error("Expected message in text format")
	}
}

func TestSingletonWithJsonFormat(t *testing.T) {
	// 重置单例
	ResetSingleton()

	os.Setenv("YGGGO_LOG_FORMAT", "json")
	defer os.Unsetenv("YGGGO_LOG_FORMAT")

	logger1 := GetLogEnv()
	logger2 := GetLogEnv()

	// 验证是同一个对象
	if logger1 != logger2 {
		t.Error("GetLogEnv should return the same singleton instance")
	}
}
