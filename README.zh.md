# ygggo_log

一个遵循“约定大于配置”的 Go 日志库：
- 开箱即用的全局日志
- 控制台彩色 + 文件 JSON 同时输出（默认）
- 日志文件大小/数量轮转
- 环境变量配置与线程安全单例
- 结构化日志、可变参数、参数彩色高亮

语言切换：中文 | [English](./README.md)

## 功能特性
- 五种日志级别：DEBUG、INFO、WARNING、ERROR、PANIC
- 约定优于配置的默认：
  - 级别：INFO
  - 控制台：彩色输出，显示 时间(毫秒)、级别、文件:行号、消息、参数
  - 文件：默认开启；路径 logs/YYYYMMDD_HHMMSS.log；JSON 格式
  - 文件大小：100MB；文件个数：3（轮转）
  - 控制台采用异步缓冲写入（文件端为轮转安全，默认同步）
- 结构化日志：文本/JSON
- 参数彩色高亮（根据类型着色）
- 环境变量配置 + 线程安全单例
- 完整单元测试覆盖

## 安装
```bash
go get github.com/yggai/ygggo_log
```

## 快速开始（全局）
```go
package main
import gglog "github.com/yggai/ygggo_log"

func main() {
    // 包导入时自动初始化
    gglog.Info("服务启动", "port=8080", map[string]any{"tries": 3, "ok": true, "pi": 3.14})
    gglog.Warning("请求较慢", "path=/api")
    gglog.Error("数据库错误", "code=E1001")
}
```
- 控制台输出（彩色）：`2025-01-01 10:11:12.345 [INFO] main.go:12 message key=value ...`
- 文件输出（JSON）：默认写入 logs/ 下并按大小与数量轮转

## 环境变量
- YGGGO_LOG_LEVEL: DEBUG|INFO|WARNING|ERROR|PANIC（默认 INFO）
- YGGGO_LOG_FILE: 文件路径（为空时自动生成 logs/xxx.log）
- YGGGO_LOG_FORMAT: text|json（默认 text；但文件默认 JSON）
- YGGGO_LOG_CONSOLE: true|false（约定下控制台默认开启）
- YGGGO_LOG_COLOR: true|false（约定下控制台默认彩色）
- YGGGO_LOG_FILE_SIZE: 如 100M（默认 100M）
- YGGGO_LOG_FILE_NUM: >=1（默认 3）

## 示例
- c01_log：全局日志（彩色控制台 + JSON 文件）
- c02_env_config：环境变量配置
- c03_singleton：单例
- c04_structured_log：文本与 JSON
- c05_color_log：彩色日志

运行：
```bash
go run examples/c01_log/main.go
```

## 测试
```bash
go test -v ./...
```

## 许可证
本项目采用 PolyForm Noncommercial License 1.0.0，仅允许非商业用途。
作者：源滚滚 <1156956636@qq.com>。本项目是个人研究项目，不接受代码合并（PR），但欢迎提 issue。

