package ygggo_log

import "io"

// BridgeFormatter 允许在同一次日志中，同时将文本消息写到两个不同的 writer
// 用于默认约定：控制台彩色、文件JSON
// 注意：它仅负责将已经格式化好的字符串写出，不负责再格式化

type BridgeFormatter struct {
	console io.Writer
	file    io.Writer
}

func NewBridgeFormatter(console io.Writer, file io.Writer) *BridgeFormatter {
	return &BridgeFormatter{console: console, file: file}
}

func (bf *BridgeFormatter) WriteConsole(p []byte) (int, error) {
	if bf.console == nil { return 0, nil }
	return bf.console.Write(p)
}
func (bf *BridgeFormatter) WriteFile(p []byte) (int, error) {
	if bf.file == nil { return 0, nil }
	return bf.file.Write(p)
}

