package ygggo_log

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLoadConfigFromEnv_DefaultConsole(t *testing.T) {
	// 清理环境变量
	os.Unsetenv("YGGGO_LOG_CONSOLE")
	
	config := LoadConfigFromEnv()
	
	if config.Console != false {
		t.Errorf("Expected default console to be false, got: %v", config.Console)
	}
}

func TestLoadConfigFromEnv_ConsoleTrue(t *testing.T) {
	os.Setenv("YGGGO_LOG_CONSOLE", "true")
	defer os.Unsetenv("YGGGO_LOG_CONSOLE")
	
	config := LoadConfigFromEnv()
	
	if config.Console != true {
		t.Errorf("Expected console to be true, got: %v", config.Console)
	}
}

func TestLoadConfigFromEnv_ConsoleFalse(t *testing.T) {
	os.Setenv("YGGGO_LOG_CONSOLE", "false")
	defer os.Unsetenv("YGGGO_LOG_CONSOLE")
	
	config := LoadConfigFromEnv()
	
	if config.Console != false {
		t.Errorf("Expected console to be false, got: %v", config.Console)
	}
}

func TestLoadConfigFromEnv_DefaultColor(t *testing.T) {
	// 清理环境变量
	os.Unsetenv("YGGGO_LOG_COLOR")
	
	config := LoadConfigFromEnv()
	
	if config.Color != false {
		t.Errorf("Expected default color to be false, got: %v", config.Color)
	}
}

func TestLoadConfigFromEnv_ColorTrue(t *testing.T) {
	os.Setenv("YGGGO_LOG_COLOR", "true")
	defer os.Unsetenv("YGGGO_LOG_COLOR")
	
	config := LoadConfigFromEnv()
	
	if config.Color != true {
		t.Errorf("Expected color to be true, got: %v", config.Color)
	}
}

func TestLoadConfigFromEnv_ColorFalse(t *testing.T) {
	os.Setenv("YGGGO_LOG_COLOR", "false")
	defer os.Unsetenv("YGGGO_LOG_COLOR")
	
	config := LoadConfigFromEnv()
	
	if config.Color != false {
		t.Errorf("Expected color to be false, got: %v", config.Color)
	}
}

func TestColorFormatter_Info(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewColorFormatter()
	
	formatter.Format(&buf, InfoLevel, "Test message")
	
	output := buf.String()
	
	// 验证包含ANSI颜色代码
	if !strings.Contains(output, "\033[") {
		t.Error("Expected color formatter to include ANSI color codes")
	}
	
	// 验证包含消息内容
	if !strings.Contains(output, "Test message") {
		t.Error("Expected output to contain the message")
	}
	
	// 验证包含INFO级别
	if !strings.Contains(output, "INFO") {
		t.Error("Expected output to contain INFO level")
	}
}

func TestColorFormatter_AllLevels(t *testing.T) {
	var buf bytes.Buffer
	formatter := NewColorFormatter()
	
	testCases := []struct {
		level LogLevel
		name  string
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
		
		output := buf.String()
		
		// 验证包含ANSI颜色代码
		if !strings.Contains(output, "\033[") {
			t.Errorf("Expected color formatter for %s to include ANSI color codes", tc.name)
		}
		
		// 验证包含级别名称
		if !strings.Contains(output, tc.name) {
			t.Errorf("Expected output to contain %s level", tc.name)
		}
	}
}

func TestLoggerWithColorConsole(t *testing.T) {
	// 重置单例
	ResetSingleton()
	
	os.Setenv("YGGGO_LOG_CONSOLE", "true")
	os.Setenv("YGGGO_LOG_COLOR", "true")
	defer os.Unsetenv("YGGGO_LOG_CONSOLE")
	defer os.Unsetenv("YGGGO_LOG_COLOR")
	
	var buf bytes.Buffer
	logger := NewLoggerFromEnvWithOutput(&buf)
	
	logger.Info("Test color message")
	
	output := buf.String()
	
	// 验证输出包含颜色代码
	if !strings.Contains(output, "\033[") {
		t.Error("Expected colored output to contain ANSI color codes")
	}
}

func TestLoggerWithConsoleOnly(t *testing.T) {
	// 重置单例
	ResetSingleton()
	
	os.Setenv("YGGGO_LOG_CONSOLE", "true")
	os.Setenv("YGGGO_LOG_COLOR", "false")
	defer os.Unsetenv("YGGGO_LOG_CONSOLE")
	defer os.Unsetenv("YGGGO_LOG_COLOR")
	
	var buf bytes.Buffer
	logger := NewLoggerFromEnvWithOutput(&buf)
	
	logger.Info("Test console message")
	
	output := buf.String()
	
	// 验证输出不包含颜色代码
	if strings.Contains(output, "\033[") {
		t.Error("Expected non-colored output to not contain ANSI color codes")
	}
	
	// 验证包含消息内容
	if !strings.Contains(output, "Test console message") {
		t.Error("Expected output to contain the message")
	}
}

func TestColorCodes(t *testing.T) {
	testCases := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, "\033[36m"},    // 青色
		{InfoLevel, "\033[32m"},     // 绿色
		{WarningLevel, "\033[33m"},  // 黄色
		{ErrorLevel, "\033[31m"},    // 红色
		{PanicLevel, "\033[35m"},    // 紫色
	}
	
	for _, tc := range testCases {
		color := getColorCode(tc.level)
		if color != tc.expected {
			t.Errorf("Expected color code %s for level %v, got: %s", tc.expected, tc.level, color)
		}
	}
}

func TestResetColorCode(t *testing.T) {
	reset := getResetCode()
	expected := "\033[0m"
	
	if reset != expected {
		t.Errorf("Expected reset code %s, got: %s", expected, reset)
	}
}

func TestMultiWriter(t *testing.T) {
	// 重置单例
	ResetSingleton()
	
	// 设置控制台输出和文件输出
	os.Setenv("YGGGO_LOG_CONSOLE", "true")
	os.Setenv("YGGGO_LOG_FILE", "test_multi.log")
	defer os.Unsetenv("YGGGO_LOG_CONSOLE")
	defer os.Unsetenv("YGGGO_LOG_FILE")
	defer os.Remove("test_multi.log")
	
	logger := GetLogEnv()
	logger.Info("Test multi output")
	
	// 检查文件是否创建
	if _, err := os.Stat("test_multi.log"); os.IsNotExist(err) {
		t.Error("Expected log file to be created when both console and file output are enabled")
	}
}
