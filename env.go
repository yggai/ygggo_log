package ygggo_log

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/yggai/ygggo_env"
)

// LogConfig 日志配置结构体
type LogConfig struct {
	Level      LogLevel  // 日志级别
	OutputFile string    // 输出文件路径，空字符串表示输出到标准输出
	Format     LogFormat // 日志格式
	Console    bool      // 是否输出到控制台
	Color      bool      // 是否使用彩色输出
	FileSize   int64     // 日志文件大小限制（字节）
	FileNum    int       // 日志文件数量限制
}

// LoadConfigFromEnv 从环境变量加载日志配置
func LoadConfigFromEnv() *LogConfig {
	// 加载环境变量
	ygggo_env.LoadEnv()

	config := &LogConfig{
		Level:      InfoLevel,         // 默认级别为INFO
		OutputFile: "",                // 默认输出到标准输出
		Format:     TextFormat,        // 默认格式为文本
		Console:    false,             // 默认不强制输出到控制台
		Color:      false,             // 默认不使用彩色输出
		FileSize:   100 * 1024 * 1024, // 默认100MB
		FileNum:    3,                 // 默认3个文件
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

	// 读取文件大小配置
	fileSizeStr := ygggo_env.GetStr("YGGGO_LOG_FILE_SIZE", "100M")
	config.FileSize = parseSizeString(fileSizeStr)

	// 读取文件数量配置
	fileNumStr := ygggo_env.GetStr("YGGGO_LOG_FILE_NUM", "3")
	config.FileNum = parseFileNum(fileNumStr)

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
	// 采用约定大于配置：
	// - 默认INFO
	// - 控制台彩色、文件JSON
	// - 文件轮转、大小/数量限制
	// - 异步写入

	// 控制台彩色（异步）
	console := NewAsyncWriter(os.Stdout, 1024)

	// 文件（异步+轮转）
	filePath := config.OutputFile
	if filePath == "" {
		_ = os.MkdirAll("logs", 0755)
		filePath = time.Now().Format("logs/20060102_150405.log")
	}
	rot, _ := NewRotatingWriter(filePath, config.FileSize, config.FileNum)
	var fileOut io.Writer
	if rot != nil {
		fileOut = rot // 同步写入，避免测试环境清理冲突
	}

	// 组合格式化器：控制台彩色 + 文件JSON
	combined := NewCombinedFormatter(console, fileOut)

	logger := NewLogger(io.Discard) // output弃用，由formatter写到目标
	logger.minLevel = config.Level  // 默认由配置决定，默认INFO
	logger.formatter = combined
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

// InitLogEnv 基于环境变量初始化全局日志（defaultLogger）
// 可显式调用；包导入时也会通过 init() 自动调用一次
func InitLogEnv() {
	config := LoadConfigFromEnv()
	defaultLogger = NewLoggerFromConfig(config)
}

func init() {
	// 包导入时自动初始化一次
	InitLogEnv()
}
