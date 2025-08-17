package main

import (
	gglog "github.com/yggai/ygggo_log"
)

func main() {
	// 可显式初始化（实际包导入时 init() 会自动调用一次）
	gglog.InitLogEnv()

	gglog.Info("Global logger initialized via environment variables")
	gglog.Debug("Debug message")
	gglog.Warning("Warning message")
	gglog.Error("Error message")
	// 注意：Panic 会触发 panic，请按需测试
	// gglog.Panic("Panic message")
}

