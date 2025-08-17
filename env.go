package ygggo_log

import (
	"io"
	"os"
	"strings"

	"github.com/yggai/ygggo_env"
)

// LogConfig 日志配置结构体
type LogConfig struct {
	Level      LogLevel  // 日志级别
	OutputFile string    // 输出文件路径，空字符串表示输出到标准输出
	Format     LogFormat // 日志格式
}

// LoadConfigFromEnv 从环境变量加载日志配置
func LoadConfigFromEnv() *LogConfig {
	// 加载环境变量
	ygggo_env.LoadEnv()

	config := &LogConfig{
		Level:      InfoLevel,  // 默认级别为INFO
		OutputFile: "",         // 默认输出到标准输出
		Format:     TextFormat, // 默认格式为文本
	}

	// 读取日志级别
	levelStr := ygggo_env.GetStr("YGGGO_LOG_LEVEL", "INFO")
	config.Level = parseLogLevel(levelStr)

	// 读取输出文件
	config.OutputFile = ygggo_env.GetStr("YGGGO_LOG_FILE", "")

	// 读取日志格式
	formatStr := ygggo_env.GetStr("YGGGO_LOG_FORMAT", "text")
	config.Format = parseLogFormat(formatStr)

	return config
}

// parseLogLevel 解析日志级别字符串
func parseLogLevel(levelStr string) LogLevel {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "WARNING":
		return WarningLevel
	case "ERROR":
		return ErrorLevel
	case "PANIC":
		return PanicLevel
	default:
		return InfoLevel // 默认返回INFO级别
	}
}

// GetLogEnv 现在在 singleton.go 中实现为单例模式

// NewLoggerFromEnvWithOutput 根据环境变量创建日志记录器，但使用指定的输出
func NewLoggerFromEnvWithOutput(output io.Writer) *Logger {
	config := LoadConfigFromEnv()
	logger := NewLogger(output)
	logger.minLevel = config.Level
	logger.formatter = createFormatter(config.Format)
	return logger
}

// NewLoggerFromConfig 根据配置创建日志记录器
func NewLoggerFromConfig(config *LogConfig) *Logger {
	var output io.Writer = os.Stdout

	// 如果指定了输出文件，则创建文件
	if config.OutputFile != "" {
		file, err := os.OpenFile(config.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// 如果创建文件失败，回退到标准输出
			output = os.Stdout
		} else {
			output = file
		}
	}

	logger := NewLogger(output)
	logger.minLevel = config.Level
	logger.formatter = createFormatter(config.Format)
	return logger
}
