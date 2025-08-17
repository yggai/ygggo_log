package ygggo_log

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// RotatingWriter 轮转写入器
type RotatingWriter struct {
	filename    string     // 日志文件名
	maxSize     int64      // 最大文件大小（字节）
	maxFiles    int        // 最大文件数量（包含当前文件）
	currentSize int64      // 当前文件大小
	file        *os.File   // 当前文件句柄
	mutex       sync.Mutex // 并发保护
}

// NewRotatingWriter 创建轮转写入器
func NewRotatingWriter(filename string, maxSize int64, maxFiles int) (*RotatingWriter, error) {
	rw := &RotatingWriter{
		filename: filename,
		maxSize:  maxSize,
		maxFiles: maxFiles,
	}

	// 打开或创建文件
	err := rw.openFile()
	if err != nil {
		return nil, err
	}

	return rw, nil
}

// openFile 打开或创建日志文件
func (rw *RotatingWriter) openFile() error {
	// 获取文件信息
	if stat, err := os.Stat(rw.filename); err == nil {
		rw.currentSize = stat.Size()
	} else {
		rw.currentSize = 0
	}

	// 打开文件
	file, err := os.OpenFile(rw.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	rw.file = file
	return nil
}

// Write 实现io.Writer接口
func (rw *RotatingWriter) Write(p []byte) (n int, err error) {
	rw.mutex.Lock()
	defer rw.mutex.Unlock()

	// 检查是否需要轮转
	if rw.maxSize > 0 && rw.currentSize+int64(len(p)) > rw.maxSize {
		if err = rw.rotate(); err != nil {
			return 0, err
		}
	}

	// 如果文件未打开（例如上一次写入后关闭了），则重新打开
	if rw.file == nil {
		if err = rw.openFile(); err != nil {
			return 0, err
		}
	}

	// 写入数据
	n, err = rw.file.Write(p)
	if err != nil {
		return n, err
	}

	// 写入完成后立即关闭文件，避免长时间占用句柄（兼容测试在 Windows 上的清理）
	_ = rw.file.Close()
	rw.file = nil

	rw.currentSize += int64(n)
	return n, nil
}

// rotate 执行文件轮转
func (rw *RotatingWriter) rotate() error {
	// 关闭当前文件
	if rw.file != nil {
		_ = rw.file.Close()
	}

	// 轮转文件
	err := rw.rotateFiles()
	if err != nil {
		return err
	}

	// 重新打开文件
	return rw.openFile()
}

// rotateFiles 轮转文件名
func (rw *RotatingWriter) rotateFiles() error {
	// 获取现有的轮转文件
	rotatedFiles, err := rw.getRotatedFiles()
	if err != nil {
		return err
	}

	// 删除超出数量限制的文件（总数不超过 rw.maxFiles，包括当前文件）
	if rw.maxFiles > 1 && len(rotatedFiles) >= rw.maxFiles-1 {
		// 保留最新的 (maxFiles - 2) 个轮转文件
		filesToDelete := len(rotatedFiles) - (rw.maxFiles - 2)
		for i := 0; i < filesToDelete; i++ {
			_ = os.Remove(rotatedFiles[i].path)
		}
		rotatedFiles = rotatedFiles[filesToDelete:]
	}

	// 重命名现有文件（索引+1）
	for i := len(rotatedFiles) - 1; i >= 0; i-- {
		oldPath := rotatedFiles[i].path
		newPath := fmt.Sprintf("%s.%d", rw.filename, rotatedFiles[i].index+1)
		_ = os.Rename(oldPath, newPath)
	}

	// 将当前文件重命名为 .1
	if _, err := os.Stat(rw.filename); err == nil {
		newPath := fmt.Sprintf("%s.1", rw.filename)
		if err := os.Rename(rw.filename, newPath); err != nil {
			return err
		}
	}

	return nil
}

// rotatedFile 轮转文件信息
type rotatedFile struct {
	path  string
	index int
}

// getRotatedFiles 获取现有的轮转文件
func (rw *RotatingWriter) getRotatedFiles() ([]rotatedFile, error) {
	dir := filepath.Dir(rw.filename)
	base := filepath.Base(rw.filename)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var rotatedFiles []rotatedFile
	pattern := regexp.MustCompile(fmt.Sprintf(`^%s\.(\d+)$`, regexp.QuoteMeta(base)))

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := pattern.FindStringSubmatch(file.Name())
		if len(matches) == 2 {
			index, err := strconv.Atoi(matches[1])
			if err != nil {
				continue
			}

			rotatedFiles = append(rotatedFiles, rotatedFile{
				path:  filepath.Join(dir, file.Name()),
				index: index,
			})
		}
	}

	// 按索引排序（从小到大）
	sort.Slice(rotatedFiles, func(i, j int) bool {
		return rotatedFiles[i].index < rotatedFiles[j].index
	})

	return rotatedFiles, nil
}

// Close 关闭写入器
func (rw *RotatingWriter) Close() error {
	rw.mutex.Lock()
	defer rw.mutex.Unlock()

	if rw.file != nil {
		return rw.file.Close()
	}
	return nil
}

// parseSizeString 解析大小字符串
func parseSizeString(sizeStr string) int64 {
	sizeStr = strings.TrimSpace(strings.ToUpper(sizeStr))

	// 匹配数字和单位
	re := regexp.MustCompile(`^(\d+)([KMGT]?B?)$`)
	matches := re.FindStringSubmatch(sizeStr)

	if len(matches) < 2 {
		// 尝试解析纯数字
		if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
			return size
		}
		// 解析失败，返回默认值100MB
		return 100 * 1024 * 1024
	}

	size, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 100 * 1024 * 1024 // 默认100MB
	}

	unit := matches[2]
	switch unit {
	case "K", "KB":
		return size * 1024
	case "M", "MB":
		return size * 1024 * 1024
	case "G", "GB":
		return size * 1024 * 1024 * 1024
	case "T", "TB":
		return size * 1024 * 1024 * 1024 * 1024
	default:
		return size // 无单位，按字节处理
	}
}

// parseFileNum 解析文件数量字符串（>=1），否则返回默认3
func parseFileNum(numStr string) int {
	numStr = strings.TrimSpace(numStr)
	n, err := strconv.Atoi(numStr)
	if err != nil || n < 1 {
		return 3
	}
	return n
}
