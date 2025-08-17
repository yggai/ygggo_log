package ygggo_log

import (
	"io"
)

// AsyncWriter 简单的异步写入器，使用channel缓冲
// 注意：无关闭与flush机制，适合作为默认异步优化；需要严谨生命周期的场景可后续扩展
type AsyncWriter struct {
	w   io.Writer
	ch  chan []byte
}

func NewAsyncWriter(w io.Writer, bufSize int) *AsyncWriter {
	aw := &AsyncWriter{
		w:  w,
		ch: make(chan []byte, bufSize),
	}
	go aw.loop()
	return aw
}

func (aw *AsyncWriter) Write(p []byte) (int, error) {
	// 拷贝避免调用方复用切片带来的数据竞争
	cp := make([]byte, len(p))
	copy(cp, p)
	aw.ch <- cp
	return len(p), nil
}

func (aw *AsyncWriter) loop() {
	for b := range aw.ch {
		_, _ = aw.w.Write(b)
	}
}

