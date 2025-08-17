package main

import (
	"bytes"
	"fmt"

	"github.com/yggai/ygggo_log"
)

func main() {
	fmt.Println("=== 日志库使用示例 ===")

	// 示例1: 使用默认日志记录器（输出到标准输出）
	fmt.Println("1. 使用默认日志记录器:")
	ygggo_log.Debug("这是一条调试信息", map[string]any{"a": 1}, "b=1.1", true)
	ygggo_log.Info("这是一条信息", "user=alice", 123)
	ygggo_log.Warning("这是一条警告信息", map[string]any{"warn": true})
	ygggo_log.Error("这是一条错误信息", "code=E123", 3.14)

	fmt.Println("\n2. 使用自定义输出的日志记录器:")
	// 示例2: 创建自定义日志记录器（输出到缓冲区）
	var buf bytes.Buffer
	customLogger := ygggo_log.NewLogger(&buf)

	customLogger.Debug("自定义日志记录器的调试信息")
	customLogger.Info("自定义日志记录器的信息")
	customLogger.Warning("自定义日志记录器的警告信息")
	customLogger.Error("自定义日志记录器的错误信息")

	fmt.Printf("缓冲区中的日志内容:\n%s", buf.String())

	fmt.Println("\n3. 全局日志（默认约定：控制台彩色+文件JSON+轮转+异步）:")
	// 使用包默认全局日志（InitLogEnv 已在导入时自动调用）
	ygggo_log.Info("全局日志：系统启动", "user=alice", map[string]any{"tries": 3, "pi": 3.14, "ok": true})
	ygggo_log.Warning("全局日志：警告", "module=core")
	ygggo_log.Error("全局日志：错误", "code=E1001")

	fmt.Println("\n4. Panic日志示例 (注意：这会导致程序崩溃):")
	fmt.Println("取消注释下面的代码来测试Panic功能:")
	fmt.Println("// ygggo_log.Panic(\"这是一条panic信息，程序将崩溃\")")

	// 如果要测试panic功能，请取消注释下面这行
	// ygggo_log.Panic("这是一条panic信息，程序将崩溃")

	fmt.Println("\n=== 示例结束 ===")
}
