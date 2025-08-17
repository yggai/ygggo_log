package ygggo_log

import (
	"os"
	"testing"
)

// Test defaults via InitLogEnv conventions
func TestDefaults_InitLogEnv(t *testing.T) {
	ResetSingleton()
	// 清空相关环境变量，确保走默认约定
	os.Unsetenv("YGGGO_LOG_LEVEL")
	os.Unsetenv("YGGGO_LOG_FILE")
	os.Unsetenv("YGGGO_LOG_FORMAT")
	os.Unsetenv("YGGGO_LOG_CONSOLE")
	os.Unsetenv("YGGGO_LOG_COLOR")

	InitLogEnv()
	lg := GetLogEnv()
	if lg == nil { t.Fatal("logger should not be nil") }
	// 至少确保默认级别为 INFO
	if lg.minLevel != InfoLevel {
		t.Fatalf("expected default level INFO, got %v", lg.minLevel)
	}
}

