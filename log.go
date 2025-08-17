package ygggo_log

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// LogLevel 定义日志级别
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	PanicLevel
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarningLevel:
		return "WARNING"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	default:
		return "UNKNOWN"
	}
}

// Logger 日志记录器结构体
type Logger struct {
	output    io.Writer
	minLevel  LogLevel  // 最小日志级别，低于此级别的日志将被过滤
	formatter Formatter // 日志格式化器
}

// NewLogger 创建一个新的日志记录器
func NewLogger(output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}
	return &Logger{
		output:    output,
		minLevel:  DebugLevel,         // 默认显示所有级别的日志
		formatter: NewTextFormatter(), // 默认使用文本格式化器
	}
}

// log 内部日志记录方法（支持可选参数）
func (l *Logger) log(level LogLevel, message string, args ...any) {
	// 检查日志级别是否满足最小级别要求
	if level < l.minLevel {
		return
	}

	// 构建包含参数的消息
	fullMsg := l.buildMessage(message, args...)

	// 使用格式化器格式化日志
	l.formatter.Format(l.output, level, fullMsg)
}

// buildMessage 将 message 与可选参数拼接为一个字符串
// 支持的参数形式：
// - 字符串（直接追加）
// - 形如 "key=value" 的字符串
// - map[string]any
// - 任意值（按 v1 v2 v3 顺序追加）
func (l *Logger) buildMessage(message string, args ...any) string {
	if len(args) == 0 {
		return message
	}
	// 简易拼接：message + " " + params
	var params string
	switch l.formatter.(type) {
	case *ColorFormatter:
		params = formatParamsColor(args...)
	default:
		params = formatParamsPlain(args...)
	}
	if params == "" {
		return message
	}
	return message + " " + params
}

// formatParamsPlain 将参数格式化为不带颜色的 key=value 串
func formatParamsPlain(args ...any) string {
	if len(args) == 0 {
		return ""
	}
	var b strings.Builder
	first := true
	for _, a := range args {
		if a == nil {
			continue
		}
		switch v := a.(type) {
		case map[string]any:
			for k, val := range v {
				if !first {
					b.WriteString(" ")
				}
				first = false
				b.WriteString(k)
				b.WriteString("=")
				b.WriteString(fmt.Sprintf("%v", val))
			}
		case string:
			if !first {
				b.WriteString(" ")
			}
			first = false
			b.WriteString(v)
		default:
			if !first {
				b.WriteString(" ")
			}
			first = false
			b.WriteString(fmt.Sprintf("%v", v))
		}
	}
	return b.String()
}

// formatParamsColor 将参数格式化为带颜色的 key=value 串
func formatParamsColor(args ...any) string {
	if len(args) == 0 {
		return ""
	}
	var b strings.Builder
	first := true
	for _, a := range args {
		if a == nil {
			continue
		}
		switch v := a.(type) {
		case map[string]any:
			for k, val := range v {
				if !first {
					b.WriteString(" ")
				}
				first = false
				b.WriteString(ColorCyan)
				b.WriteString(k)
				b.WriteString(ColorReset)
				b.WriteString("=")
				b.WriteString(colorizeValue(val))
			}
		case string:
			if !first {
				b.WriteString(" ")
			}
			first = false
			parts := strings.SplitN(v, "=", 2)
			if len(parts) == 2 {
				b.WriteString(ColorCyan)
				b.WriteString(parts[0])
				b.WriteString(ColorReset)
				b.WriteString("=")
				b.WriteString(colorizeValue(parts[1]))
			} else {
				b.WriteString(v)
			}
		default:
			if !first {
				b.WriteString(" ")
			}
			first = false
			b.WriteString(colorizeValue(v))
		}
	}
	return b.String()
}

// colorizeValue 根据类型为值着色（用于彩色输出）
func colorizeValue(v any) string {
	switch val := v.(type) {
	case bool:
		return ColorYellow + fmt.Sprintf("%v", val) + ColorReset
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return ColorGreen + fmt.Sprintf("%v", val) + ColorReset
	case float32, float64:
		return ColorPurple + fmt.Sprintf("%v", val) + ColorReset
	case string:
		return ColorWhite + val + ColorReset
	default:
		return ColorBlue + fmt.Sprintf("%v", val) + ColorReset
	}
}

// Debug 生成DEBUG级别的日志（支持参数）
func (l *Logger) Debug(message string, args ...any) {
	l.log(DebugLevel, message, args...)
}

// Info 生成INFO级别的日志（支持参数）
func (l *Logger) Info(message string, args ...any) {
	l.log(InfoLevel, message, args...)
}

// Warning 生成WARNING级别的日志（支持参数）
func (l *Logger) Warning(message string, args ...any) {
	l.log(WarningLevel, message, args...)
}

// Error 生成ERROR级别的日志（支持参数）
func (l *Logger) Error(message string, args ...any) {
	l.log(ErrorLevel, message, args...)
}

// Panic 生成Panic级别的日志并触发panic（支持参数）
func (l *Logger) Panic(message string, args ...any) {
	l.log(PanicLevel, message, args...)
	panic(message)
}

// 默认日志记录器实例
var defaultLogger = NewLogger(os.Stdout)

// Debug 使用默认日志记录器生成DEBUG级别的日志（支持参数）
func Debug(message string, args ...any) {
	defaultLogger.Debug(message, args...)
}

// Info 使用默认日志记录器生成INFO级别的日志（支持参数）
func Info(message string, args ...any) {
	defaultLogger.Info(message, args...)
}

// Warning 使用默认日志记录器生成WARNING级别的日志（支持参数）
func Warning(message string, args ...any) {
	defaultLogger.Warning(message, args...)
}

// Error 使用默认日志记录器生成ERROR级别的日志（支持参数）
func Error(message string, args ...any) {
	defaultLogger.Error(message, args...)
}

// Panic 使用默认日志记录器生成Panic级别的日志并触发panic（支持参数）
func Panic(message string, args ...any) {
	defaultLogger.Panic(message, args...)
}
