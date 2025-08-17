package ygggo_log

import (
	"os"
	"testing"
)

func TestGetLogEnv_Singleton(t *testing.T) {
	// 清理环境变量
	os.Unsetenv("YGGGO_LOG_LEVEL")
	os.Unsetenv("YGGGO_LOG_FILE")
	
	// 多次调用 GetLogEnv，应该返回同一个对象
	logger1 := GetLogEnv()
	logger2 := GetLogEnv()
	logger3 := GetLogEnv()
	
	// 检查是否是同一个对象（比较指针地址）
	if logger1 != logger2 {
		t.Error("GetLogEnv() should return the same instance, but got different instances")
	}
	
	if logger2 != logger3 {
		t.Error("GetLogEnv() should return the same instance, but got different instances")
	}
	
	if logger1 != logger3 {
		t.Error("GetLogEnv() should return the same instance, but got different instances")
	}
}

func TestGetLogEnv_SingletonWithDifferentEnvValues(t *testing.T) {
	// 测试即使环境变量改变，也应该返回同一个对象
	
	// 第一次调用
	os.Setenv("YGGGO_LOG_LEVEL", "DEBUG")
	logger1 := GetLogEnv()
	
	// 改变环境变量后再次调用
	os.Setenv("YGGGO_LOG_LEVEL", "ERROR")
	logger2 := GetLogEnv()
	
	// 应该是同一个对象
	if logger1 != logger2 {
		t.Error("GetLogEnv() should return the same instance even when environment variables change")
	}
	
	// 清理
	os.Unsetenv("YGGGO_LOG_LEVEL")
}

func TestGetLogEnv_SingletonConcurrency(t *testing.T) {
	// 测试并发安全性
	const numGoroutines = 100
	loggers := make([]*Logger, numGoroutines)
	done := make(chan bool, numGoroutines)
	
	// 启动多个goroutine同时调用GetLogEnv
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			loggers[index] = GetLogEnv()
			done <- true
		}(i)
	}
	
	// 等待所有goroutine完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	// 检查所有返回的logger是否是同一个对象
	firstLogger := loggers[0]
	for i := 1; i < numGoroutines; i++ {
		if loggers[i] != firstLogger {
			t.Errorf("Concurrent calls to GetLogEnv() should return the same instance, but got different instances at index %d", i)
		}
	}
}

func TestResetSingleton(t *testing.T) {
	// 测试重置单例的功能（用于测试）
	
	// 获取第一个实例
	logger1 := GetLogEnv()
	
	// 重置单例
	ResetSingleton()
	
	// 获取新的实例
	logger2 := GetLogEnv()
	
	// 应该是不同的对象
	if logger1 == logger2 {
		t.Error("After ResetSingleton(), GetLogEnv() should return a new instance")
	}
}

func TestSingletonConfiguration(t *testing.T) {
	// 测试单例对象的配置是否正确应用
	
	// 重置单例
	ResetSingleton()
	
	// 设置环境变量
	os.Setenv("YGGGO_LOG_LEVEL", "ERROR")
	defer os.Unsetenv("YGGGO_LOG_LEVEL")
	
	// 获取单例
	logger := GetLogEnv()
	
	// 检查配置是否正确应用
	if logger.minLevel != ErrorLevel {
		t.Errorf("Expected singleton logger to have ERROR level, got: %v", logger.minLevel)
	}
}

func TestSingletonConfigurationOnlyAppliedOnce(t *testing.T) {
	// 测试配置只在第一次创建时应用
	
	// 重置单例
	ResetSingleton()
	
	// 第一次设置环境变量并获取单例
	os.Setenv("YGGGO_LOG_LEVEL", "DEBUG")
	logger1 := GetLogEnv()
	
	// 改变环境变量
	os.Setenv("YGGGO_LOG_LEVEL", "ERROR")
	logger2 := GetLogEnv()
	
	// 应该是同一个对象，且配置应该保持第一次的设置
	if logger1 != logger2 {
		t.Error("Should return the same singleton instance")
	}
	
	if logger2.minLevel != DebugLevel {
		t.Errorf("Singleton configuration should not change after first initialization, expected DEBUG level, got: %v", logger2.minLevel)
	}
	
	// 清理
	os.Unsetenv("YGGGO_LOG_LEVEL")
}
