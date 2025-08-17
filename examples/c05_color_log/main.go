package main

import (
	"fmt"
	"os"

	"github.com/yggai/ygggo_log"
)

func main() {
	fmt.Println("=== 彩色日志示例 ===")

	// 示例1: 默认输出（无彩色）
	fmt.Println("\n1. 默认输出（无彩色）:")
	ygggo_log.ResetSingleton()
	os.Unsetenv("YGGGO_LOG_COLOR")
	os.Unsetenv("YGGGO_LOG_CONSOLE")
	
	logger1 := ygggo_log.GetLogEnv()
	logger1.Debug("这是DEBUG信息")
	logger1.Info("这是INFO信息")
	logger1.Warning("这是WARNING信息")
	logger1.Error("这是ERROR信息")

	// 示例2: 启用彩色输出
	fmt.Println("\n2. 启用彩色输出:")
	ygggo_log.ResetSingleton()
	os.Setenv("YGGGO_LOG_COLOR", "true")
	
	logger2 := ygggo_log.GetLogEnv()
	logger2.Debug("这是彩色的DEBUG信息")
	logger2.Info("这是彩色的INFO信息")
	logger2.Warning("这是彩色的WARNING信息")
	logger2.Error("这是彩色的ERROR信息")

	// 示例3: 同时输出到文件和控制台（彩色）
	fmt.Println("\n3. 同时输出到文件和控制台（彩色）:")
	ygggo_log.ResetSingleton()
	os.Setenv("YGGGO_LOG_COLOR", "true")
	os.Setenv("YGGGO_LOG_CONSOLE", "true")
	os.Setenv("YGGGO_LOG_FILE", "color.log")
	
	logger3 := ygggo_log.GetLogEnv()
	logger3.Info("这条日志同时输出到控制台（彩色）和文件（无彩色）")
	logger3.Warning("控制台显示彩色，文件保存纯文本")
	logger3.Error("这样便于阅读和存储")

	// 读取并显示文件内容
	if content, err := os.ReadFile("color.log"); err == nil {
		fmt.Printf("\n文件内容（无彩色）:\n%s", string(content))
	} else {
		fmt.Printf("读取日志文件失败: %v\n", err)
	}

	// 示例4: 只输出到控制台（彩色）
	fmt.Println("\n4. 只输出到控制台（彩色）:")
	ygggo_log.ResetSingleton()
	os.Setenv("YGGGO_LOG_COLOR", "true")
	os.Setenv("YGGGO_LOG_CONSOLE", "true")
	os.Unsetenv("YGGGO_LOG_FILE")
	
	logger4 := ygggo_log.GetLogEnv()
	logger4.Debug("青色的DEBUG信息")
	logger4.Info("绿色的INFO信息")
	logger4.Warning("黄色的WARNING信息")
	logger4.Error("红色的ERROR信息")

	// 示例5: JSON格式 + 彩色输出
	fmt.Println("\n5. JSON格式 + 彩色输出:")
	ygggo_log.ResetSingleton()
	os.Setenv("YGGGO_LOG_COLOR", "true")
	os.Setenv("YGGGO_LOG_FORMAT", "json")
	
	logger5 := ygggo_log.GetLogEnv()
	logger5.Info("彩色的JSON格式日志")
	logger5.Warning("JSON + 彩色的组合")

	// 示例6: 显示当前配置
	fmt.Println("\n6. 显示当前配置:")
	config := ygggo_log.LoadConfigFromEnv()
	fmt.Printf("日志级别: %s\n", config.Level.String())
	fmt.Printf("日志格式: %s\n", config.Format.String())
	fmt.Printf("输出文件: %s\n", config.OutputFile)
	fmt.Printf("控制台输出: %t\n", config.Console)
	fmt.Printf("彩色输出: %t\n", config.Color)

	// 清理环境变量
	os.Unsetenv("YGGGO_LOG_COLOR")
	os.Unsetenv("YGGGO_LOG_CONSOLE")
	os.Unsetenv("YGGGO_LOG_FILE")
	os.Unsetenv("YGGGO_LOG_FORMAT")

	fmt.Println("\n=== 示例结束 ===")
	fmt.Println("\n说明：")
	fmt.Println("- YGGGO_LOG_COLOR=true 启用彩色输出")
	fmt.Println("- YGGGO_LOG_CONSOLE=true 强制输出到控制台")
	fmt.Println("- 彩色输出只在控制台显示，文件中保存纯文本")
	fmt.Println("- 不同日志级别使用不同颜色：")
	fmt.Println("  * DEBUG: 青色")
	fmt.Println("  * INFO: 绿色")
	fmt.Println("  * WARNING: 黄色")
	fmt.Println("  * ERROR: 红色")
	fmt.Println("  * PANIC: 紫色")
}
