package ygggo_log

import (
	"fmt"
	"io"
	"time"
)

// ANSI颜色代码常量
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// ColorFormatter 彩色格式化器
type ColorFormatter struct{}

// NewColorFormatter 创建彩色格式化器
func NewColorFormatter() *ColorFormatter {
	return &ColorFormatter{}
}

// Format 格式化为彩色文本格式
func (f *ColorFormatter) Format(writer io.Writer, level LogLevel, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	colorCode := getColorCode(level)
	resetCode := getResetCode()

	logEntry := fmt.Sprintf("%s%s [%s] %s%s\n",
		colorCode, timestamp, level.String(), message, resetCode)
	writer.Write([]byte(logEntry))
}

// getColorCode 根据日志级别获取颜色代码
func getColorCode(level LogLevel) string {
	switch level {
	case DebugLevel:
		return ColorCyan // 青色
	case InfoLevel:
		return ColorGreen // 绿色
	case WarningLevel:
		return ColorYellow // 黄色
	case ErrorLevel:
		return ColorRed // 红色
	case PanicLevel:
		return ColorPurple // 紫色
	default:
		return ColorWhite // 白色
	}
}

// getResetCode 获取重置颜色代码
func getResetCode() string {
	return ColorReset
}

// ColorAwareMultiWriter 颜色感知的多重写入器
// 对控制台输出使用彩色格式化器，对文件输出使用普通格式化器
type ColorAwareMultiWriter struct {
	consoleWriter io.Writer
	fileWriter    io.Writer
	useColor      bool
}

// NewColorAwareMultiWriter 创建颜色感知的多重写入器
func NewColorAwareMultiWriter(consoleWriter, fileWriter io.Writer, useColor bool) *ColorAwareMultiWriter {
	return &ColorAwareMultiWriter{
		consoleWriter: consoleWriter,
		fileWriter:    fileWriter,
		useColor:      useColor,
	}
}

// WriteWithFormatter 使用指定的格式化器写入
func (cw *ColorAwareMultiWriter) WriteWithFormatter(level LogLevel, message string) {
	// 写入控制台（可能带颜色）
	if cw.consoleWriter != nil {
		if cw.useColor {
			colorFormatter := NewColorFormatter()
			colorFormatter.Format(cw.consoleWriter, level, message)
		} else {
			textFormatter := NewTextFormatter()
			textFormatter.Format(cw.consoleWriter, level, message)
		}
	}

	// 写入文件（不带颜色）
	if cw.fileWriter != nil {
		textFormatter := NewTextFormatter()
		textFormatter.Format(cw.fileWriter, level, message)
	}
}
