package ygggo_log

import (
	"io"
	"os"
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

// log 内部日志记录方法
func (l *Logger) log(level LogLevel, message string) {
	// 检查日志级别是否满足最小级别要求
	if level < l.minLevel {
		return
	}

	// 使用格式化器格式化日志
	l.formatter.Format(l.output, level, message)
}

// Debug 生成DEBUG级别的日志
func (l *Logger) Debug(message string) {
	l.log(DebugLevel, message)
}

// Info 生成INFO级别的日志
func (l *Logger) Info(message string) {
	l.log(InfoLevel, message)
}

// Warning 生成WARNING级别的日志
func (l *Logger) Warning(message string) {
	l.log(WarningLevel, message)
}

// Error 生成ERROR级别的日志
func (l *Logger) Error(message string) {
	l.log(ErrorLevel, message)
}

// Panic 生成Panic级别的日志并触发panic
func (l *Logger) Panic(message string) {
	l.log(PanicLevel, message)
	panic(message)
}

// 默认日志记录器实例
var defaultLogger = NewLogger(os.Stdout)

// Debug 使用默认日志记录器生成DEBUG级别的日志
func Debug(message string) {
	defaultLogger.Debug(message)
}

// Info 使用默认日志记录器生成INFO级别的日志
func Info(message string) {
	defaultLogger.Info(message)
}

// Warning 使用默认日志记录器生成WARNING级别的日志
func Warning(message string) {
	defaultLogger.Warning(message)
}

// Error 使用默认日志记录器生成ERROR级别的日志
func Error(message string) {
	defaultLogger.Error(message)
}

// Panic 使用默认日志记录器生成Panic级别的日志并触发panic
func Panic(message string) {
	defaultLogger.Panic(message)
}
