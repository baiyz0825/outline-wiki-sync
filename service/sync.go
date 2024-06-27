// outline-wiki-sync
//
// @(#)sync.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package service

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/baiyz0825/outline-wiki-sync/utils"
	"github.com/patrickmn/go-cache"
)

type SyncMarkDownFile struct {
	FileRootPath []string
	filePathLock *sync.Map
}

func NewSyncMarkDownFile(fileRootPaths []string) *SyncMarkDownFile {
	return &SyncMarkDownFile{
		FileRootPath: fileRootPaths,
		filePathLock: &sync.Map{},
	}
}

func (s *SyncMarkDownFile) getMutexForPath(path string) *sync.Mutex {
	mutex, _ := s.filePathLock.LoadOrStore(path, &sync.Mutex{})
	return mutex.(*sync.Mutex)
}

// SyncMarkdownFile 同步Markdown文件
func (s *SyncMarkDownFile) SyncMarkdownFile() {
	for _, fileRootPath := range s.FileRootPath {
		go func(path string) {
			fileSystem := os.DirFS(path)
			err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
				return s.processOneFileOrPath(path, d, err)
			})
			utils.Log.Errorf("遍历文件夹出错：%v", err)
		}(fileRootPath)
	}
}

// processOneFileOrPath 处理单个文件或文件夹
func (s *SyncMarkDownFile) processOneFileOrPath(path string, d fs.DirEntry, err error) error {
	if err != nil {
		utils.Log.Error("遍历文件路径失败: %v", err)
	}

	if !d.IsDir() {
		go s.processFile(path, d)
	}
	return nil
}

// processDir 处理文件夹 这里对文件夹上锁，有一个处理中或者处理成功就不处理了
func (s *SyncMarkDownFile) processDir(path string) (collectionId string) {
	// 检查缓存
	cacheKey := utils.XCache.GenCollectionCacheKey(path)
	fromCache := utils.XCache.GetDataFromCache(cacheKey)
	if fromCache != nil {
		return fromCache.(string)
	}
	// lock 更新并查看db
	s.getMutexForPath(path).Lock()
	defer s.getMutexForPath(path).Unlock()

	// 检查数据库是否创建了这个Id
	// 获取最后一层文件夹名称 数据库存储全路径映射
	// lastPathName := filepath.Base(path)

	// TODO 创建数据并更新数据库

	// 更新缓存
	utils.XCache.SetDataToCache(cacheKey, collectionId, cache.NoExpiration)
	return ""
}

// processFile 处理文件
func (s *SyncMarkDownFile) processFile(path string, file fs.DirEntry) {
	if filepath.Ext(file.Name()) != ".md" {
		utils.Log.Infof("当前文件不是md文件,跳过处理: %s", file.Name())
		return
	}
	// 打开文件
	// bytes, err := os.ReadFile(path)
	// if err != nil {
	// 	utils.Log.Errorf("读取文件失败: filePath: %s error: %v", path, err)
	// 	return
	// }
	// 获取目录的collection Id
	// collectionId := s.processDir(filepath.Dir(path))

	// 请求接口创建文档

	// 存储数据库文件更新情况
}
