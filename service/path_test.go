package service

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestPat(t *testing.T) {

	// 示例路径
	currentPath := "/path/to/dir1"

	// 获取父目录信息
	parentDir := filepath.Dir(currentPath)

	// 输出父目录
	fmt.Println("父目录:", parentDir)
}
