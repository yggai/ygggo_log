# ygggo_log

A pragmatic, convention-over-configuration logging library for Go.

- Global logger out-of-the-box
- Color console + JSON file logging by default
- File rotation by size and count (100MB / 3 files by default)
- Env-based config and thread-safe singleton
- Structured logs, variadic APIs with colorized parameters

Switch language: English | [中文](./README.zh.md)

## Table of Contents
- [ygggo\_log](#ygggo_log)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
  - [Quick Start (Global Logger)](#quick-start-global-logger)
  - [Conventions (Defaults)](#conventions-defaults)
  - [Environment Variables](#environment-variables)
  - [Examples](#examples)
  - [Testing](#testing)
  - [License](#license)

## Features
- Five levels: DEBUG, INFO, WARNING, ERROR, PANIC
- Structured logs: text or JSON
- Colorized parameters with type-aware coloring
- Environment-based configuration and a thread-safe singleton

## Installation
```bash
go get github.com/yggai/ygggo_log
```

## Quick Start (Global Logger)
```go
package main
import gglog "github.com/yggai/ygggo_log"

func main() {
    // Initialized automatically at import via package init()
    gglog.Info("service started", "port=8080", map[string]any{"tries": 3, "ok": true, "pi": 3.14})
    gglog.Warning("slow request", "path=/api")
    gglog.Error("db error", "code=E1001")
}
```
- Console: colored, e.g. `2025-01-01 10:11:12.345 [INFO] main.go:12 message key=value ...`
- File: JSON records under `logs/` with rotation

## Conventions (Defaults)
- Level: INFO
- Console: colored output with time (milliseconds), level, file:line, message, params
- File: enabled by default, path `logs/YYYYMMDD_HHMMSS.log`, JSON format
- File rotation: size 100MB, count 3 files
- High-performance buffering + async for console; file is rotation-safe (synchronous by default for stability)

## Environment Variables
- YGGGO_LOG_LEVEL: DEBUG|INFO|WARNING|ERROR|PANIC (default INFO)
- YGGGO_LOG_FILE: file path (auto-generated under `logs/` when empty)
- YGGGO_LOG_FORMAT: text|json (defaults to text; file uses JSON under conventions)
- YGGGO_LOG_CONSOLE: true|false (console enabled by default under conventions)
- YGGGO_LOG_COLOR: true|false (console colors enabled by default under conventions)
- YGGGO_LOG_FILE_SIZE: e.g. 100M (default 100M)
- YGGGO_LOG_FILE_NUM: integer >=1 (default 3)

## Examples
See `examples/`:
- c01_log: basic global usage (colored console + JSON file)
- c02_env_config: configure via env
- c03_singleton: singleton
- c04_structured_log: text vs JSON
- c05_color_log: color demo

Run:
```bash
go run examples/c01_log/main.go
```

## Testing
```bash
go test -v ./...
```

## License
PolyForm Noncommercial License 1.0.0 — noncommercial use only. Author: 源滚滚 <1156956636@qq.com>.
Issues welcome; PRs not accepted (personal research project).
