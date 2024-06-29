package fileutils

import (
	"os"

	"github.com/baiyz0825/outline-wiki-sync/utils/xlog"
)

func CheckIsEmptyDir(absPath string) bool {
	// 打开当前文件夹
	dir, err := os.Open(absPath)
	if err != nil {
		return true
	}
	defer func(dir *os.File) {
		_ = dir.Close()
	}(dir)

	// 读取文件夹内容的名称
	names, err := dir.Readdirnames(0) // 0 表示不限制返回的名称数
	if err != nil {
		xlog.Log.Errorf("读取文件夹内容失败: %v\n", err)
		return true
	}

	// 检查文件夹是否为空
	if len(names) == 0 {
		return true
	}
	return false
}
