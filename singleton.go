package ygggo_log

import (
	"sync"
)

var (
	// singletonLogger 单例日志记录器实例
	singletonLogger *Logger
	
	// once 确保单例只被初始化一次
	once sync.Once
	
	// mutex 用于保护重置操作的并发安全
	mutex sync.RWMutex
)

// GetLogEnv 获取单例日志记录器
// 无论从代码的任何地方调用此方法，都会返回同一个日志对象
func GetLogEnv() *Logger {
	mutex.RLock()
	if singletonLogger != nil {
		defer mutex.RUnlock()
		return singletonLogger
	}
	mutex.RUnlock()
	
	// 使用 sync.Once 确保只初始化一次
	once.Do(func() {
		mutex.Lock()
		defer mutex.Unlock()
		
		// 双重检查锁定模式
		if singletonLogger == nil {
			config := LoadConfigFromEnv()
			singletonLogger = NewLoggerFromConfig(config)
		}
	})
	
	mutex.RLock()
	defer mutex.RUnlock()
	return singletonLogger
}

// ResetSingleton 重置单例实例（主要用于测试）
// 注意：这个方法不是线程安全的，只应该在测试环境中使用
func ResetSingleton() {
	mutex.Lock()
	defer mutex.Unlock()
	
	singletonLogger = nil
	// 重置 sync.Once，使其可以再次执行
	once = sync.Once{}
}
