package main

import (
	"fmt"

	"github.com/yggai/ygggo_log"
)

func main() {
	fmt.Println("=== 单例日志对象示例 ===")

	// 示例1: 多次调用 GetLogEnv 获取同一个对象
	fmt.Println("\n1. 验证单例模式:")
	logger1 := ygggo_log.GetLogEnv()
	logger2 := ygggo_log.GetLogEnv()
	logger3 := ygggo_log.GetLogEnv()

	fmt.Printf("logger1 地址: %p\n", logger1)
	fmt.Printf("logger2 地址: %p\n", logger2)
	fmt.Printf("logger3 地址: %p\n", logger3)

	if logger1 == logger2 && logger2 == logger3 {
		fmt.Println("✓ 验证成功：所有调用返回同一个对象")
	} else {
		fmt.Println("✗ 验证失败：返回了不同的对象")
	}

	// 示例2: 在不同函数中调用
	fmt.Println("\n2. 在不同函数中调用:")
	loggerFromFunc1 := getLoggerFromFunction1()
	loggerFromFunc2 := getLoggerFromFunction2()

	fmt.Printf("函数1中的logger地址: %p\n", loggerFromFunc1)
	fmt.Printf("函数2中的logger地址: %p\n", loggerFromFunc2)

	if loggerFromFunc1 == loggerFromFunc2 {
		fmt.Println("✓ 验证成功：不同函数中获取的是同一个对象")
	} else {
		fmt.Println("✗ 验证失败：不同函数中获取了不同的对象")
	}

	// 示例3: 使用单例日志记录器
	fmt.Println("\n3. 使用单例日志记录器:")
	logger := ygggo_log.GetLogEnv()
	logger.Info("这是来自单例日志记录器的信息")
	logger.Warning("这是来自单例日志记录器的警告")
	logger.Error("这是来自单例日志记录器的错误")

	// 示例4: 在不同地方使用，都是同一个对象
	fmt.Println("\n4. 在不同地方使用日志:")
	useLoggerInDifferentPlaces()

	fmt.Println("\n=== 示例结束 ===")
	fmt.Println("\n说明：")
	fmt.Println("- GetLogEnv() 方法实现了单例模式")
	fmt.Println("- 无论在代码的任何地方调用，都会返回同一个日志对象")
	fmt.Println("- 这确保了日志配置的一致性和资源的有效利用")
	fmt.Println("- 单例模式是线程安全的，支持并发调用")
}

func getLoggerFromFunction1() *ygggo_log.Logger {
	fmt.Println("在函数1中获取日志记录器...")
	logger := ygggo_log.GetLogEnv()
	logger.Debug("这是来自函数1的调试信息")
	return logger
}

func getLoggerFromFunction2() *ygggo_log.Logger {
	fmt.Println("在函数2中获取日志记录器...")
	logger := ygggo_log.GetLogEnv()
	logger.Debug("这是来自函数2的调试信息")
	return logger
}

func useLoggerInDifferentPlaces() {
	// 模拟在不同的代码位置使用日志
	
	// 位置1
	func() {
		logger := ygggo_log.GetLogEnv()
		logger.Info("位置1：匿名函数中的日志")
	}()
	
	// 位置2
	go func() {
		logger := ygggo_log.GetLogEnv()
		logger.Info("位置2：goroutine中的日志")
	}()
	
	// 位置3
	type MyStruct struct{}
	myStruct := MyStruct{}
	func(ms MyStruct) {
		logger := ygggo_log.GetLogEnv()
		logger.Info("位置3：结构体方法中的日志")
	}(myStruct)
	
	// 等待goroutine完成
	// 在实际应用中，你可能需要使用sync.WaitGroup或其他同步机制
	fmt.Println("所有位置都使用了同一个单例日志对象")
}
