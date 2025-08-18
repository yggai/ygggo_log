package ygggo_log

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// LogLevel represents severity for log records in ascending order.
// Values are ordered so that higher levels represent more severe events.
type LogLevel int

const (
	// DebugLevel is used for verbose diagnostic information to help troubleshooting.
	DebugLevel LogLevel = iota
	// InfoLevel is used for routine information, startup messages, progress, etc.
	InfoLevel
	// WarningLevel indicates something unexpected happened, but the application continues.
	WarningLevel
	// ErrorLevel indicates an error occurred that prevented an operation from succeeding.
	ErrorLevel
	// PanicLevel logs the message and triggers a panic.
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

// Logger is a minimal, pluggable logger with level filtering and a formatter.
// It is concurrency-safe as long as the configured output is safe for concurrent writes.
type Logger struct {
	output    io.Writer
	minLevel  LogLevel  // Minimum level to emit; messages below are discarded.
	formatter Formatter // Responsible for rendering a log entry to the output.
}

// NewLogger creates a new Logger that writes to the provided output.
// If output is nil, os.Stdout is used. By default, the logger prints all levels
// using a TextFormatter.
func NewLogger(output io.Writer) *Logger {
	if output == nil {
		output = os.Stdout
	}
	return &Logger{
		output:    output,
		minLevel:  DebugLevel,         // default: emit all levels
		formatter: NewTextFormatter(), // default: text formatter
	}
}

// log writes a log entry at the given level after level filtering. It accepts
// variadic arguments and appends them to the message as key=value pairs.
func (l *Logger) log(level LogLevel, message string, args ...any) {
	if level < l.minLevel {
		return
	}
	fullMsg := l.buildMessage(message, args...)
	l.formatter.Format(l.output, level, fullMsg)
}

// buildMessage joins message with formatted parameters. Supported argument forms:
//   - string (appended directly)
//   - "key=value" strings
//   - map[string]any
//   - any other value, appended in order
//
// When ColorFormatter is in use, keys and values are colorized.
func (l *Logger) buildMessage(message string, args ...any) string {
	if len(args) == 0 {
		return message
	}
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
