package ygggo_log

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/yggai/ygggo_env"
)

// LogConfig holds logger configuration values loaded from environment variables.
// See LoadConfigFromEnv for defaults and supported variables.
type LogConfig struct {
	Level      LogLevel  // minimum log level
	OutputFile string    // output file path; empty means stdout only
	Format     LogFormat // text or json for non-colored outputs
	Console    bool      // force console output
	Color      bool      // color output for console
	FileSize   int64     // max file size in bytes (rotation)
	FileNum    int       // max number of files (rotation)
}

// LoadConfigFromEnv loads configuration from environment variables, applying
// sensible defaults. It also ensures the .env file is read, if present.
// Defaults:
//   - Level: INFO
//   - OutputFile: "" (stdout only; conventions may choose a default file path)
//   - Format: text
//   - Console: false
//   - Color: false
//   - FileSize: 100MB
//   - FileNum: 3
func LoadConfigFromEnv() *LogConfig {
	// Load .env and OS environment
	ygggo_env.LoadEnv()

	config := &LogConfig{
		Level:      InfoLevel,
		OutputFile: "",
		Format:     TextFormat,
		Console:    false,
		Color:      false,
		FileSize:   100 * 1024 * 1024,
		FileNum:    3,
	}

	// Level
	levelStr := ygggo_env.GetStr("YGGGO_LOG_LEVEL", "INFO")
	config.Level = parseLogLevel(levelStr)

	// Output file
	config.OutputFile = ygggo_env.GetStr("YGGGO_LOG_FILE", "")

	// Format
	formatStr := ygggo_env.GetStr("YGGGO_LOG_FORMAT", "text")
	config.Format = parseLogFormat(formatStr)

	// Console
	consoleStr := ygggo_env.GetStr("YGGGO_LOG_CONSOLE", "false")
	config.Console = parseBool(consoleStr)

	// Color
	colorStr := ygggo_env.GetStr("YGGGO_LOG_COLOR", "false")
	config.Color = parseBool(colorStr)

	// File size
	fileSizeStr := ygggo_env.GetStr("YGGGO_LOG_FILE_SIZE", "100M")
	config.FileSize = parseSizeString(fileSizeStr)

	// File num
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

// NewLoggerFromEnvWithOutput creates a Logger using environment settings but
// writing to the provided output. Useful for tests or custom sinks.
func NewLoggerFromEnvWithOutput(output io.Writer) *Logger {
	config := LoadConfigFromEnv()
	logger := NewLogger(output)
	logger.minLevel = config.Level

	if config.Color {
		logger.formatter = NewColorFormatter()
	} else {
		logger.formatter = createFormatter(config.Format)
	}
	return logger
}

// NewLoggerFromConfig creates a Logger that follows convention-over-configuration
// defaults: colored console + JSON file with rotation, INFO level, and async console
// writes. File path defaults to logs/YYYYMMDD_HHMMSS.log when not provided.
func NewLoggerFromConfig(config *LogConfig) *Logger {
	// Console: colored + async buffering
	console := NewAsyncWriter(os.Stdout, 1024)

	// File: rotation (size/count), default path under logs/
	filePath := config.OutputFile
	if filePath == "" {
		_ = os.MkdirAll("logs", 0755)
		filePath = time.Now().Format("logs/20060102_150405.log")
	}
	rot, _ := NewRotatingWriter(filePath, config.FileSize, config.FileNum)
	var fileOut io.Writer
	if rot != nil {
		fileOut = rot // sync by default; stable for tests & Windows
	}

	// Combined formatter: console (color) + file (JSON)
	combined := NewCombinedFormatter(console, fileOut)

	logger := NewLogger(io.Discard) // formatter writes to destinations
	logger.minLevel = config.Level
	logger.formatter = combined
	return logger
}

// parseBool parses a boolean-ish string into a bool.
func parseBool(boolStr string) bool {
	switch strings.ToLower(boolStr) {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return false
	}
}

// InitLogEnv initializes the package-level defaultLogger from environment
// variables and conventions. It runs automatically on package import, and can
// also be called explicitly by applications.
func InitLogEnv() {
	config := LoadConfigFromEnv()
	defaultLogger = NewLoggerFromConfig(config)
}

func init() {
	InitLogEnv()
}
