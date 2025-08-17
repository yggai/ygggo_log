package ygggo_log

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigFromEnv_DefaultFileSize(t *testing.T) {
	// 清理环境变量
	os.Unsetenv("YGGGO_LOG_FILE_SIZE")
	
	config := LoadConfigFromEnv()
	
	if config.FileSize != 100*1024*1024 { // 默认100MB
		t.Errorf("Expected default file size to be 100MB, got: %d", config.FileSize)
	}
}

func TestLoadConfigFromEnv_DefaultFileNum(t *testing.T) {
	// 清理环境变量
	os.Unsetenv("YGGGO_LOG_FILE_NUM")
	
	config := LoadConfigFromEnv()
	
	if config.FileNum != 3 { // 默认3个文件
		t.Errorf("Expected default file num to be 3, got: %d", config.FileNum)
	}
}

func TestLoadConfigFromEnv_FileSizeParsing(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"100M", 100 * 1024 * 1024},
		{"50MB", 50 * 1024 * 1024},
		{"1G", 1024 * 1024 * 1024},
		{"1GB", 1024 * 1024 * 1024},
		{"500K", 500 * 1024},
		{"500KB", 500 * 1024},
		{"1024", 1024}, // 纯数字，按字节处理
		{"invalid", 100 * 1024 * 1024}, // 无效值，使用默认值
	}
	
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			os.Setenv("YGGGO_LOG_FILE_SIZE", tc.input)
			defer os.Unsetenv("YGGGO_LOG_FILE_SIZE")
			
			config := LoadConfigFromEnv()
			
			if config.FileSize != tc.expected {
				t.Errorf("Expected file size %d for input %s, got: %d", 
					tc.expected, tc.input, config.FileSize)
			}
		})
	}
}

func TestLoadConfigFromEnv_FileNumParsing(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"5", 5},
		{"10", 10},
		{"1", 1},
		{"0", 3},       // 无效值，使用默认值
		{"invalid", 3}, // 无效值，使用默认值
	}
	
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			os.Setenv("YGGGO_LOG_FILE_NUM", tc.input)
			defer os.Unsetenv("YGGGO_LOG_FILE_NUM")
			
			config := LoadConfigFromEnv()
			
			if config.FileNum != tc.expected {
				t.Errorf("Expected file num %d for input %s, got: %d", 
					tc.expected, tc.input, config.FileNum)
			}
		})
	}
}

func TestRotatingWriter_Write(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")
	
	// 创建小文件大小的轮转写入器（100字节）
	writer, err := NewRotatingWriter(logFile, 100, 3)
	if err != nil {
		t.Fatalf("Failed to create rotating writer: %v", err)
	}
	defer writer.Close()
	
	// 写入数据，超过文件大小限制
	data := strings.Repeat("Hello World!\n", 10) // 约130字节
	
	n, err := writer.Write([]byte(data))
	if err != nil {
		t.Errorf("Failed to write data: %v", err)
	}
	
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, got: %d", len(data), n)
	}
	
	// 检查是否创建了轮转文件
	files, err := filepath.Glob(filepath.Join(tempDir, "test.log*"))
	if err != nil {
		t.Fatalf("Failed to list log files: %v", err)
	}
	
	if len(files) < 1 {
		t.Error("Expected at least one log file to be created")
	}
}

func TestRotatingWriter_Rotation(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "rotate.log")
	
	// 创建小文件大小的轮转写入器（50字节，最多2个文件）
	writer, err := NewRotatingWriter(logFile, 50, 2)
	if err != nil {
		t.Fatalf("Failed to create rotating writer: %v", err)
	}
	defer writer.Close()
	
	// 写入多次数据，触发轮转
	for i := 0; i < 5; i++ {
		data := strings.Repeat("Test data line\n", 2) // 约30字节
		writer.Write([]byte(data))
	}
	
	// 检查文件数量不超过限制
	files, err := filepath.Glob(filepath.Join(tempDir, "rotate.log*"))
	if err != nil {
		t.Fatalf("Failed to list log files: %v", err)
	}
	
	if len(files) > 2 {
		t.Errorf("Expected at most 2 log files, got: %d", len(files))
	}
}

func TestRotatingWriter_FileNaming(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "naming.log")
	
	writer, err := NewRotatingWriter(logFile, 30, 3)
	if err != nil {
		t.Fatalf("Failed to create rotating writer: %v", err)
	}
	defer writer.Close()
	
	// 写入数据触发轮转
	for i := 0; i < 3; i++ {
		data := strings.Repeat("Data\n", 10) // 50字节
		writer.Write([]byte(data))
	}
	
	// 检查文件命名
	files, err := filepath.Glob(filepath.Join(tempDir, "naming.log*"))
	if err != nil {
		t.Fatalf("Failed to list log files: %v", err)
	}
	
	// 应该有当前文件和轮转文件
	hasCurrentFile := false
	hasRotatedFile := false
	
	for _, file := range files {
		basename := filepath.Base(file)
		if basename == "naming.log" {
			hasCurrentFile = true
		} else if strings.HasPrefix(basename, "naming.log.") {
			hasRotatedFile = true
		}
	}
	
	if !hasCurrentFile {
		t.Error("Expected current log file to exist")
	}
	
	if !hasRotatedFile {
		t.Error("Expected rotated log file to exist")
	}
}

func TestLoggerWithRotation(t *testing.T) {
	// 重置单例
	ResetSingleton()
	
	// 创建临时目录
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "logger_rotate.log")
	
	os.Setenv("YGGGO_LOG_FILE", logFile)
	os.Setenv("YGGGO_LOG_FILE_SIZE", "100")  // 100字节
	os.Setenv("YGGGO_LOG_FILE_NUM", "2")     // 2个文件
	defer os.Unsetenv("YGGGO_LOG_FILE")
	defer os.Unsetenv("YGGGO_LOG_FILE_SIZE")
	defer os.Unsetenv("YGGGO_LOG_FILE_NUM")
	
	logger := GetLogEnv()
	
	// 写入多条日志，触发轮转
	for i := 0; i < 10; i++ {
		logger.Info("This is a test log message for rotation testing")
	}
	
	// 检查文件是否创建
	files, err := filepath.Glob(filepath.Join(tempDir, "logger_rotate.log*"))
	if err != nil {
		t.Fatalf("Failed to list log files: %v", err)
	}
	
	if len(files) == 0 {
		t.Error("Expected log files to be created")
	}
	
	if len(files) > 2 {
		t.Errorf("Expected at most 2 log files, got: %d", len(files))
	}
}

func TestParseSizeString(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{"100", 100},
		{"1K", 1024},
		{"1KB", 1024},
		{"1M", 1024 * 1024},
		{"1MB", 1024 * 1024},
		{"1G", 1024 * 1024 * 1024},
		{"1GB", 1024 * 1024 * 1024},
		{"invalid", 100 * 1024 * 1024}, // 默认值
	}
	
	for _, tc := range testCases {
		result := parseSizeString(tc.input)
		if result != tc.expected {
			t.Errorf("Expected %d for input %s, got: %d", tc.expected, tc.input, result)
		}
	}
}
