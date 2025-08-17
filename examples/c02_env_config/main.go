package main

import (
	"fmt"
	"os"

	"github.com/yggai/ygggo_log"
)

func main() {
	fmt.Println("=== .env 文件配置日志示例 ===")

	// 显示当前工作目录
	if wd, err := os.Getwd(); err == nil {
		fmt.Printf("当前工作目录: %s\n", wd)
	}

	// 检查 .env 文件是否存在
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("警告: .env 文件不存在，将使用默认配置")
	} else {
		fmt.Println("找到 .env 文件，正在加载配置...")
	}

	// 使用 GetLogEnv 创建日志记录器，它会自动调用 ygggo_env.LoadEnv()
	logger := ygggo_log.GetLogEnv()

	// 显示当前配置
	config := ygggo_log.LoadConfigFromEnv()
	fmt.Printf("当前日志级别: %s\n", config.Level.String())
	fmt.Printf("当前输出文件: %s\n", config.OutputFile)
	if config.OutputFile == "" {
		fmt.Println("（空字符串表示输出到标准输出）")
	}

	fmt.Println("\n开始记录日志:")

	// 记录不同级别的日志
	logger.Debug("这是DEBUG级别的日志")
	logger.Info("这是INFO级别的日志")
	logger.Warning("这是WARNING级别的日志")
	logger.Error("这是ERROR级别的日志")
	// 如果配置了输出文件，显示文件内容
	if config.OutputFile != "" {
		fmt.Printf("\n检查日志文件 %s:\n", config.OutputFile)
		if content, err := os.ReadFile(config.OutputFile); err == nil {
			fmt.Printf("文件内容:\n%s", string(content))
		} else {
			fmt.Printf("读取日志文件失败: %v\n", err)
		}
	}

	fmt.Println("\n=== 示例结束 ===")
	fmt.Println("\n提示：")
	fmt.Println("1. 修改 .env 文件中的 YGGGO_LOG_LEVEL 来改变日志级别")
	fmt.Println("2. 修改 .env 文件中的 YGGGO_LOG_FILE 来改变输出文件")
	fmt.Println("3. 重新运行程序查看效果")
}
