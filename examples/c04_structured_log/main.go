package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/yggai/ygggo_log"
)

func main() {
	fmt.Println("=== 结构化日志示例 ===")

	// 示例1: 文本格式
	fmt.Println("\n1. 文本格式:")
	var textBuf bytes.Buffer
	os.Setenv("YGGGO_LOG_FORMAT", "text")
	textLogger := ygggo_log.NewLoggerFromEnvWithOutput(&textBuf)
	textLogger.Info("这是文本格式的日志")
	textLogger.Warning("这是文本格式的警告")
	textLogger.Error("这是文本格式的错误")
	fmt.Print(textBuf.String())

	// 示例2: JSON格式
	fmt.Println("\n2. JSON格式:")
	var jsonBuf bytes.Buffer
	os.Setenv("YGGGO_LOG_FORMAT", "json")
	jsonLogger := ygggo_log.NewLoggerFromEnvWithOutput(&jsonBuf)
	jsonLogger.Info("这是JSON格式的日志")
	jsonLogger.Warning("这是JSON格式的警告")
	jsonLogger.Error("这是JSON格式的错误")
	fmt.Print(jsonBuf.String())

	// 示例3: JSON格式输出到文件
	fmt.Println("\n3. JSON格式输出到文件:")
	ygggo_log.ResetSingleton()
	os.Setenv("YGGGO_LOG_FORMAT", "json")
	os.Setenv("YGGGO_LOG_FILE", "structured.log")

	logger3 := ygggo_log.GetLogEnv()
	logger3.Info("这条JSON日志将写入文件")
	logger3.Warning("JSON格式的警告信息")
	logger3.Error("JSON格式的错误信息")

	// 读取并显示JSON文件内容
	if content, err := os.ReadFile("structured.log"); err == nil {
		fmt.Printf("JSON日志文件内容:\n%s", string(content))
	} else {
		fmt.Printf("读取JSON日志文件失败: %v\n", err)
	}

	// 示例4: 文本格式输出到文件
	fmt.Println("\n4. 文本格式输出到文件:")
	ygggo_log.ResetSingleton()
	os.Setenv("YGGGO_LOG_FORMAT", "text")
	os.Setenv("YGGGO_LOG_FILE", "text.log")

	logger4 := ygggo_log.GetLogEnv()
	logger4.Info("这条文本日志将写入文件")
	logger4.Warning("文本格式的警告信息")
	logger4.Error("文本格式的错误信息")

	// 读取并显示文本文件内容
	if content, err := os.ReadFile("text.log"); err == nil {
		fmt.Printf("文本日志文件内容:\n%s", string(content))
	} else {
		fmt.Printf("读取文本日志文件失败: %v\n", err)
	}

	// 示例5: 显示当前配置
	fmt.Println("\n5. 显示当前配置:")
	config := ygggo_log.LoadConfigFromEnv()
	fmt.Printf("当前日志级别: %s\n", config.Level.String())
	fmt.Printf("当前日志格式: %s\n", config.Format.String())
	fmt.Printf("当前输出文件: %s\n", config.OutputFile)

	// 清理环境变量
	os.Unsetenv("YGGGO_LOG_FORMAT")
	os.Unsetenv("YGGGO_LOG_FILE")

	fmt.Println("\n=== 示例结束 ===")
	fmt.Println("\n说明：")
	fmt.Println("- 支持两种日志格式：text（文本）和 json（JSON）")
	fmt.Println("- 通过环境变量 YGGGO_LOG_FORMAT 控制格式")
	fmt.Println("- JSON格式便于日志分析和处理")
	fmt.Println("- 文本格式便于人类阅读")
	fmt.Println("- 可以与文件输出结合使用")
}
