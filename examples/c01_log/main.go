package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/yggai/ygggo_log"
)

func main() {
	fmt.Println("=== 日志库使用示例 ===")

	// 示例1: 使用默认日志记录器（输出到标准输出）
	fmt.Println("1. 使用默认日志记录器:")
	ygggo_log.Debug("这是一条调试信息")
	ygggo_log.Info("这是一条信息")
	ygggo_log.Warning("这是一条警告信息")
	ygggo_log.Error("这是一条错误信息")

	fmt.Println("\n2. 使用自定义输出的日志记录器:")
	// 示例2: 创建自定义日志记录器（输出到缓冲区）
	var buf bytes.Buffer
	customLogger := ygggo_log.NewLogger(&buf)

	customLogger.Debug("自定义日志记录器的调试信息")
	customLogger.Info("自定义日志记录器的信息")
	customLogger.Warning("自定义日志记录器的警告信息")
	customLogger.Error("自定义日志记录器的错误信息")

	fmt.Printf("缓冲区中的日志内容:\n%s", buf.String())

	fmt.Println("\n3. 使用文件输出的日志记录器:")
	// 示例3: 输出到文件
	file, err := os.Create("app.log")
	if err != nil {
		fmt.Printf("创建日志文件失败: %v\n", err)
		return
	}
	defer file.Close()

	fileLogger := ygggo_log.NewLogger(file)
	fileLogger.Info("这条日志将写入到文件中")
	fileLogger.Warning("文件日志警告信息")
	fileLogger.Error("文件日志错误信息")

	fmt.Println("日志已写入到 app.log 文件中")

	// 读取并显示文件内容
	content, err := os.ReadFile("app.log")
	if err != nil {
		fmt.Printf("读取日志文件失败: %v\n", err)
		return
	}
	fmt.Printf("文件内容:\n%s", string(content))

	fmt.Println("\n4. Panic日志示例 (注意：这会导致程序崩溃):")
	fmt.Println("取消注释下面的代码来测试Panic功能:")
	fmt.Println("// ygggo_log.Panic(\"这是一条panic信息，程序将崩溃\")")

	// 如果要测试panic功能，请取消注释下面这行
	// ygggo_log.Panic("这是一条panic信息，程序将崩溃")

	fmt.Println("\n=== 示例结束 ===")
}
