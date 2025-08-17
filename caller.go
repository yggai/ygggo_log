package ygggo_log

import (
	"path/filepath"
	"runtime"
)

func callerInfo() (string, int) {
	// 3 层栈：callerInfo -> ColorFormatter.Format -> Logger.log -> Debug/Info
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "?", 0
	}
	return filepath.Base(file), line
}

