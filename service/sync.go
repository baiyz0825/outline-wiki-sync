// outline-wiki-sync
//
// @(#)sync.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package service

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/baiyz0825/outline-wiki-sync/utils"
)

type SyncMarkDownFile struct {
	FileRootPath []string
}

func NewSyncMarkDownFile(fileRootPaths []string) *SyncMarkDownFile {
	return &SyncMarkDownFile{
		FileRootPath: fileRootPaths,
	}
}

// SyncMarkdownFile 同步Markdown文件
func (s *SyncMarkDownFile) SyncMarkdownFile() {
	for _, fileRootPath := range s.FileRootPath {
		go func(path string) {
			fileSystem := os.DirFS(path)
			err := fs.WalkDir(fileSystem, ".", processOneFileOrPath)
			utils.Log.Errorf("遍历文件夹出错：%v", err)
		}(fileRootPath)
	}
}

// processOneFileOrPath 处理单个文件或文件夹
func processOneFileOrPath(path string, d fs.DirEntry, err error) error {
	if err != nil {
		utils.Log.Error("遍历文件路径失败: %v", err)
	}
	if d.IsDir() {
		processDir(path, d)
	} else {
		go processFile(path, d)
	}
	return nil
}

// processDir 处理文件夹 这里对文件夹上锁，有一个处理中或者处理成功就不处理了
func processDir(path string, dir fs.DirEntry) {

}

// processFile 处理文件
func processFile(path string, file fs.DirEntry) {
	if filepath.Ext(file.Name()) != ".md" {
		utils.Log.Infof("当前文件不是md文件,跳过处理: %s", file.Name())
		return
	}

	// 打开文件

	// 获取目录的collection Id

	// 存储数据库遍历结果

	// 请求接口创建文档
}
