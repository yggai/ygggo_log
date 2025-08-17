package ygggo_log

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLoadConfigFromEnv_DefaultValues(t *testing.T) {
	// 清理环境变量
	os.Unsetenv("YGGGO_LOG_LEVEL")
	os.Unsetenv("YGGGO_LOG_FILE")

	config := LoadConfigFromEnv()

	if config.Level != InfoLevel {
		t.Errorf("Expected default log level to be InfoLevel, got: %v", config.Level)
	}

	if config.OutputFile != "" {
		t.Errorf("Expected default output file to be empty, got: %s", config.OutputFile)
	}
}

func TestLoadConfigFromEnv_LogLevel(t *testing.T) {
	tests := []struct {
		envValue      string
		expectedLevel LogLevel
	}{
		{"DEBUG", DebugLevel},
		{"INFO", InfoLevel},
		{"WARNING", WarningLevel},
		{"ERROR", ErrorLevel},
		{"PANIC", PanicLevel},
		{"debug", DebugLevel}, // 测试大小写不敏感
		{"info", InfoLevel},
		{"INVALID", InfoLevel}, // 无效值应该使用默认值
	}

	for _, test := range tests {
		t.Run(test.envValue, func(t *testing.T) {
			os.Setenv("YGGGO_LOG_LEVEL", test.envValue)
			defer os.Unsetenv("YGGGO_LOG_LEVEL")

			config := LoadConfigFromEnv()

			if config.Level != test.expectedLevel {
				t.Errorf("Expected log level %v for env value %s, got: %v",
					test.expectedLevel, test.envValue, config.Level)
			}
		})
	}
}

func TestLoadConfigFromEnv_LogFile(t *testing.T) {
	testFile := "test.log"
	os.Setenv("YGGGO_LOG_FILE", testFile)
	defer os.Unsetenv("YGGGO_LOG_FILE")

	config := LoadConfigFromEnv()

	if config.OutputFile != testFile {
		t.Errorf("Expected output file to be %s, got: %s", testFile, config.OutputFile)
	}
}

func TestGetLogEnv_DefaultConfig(t *testing.T) {
	// 清理环境变量
	os.Unsetenv("YGGGO_LOG_LEVEL")
	os.Unsetenv("YGGGO_LOG_FILE")

	logger := GetLogEnv()

	if logger == nil {
		t.Error("Expected logger to be created, got nil")
	}
}

func TestGetLogEnv_WithFileOutput(t *testing.T) {
	testFile := "test_env.log"
	os.Setenv("YGGGO_LOG_FILE", testFile)
	defer os.Unsetenv("YGGGO_LOG_FILE")
	defer os.Remove(testFile) // 清理测试文件

	logger := GetLogEnv()

	if logger == nil {
		t.Error("Expected logger to be created, got nil")
	}

	// 测试写入文件
	logger.Info("Test message")

	// 检查文件是否存在
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Expected log file %s to be created", testFile)
	}
}

func TestLoggerWithEnvLevel_Debug(t *testing.T) {
	os.Setenv("YGGGO_LOG_LEVEL", "DEBUG")
	defer os.Unsetenv("YGGGO_LOG_LEVEL")

	var buf bytes.Buffer
	logger := NewLoggerFromEnvWithOutput(&buf)

	logger.Debug("Debug message")
	logger.Info("Info message")

	output := buf.String()
	if !strings.Contains(output, "DEBUG") {
		t.Error("Expected DEBUG message to be logged")
	}
	if !strings.Contains(output, "INFO") {
		t.Error("Expected INFO message to be logged")
	}
}

func TestLoggerWithEnvLevel_Error(t *testing.T) {
	os.Setenv("YGGGO_LOG_LEVEL", "ERROR")
	defer os.Unsetenv("YGGGO_LOG_LEVEL")

	var buf bytes.Buffer
	logger := NewLoggerFromEnvWithOutput(&buf)

	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Error("Error message")

	output := buf.String()
	if strings.Contains(output, "DEBUG") {
		t.Error("Expected DEBUG message to be filtered out")
	}
	if strings.Contains(output, "INFO") {
		t.Error("Expected INFO message to be filtered out")
	}
	if !strings.Contains(output, "ERROR") {
		t.Error("Expected ERROR message to be logged")
	}
}

func TestConfigStruct(t *testing.T) {
	config := &LogConfig{
		Level:      DebugLevel,
		OutputFile: "test.log",
	}

	if config.Level != DebugLevel {
		t.Errorf("Expected level to be DebugLevel, got: %v", config.Level)
	}

	if config.OutputFile != "test.log" {
		t.Errorf("Expected output file to be test.log, got: %s", config.OutputFile)
	}
}
