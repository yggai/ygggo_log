package ygggo_log

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// LogFormat 定义日志格式类型
type LogFormat int

const (
	TextFormat LogFormat = iota // 文本格式（默认）
	JsonFormat                  // JSON格式
)

// String 返回日志格式的字符串表示
func (f LogFormat) String() string {
	switch f {
	case TextFormat:
		return "text"
	case JsonFormat:
		return "json"
	default:
		return "text"
	}
}

// Formatter 日志格式化器接口
type Formatter interface {
	Format(writer io.Writer, level LogLevel, message string)
}

// TextFormatter 文本格式化器
type TextFormatter struct{}

// NewTextFormatter 创建文本格式化器
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

// Format 格式化为文本格式
func (f *TextFormatter) Format(writer io.Writer, level LogLevel, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("%s [%s] %s\n", timestamp, level.String(), message)
	writer.Write([]byte(logEntry))
}

// JsonFormatter JSON格式化器
type JsonFormatter struct{}

// NewJsonFormatter 创建JSON格式化器
func NewJsonFormatter() *JsonFormatter {
	return &JsonFormatter{}
}

// JsonLogEntry JSON日志条目结构
type JsonLogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// Format 格式化为JSON格式
func (f *JsonFormatter) Format(writer io.Writer, level LogLevel, message string) {
	entry := JsonLogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level.String(),
		Message:   message,
	}
	
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// 如果JSON序列化失败，回退到文本格式
		textFormatter := NewTextFormatter()
		textFormatter.Format(writer, level, message)
		return
	}
	
	writer.Write(jsonData)
	writer.Write([]byte("\n"))
}

// parseLogFormat 解析日志格式字符串
func parseLogFormat(formatStr string) LogFormat {
	switch formatStr {
	case "json":
		return JsonFormat
	case "text":
		return TextFormat
	default:
		return TextFormat // 默认返回文本格式
	}
}

// createFormatter 根据格式创建对应的格式化器
func createFormatter(format LogFormat) Formatter {
	switch format {
	case JsonFormat:
		return NewJsonFormatter()
	case TextFormat:
		return NewTextFormatter()
	default:
		return NewTextFormatter()
	}
}
