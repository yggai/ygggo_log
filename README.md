# ygggo_log

一个简单易用的Go语言日志库，支持多种日志级别和自定义输出。

## 功能特性

- 支持5种日志级别：DEBUG、INFO、WARNING、ERROR、PANIC
- 支持自定义输出目标（标准输出、文件、缓冲区等）
- 支持环境变量配置（通过 .env 文件）
- **单例模式**：`GetLogEnv()` 始终返回同一个日志对象
- 日志级别过滤功能
- 线程安全的并发支持
- 自动添加时间戳
- 简洁的API设计
- 完整的单元测试覆盖

## 安装

```bash
go get github.com/yggai/ygggo_log
```

## 快速开始

### 使用默认日志记录器

```go
package main

import "github.com/yggai/ygggo_log"

func main() {
    ygggo_log.Debug("这是一条调试信息")
    ygggo_log.Info("这是一条信息")
    ygggo_log.Warning("这是一条警告信息")
    ygggo_log.Error("这是一条错误信息")
    // ygggo_log.Panic("这会触发panic") // 谨慎使用
}
```

### 使用环境变量配置

创建 `.env` 文件：
```env
# 日志配置环境变量
YGGGO_LOG_LEVEL=DEBUG
YGGGO_LOG_FILE=app.log
```

使用环境变量配置的日志记录器：
```go
package main

import "github.com/yggai/ygggo_log"

func main() {
    // 自动从 .env 文件加载配置
    logger := ygggo_log.GetLogEnv()

    logger.Debug("这是DEBUG信息")
    logger.Info("这是INFO信息")
    logger.Warning("这是WARNING信息")
    logger.Error("这是ERROR信息")
}
```

### 使用自定义日志记录器

```go
package main

import (
    "os"
    "github.com/yggai/ygggo_log"
)

func main() {
    // 输出到文件
    file, _ := os.Create("app.log")
    defer file.Close()

    logger := ygggo_log.NewLogger(file)
    logger.Info("这条日志将写入文件")
}
```

## API 文档

### 日志级别

- `Debug(message string)` - 调试级别日志
- `Info(message string)` - 信息级别日志
- `Warning(message string)` - 警告级别日志
- `Error(message string)` - 错误级别日志
- `Panic(message string)` - Panic级别日志（会触发panic）

### 环境变量配置与单例模式

- `GetLogEnv() *Logger` - 获取单例日志记录器（基于环境变量配置）
- `LoadConfigFromEnv() *LogConfig` - 从环境变量加载配置

**重要特性**：
- `GetLogEnv()` 实现了单例模式，无论在代码的任何地方调用，都会返回同一个日志对象
- 单例模式是线程安全的，支持并发调用
- 配置只在第一次调用时加载，后续调用不会重新读取环境变量

支持的环境变量：
- `YGGGO_LOG_LEVEL` - 日志级别（DEBUG、INFO、WARNING、ERROR、PANIC）
- `YGGGO_LOG_FILE` - 输出文件路径（空值表示输出到标准输出）

### 自定义日志记录器

- `NewLogger(output io.Writer) *Logger` - 创建新的日志记录器

## 示例

查看 `examples/` 目录中的完整示例：

- `c01_log/main.go` - 基本日志功能示例
- `c02_env_config/main.go` - 环境变量配置示例
- `c03_singleton/main.go` - 单例模式示例

运行示例：

```bash
go run examples/c01_log/main.go
go run examples/c02_env_config/main.go
go run examples/c03_singleton/main.go
```

## 测试

运行所有测试：

```bash
go test -v
```

## 日志格式

日志输出格式为：`时间戳 [级别] 消息`

示例：
```
2025-08-17 20:02:00 [INFO] 这是一条信息
2025-08-17 20:02:00 [ERROR] 这是一条错误信息
```
Go语言用于记录日志的底层核心库
