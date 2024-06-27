// outline-wiki-sync
//
// @(#)sync.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package service

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/baiyz0825/outline-wiki-sdk"
	"github.com/baiyz0825/outline-wiki-sync/utils"
	cache2 "github.com/baiyz0825/outline-wiki-sync/utils/cache"
	"github.com/baiyz0825/outline-wiki-sync/utils/client"
	"github.com/patrickmn/go-cache"
)

type SyncMarkDownFile struct {
	FileRootPath  []string
	filePathLocks *sync.Map
	ctx           context.Context
}

func NewSyncMarkDownFile(ctx context.Context, fileRootPaths []string) *SyncMarkDownFile {
	return &SyncMarkDownFile{
		FileRootPath:  fileRootPaths,
		filePathLocks: &sync.Map{},
		ctx:           ctx,
	}
}

func (s *SyncMarkDownFile) getMutexForPath(path string) *sync.Mutex {
	mutex, _ := s.filePathLocks.LoadOrStore(path, &sync.Mutex{})
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
		select {
		case <-s.ctx.Done():
			{
				utils.Log.Infof("监测到退出信号, 处理单个文件或文件夹处理业务退出")
				return nil
			}
		default:
			go s.processFile(path, d)
		}
	}
	return nil
}

// processDir 处理文件夹 这里对文件夹上锁，有一个处理中或者处理成功就不处理了
func (s *SyncMarkDownFile) processDir(path string) (collectionId string) {
	// 检查缓存
	cacheKey := cache2.XCache.GenCollectionCacheKey(path)
	fromCache := cache2.XCache.GetDataFromCache(cacheKey)
	if fromCache != nil {
		return fromCache.(string)
	}
	// lock 更新并查看db
	s.getMutexForPath(path).Lock()
	defer s.getMutexForPath(path).Unlock()

	// 检查数据库是否创建了这个Id
	// 获取最后一层文件夹名称 数据库存储全路径映射
	lastPathName := filepath.Base(path)
	request := outline.PostCollectionsCreateJSONRequestBody{
		Description: utils.PtrString(fmt.Sprintf("%s-%s", "sync->", lastPathName)),
		Name:        lastPathName,
		Private:     utils.PtrBool(true),
	}
	ok, response := client.OutlineSdk.CreateCollection(s.ctx, request)
	if !ok {
		utils.Log.Errorf("创建outline文件夹失败: rawPath:%s request:%v response:%v", path, request, response)
		return ""
	}
	utils.Log.Infof("创建outline文件夹成功: rawPath:%s collectionId:%v", path, response)
	// TODO 创建数据并更新数据库

	// 更新缓存
	cache2.XCache.SetDataToCache(cacheKey, collectionId, cache.NoExpiration)
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
