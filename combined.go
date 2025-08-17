package ygggo_log

import "io"

// CombinedFormatter 同时将日志写到控制台（彩色）和文件（JSON）
type CombinedFormatter struct {
	console io.Writer
	file    io.Writer
}

func NewCombinedFormatter(console io.Writer, file io.Writer) *CombinedFormatter {
	return &CombinedFormatter{console: console, file: file}
}

func (f *CombinedFormatter) Format(_ io.Writer, level LogLevel, message string) {
	// 控制台：彩色
	if f.console != nil {
		cf := NewColorFormatter()
		cf.Format(f.console, level, message)
	}
	// 文件：JSON
	if f.file != nil {
		jf := NewJsonFormatter()
		jf.Format(f.file, level, message)
	}
}

