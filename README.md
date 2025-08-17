# ygggo_log

一个简单易用的Go语言日志库，支持多种日志级别和自定义输出。

## 功能特性

- 支持5种日志级别：DEBUG、INFO、WARNING、ERROR、PANIC
- 支持自定义输出目标（标准输出、文件、缓冲区等）
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

### 自定义日志记录器

- `NewLogger(output io.Writer) *Logger` - 创建新的日志记录器

## 示例

查看 `examples/c01_log/` 目录中的完整示例：

- `main.go` - 基本使用示例
- `panic_example.go` - Panic功能示例

运行示例：

```bash
go run examples/c01_log/main.go
go run examples/c01_log/panic_example.go
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
