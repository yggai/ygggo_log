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
	Console    bool      // 是否输出到控制台
	Color      bool      // 是否使用彩色输出
}

// LoadConfigFromEnv 从环境变量加载日志配置
func LoadConfigFromEnv() *LogConfig {
	// 加载环境变量
	ygggo_env.LoadEnv()

	config := &LogConfig{
		Level:      InfoLevel,  // 默认级别为INFO
		OutputFile: "",         // 默认输出到标准输出
		Format:     TextFormat, // 默认格式为文本
		Console:    false,      // 默认不强制输出到控制台
		Color:      false,      // 默认不使用彩色输出
	}

	// 读取日志级别
	levelStr := ygggo_env.GetStr("YGGGO_LOG_LEVEL", "INFO")
	config.Level = parseLogLevel(levelStr)

	// 读取输出文件
	config.OutputFile = ygggo_env.GetStr("YGGGO_LOG_FILE", "")

	// 读取日志格式
	formatStr := ygggo_env.GetStr("YGGGO_LOG_FORMAT", "text")
	config.Format = parseLogFormat(formatStr)

	// 读取控制台输出配置
	consoleStr := ygggo_env.GetStr("YGGGO_LOG_CONSOLE", "false")
	config.Console = parseBool(consoleStr)

	// 读取彩色输出配置
	colorStr := ygggo_env.GetStr("YGGGO_LOG_COLOR", "false")
	config.Color = parseBool(colorStr)

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

	// 根据配置选择格式化器
	if config.Color {
		logger.formatter = NewColorFormatter()
	} else {
		logger.formatter = createFormatter(config.Format)
	}

	return logger
}

// NewLoggerFromConfig 根据配置创建日志记录器
func NewLoggerFromConfig(config *LogConfig) *Logger {
	var output io.Writer = os.Stdout

	// 处理文件输出和控制台输出
	if config.OutputFile != "" && config.Console {
		// 同时输出到文件和控制台
		file, err := os.OpenFile(config.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// 如果创建文件失败，回退到标准输出
			output = os.Stdout
		} else {
			// 创建多重写入器
			output = NewMultiWriter(os.Stdout, file)
		}
	} else if config.OutputFile != "" {
		// 只输出到文件
		file, err := os.OpenFile(config.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			// 如果创建文件失败，回退到标准输出
			output = os.Stdout
		} else {
			output = file
		}
	}
	// 如果没有指定文件，默认输出到标准输出

	logger := NewLogger(output)
	logger.minLevel = config.Level

	// 根据配置选择格式化器
	if config.Color {
		logger.formatter = NewColorFormatter()
	} else {
		logger.formatter = createFormatter(config.Format)
	}

	return logger
}

// parseBool 解析布尔值字符串
func parseBool(boolStr string) bool {
	switch strings.ToLower(boolStr) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return false // 默认返回false
	}
}
