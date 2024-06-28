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
	"time"

	"github.com/baiyz0825/outline-wiki-sdk"
	"github.com/baiyz0825/outline-wiki-sync/dao"
	"github.com/baiyz0825/outline-wiki-sync/model"
	"github.com/baiyz0825/outline-wiki-sync/utils"
	cache2 "github.com/baiyz0825/outline-wiki-sync/utils/cache"
	"github.com/baiyz0825/outline-wiki-sync/utils/client"
	"github.com/baiyz0825/outline-wiki-sync/utils/xlog"
	"github.com/google/uuid"
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
	var rootWg sync.WaitGroup
	for _, fileRootPath := range s.FileRootPath {
		collectionId := ""
		// // init collectionId
		exist, err := dao.OutlineWikiCollectionMapping.WithContext(s.ctx).
			Where(dao.OutlineWikiCollectionMapping.CollectionPath.Eq(fileRootPath), dao.OutlineWikiCollectionMapping.RealCollection.Is(true)).First()
		if err == nil && exist != nil {
			collectionId = exist.CollectionId
		} else {
			base := filepath.Base(fileRootPath)
			request := outline.PostCollectionsCreateJSONRequestBody{
				Description: utils.PtrString(fmt.Sprintf("%s-%s", "sync->", base)),
				Name:        base,
				Private:     utils.PtrBool(true),
			}
			ok, response := client.OutlineSdk.CreateCollection(s.ctx, request)
			if !ok {
				xlog.Log.Errorf("创建outline文件夹失败: rawPath:%s request:%v response:%v", fileRootPath, request, response)
				continue
			}
			xlog.Log.Infof("创建outline文件夹成功: rawPath:%s collectionId:%v", fileRootPath, response)
			mapping := &model.OutlineWikiCollectionMapping{
				CollectionId:   response.JSON200.Data.Id.String(),
				CollectionPath: fileRootPath,
				CollectionName: base,
				RealCollection: true,
				Sync:           true,
				CreatedAt:      time.Time{},
				UpdatedAt:      time.Time{},
			}
			err := dao.OutlineWikiCollectionMapping.WithContext(s.ctx).Create(mapping)
			if err != nil {
				xlog.Log.Errorf("创建outline文件夹 并保存db 失败: name:%s  data:%v", base, mapping)
				continue
			}
			collectionId = mapping.CollectionId
			xlog.Log.Errorf("创建outline集合 成功: name:%s  collectionId:%v", base, mapping.CollectionId)
		}

		// process file dir
		go func(rootPath, collectionId string, wg *sync.WaitGroup) {
			rootWg.Add(1)
			defer wg.Done()

			if len(collectionId) == 0 || len(rootPath) == 0 {
				return
			}

			var parentId *string
			fileSystem := os.DirFS(rootPath)
			err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					xlog.Log.Error("遍历文件路径失败: %v", err)
					return nil
				}

				// 深度遍历结束回到最上层时候恢复parentId
				if filepath.Dir(path) == rootPath {
					parentId = utils.PtrString("")
				}

				if d.IsDir() {
					// 获取目录的 parent Id
					parentId = utils.PtrString(s.processDir(path, *parentId, collectionId))
				} else {
					go func() {
						s.processFile(path, *parentId, collectionId)
					}()
				}
				return nil
			})
			xlog.Log.Errorf("遍历文件夹出错：%v", err)
		}(fileRootPath, collectionId, &rootWg)

		// wait
		rootWg.Wait()
	}
}

// processDir 处理文件夹 这里对文件夹上锁，有一个处理中或者处理成功就不处理了
func (s *SyncMarkDownFile) processDir(path, parentId, collectionId string) string {
	// 1. 检查缓存
	cacheKey := cache2.XCache.GenCollectionCacheKey(path)
	fromCache := cache2.XCache.GetDataFromCache(cacheKey)
	if fromCache != nil {
		return fromCache.(string)
	}
	// 更新缓存
	defer func() {
		cache2.XCache.SetDataToCache(cacheKey, path, cache.NoExpiration)
	}()

	// lock update
	s.getMutexForPath(path).Lock()
	defer s.getMutexForPath(path).Unlock()

	// 2. 检查数据库是否创建了这个Id
	wikiCollectionMapping, err := dao.OutlineWikiCollectionMapping.WithContext(s.ctx).
		Where(dao.OutlineWikiCollectionMapping.CollectionPath.Eq(path), dao.OutlineWikiCollectionMapping.RealCollection.Is(false)).First()
	if err != nil {
		xlog.Log.Errorf("查询outline一般子文件夹配置: rawPath:%s", path)
		return ""
	}
	if wikiCollectionMapping != nil {
		return wikiCollectionMapping.CurrentId
	}

	// 3. create new
	// 获取最后一层文件夹名称 数据库存储全路径映射
	lastPathName := filepath.Base(path)
	// 正常创建空文档集合
	collectionUUID, _ := uuid.Parse(collectionId)
	request := outline.PostDocumentsCreateJSONRequestBody{
		CollectionId: collectionUUID,
		Publish:      utils.PtrBool(false),
		Text:         utils.PtrString(lastPathName),
		Title:        lastPathName,
	}
	ok, response := client.OutlineSdk.CreateDocument(s.ctx, request)
	if !ok {
		xlog.Log.Errorf("创建outline一般子文件夹失败: rawPath:%s request:%v response:%v", path, request, response)
		return ""
	}
	xlog.Log.Infof("创建outline一般子文件夹成功: rawPath:%s collectionId:%v", path, response.JSON200.Data.Id)
	// 4. save data
	mapping := &model.OutlineWikiCollectionMapping{
		CollectionId:   collectionId,
		CurrentId:      response.JSON200.Data.Id.String(),
		ParentId:       parentId,
		CollectionPath: path,
		CollectionName: lastPathName,
		RealCollection: false,
		Sync:           true,
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
	}
	err = dao.OutlineWikiCollectionMapping.WithContext(s.ctx).Create(mapping)
	if err != nil {
		xlog.Log.Errorf("创建outline一般子文件夹 并保存db 失败: name:%s  data:%v", lastPathName, mapping)
		return ""
	}
	return mapping.CurrentId
}

// processFile 处理文件
func (s *SyncMarkDownFile) processFile(path, parentId, collectionId string) {

	// 文件详情
	fileInfo, err := os.Open(path)
	if err != nil {
		xlog.Log.Errorf("读取文件信息失败: filePath: %s error: %v", path, err)
		return
	}
	stat, err := fileInfo.Stat()
	if err != nil {
		xlog.Log.Errorf("读取文件状态失败: filePath: %s error: %v", path, err)
		return
	}
	fileSize := float64(stat.Size()) / 1024

	if filepath.Ext(stat.Name()) != ".md" {
		xlog.Log.Infof("当前文件不是md文件,跳过处理: %s", stat.Name())
		return
	}

	// 打开文件
	fileContent, err := os.ReadFile(path)
	if err != nil {
		xlog.Log.Errorf("读取文件失败: filePath: %s error: %v", path, err)
		return
	}

	// 3. save data
	mapping := &model.FileSyncRecord{
		OutlineWikiId: "",
		CollectionId:  collectionId,
		FileName:      stat.Name(),
		FilePath:      path,
		FileContent:   string(fileContent),
		FileSize:      fileSize,
	}
	err = dao.FileSyncRecord.WithContext(s.ctx).Create(mapping)
	if err != nil {
		xlog.Log.Errorf("一般子文件夹 初始化保存db 失败: name:%s  data:%v", stat.Name(), mapping)
		return
	}

	// 请求接口创建文档
	uuidParentDocId, _ := uuid.Parse(parentId)
	collectionUUID, _ := uuid.Parse(collectionId)
	request := outline.PostDocumentsCreateJSONRequestBody{
		CollectionId:     collectionUUID,
		ParentDocumentId: &uuidParentDocId,
		Publish:          utils.PtrBool(false),
		Text:             utils.PtrString(string(fileContent)),
		Title:            stat.Name(),
	}
	ok, response := client.OutlineSdk.CreateDocument(s.ctx, request)
	if !ok {
		xlog.Log.Errorf("创建outline文档失败: rawPath:%s request:%v response:%v", path, request, response)
		return
	}
	xlog.Log.Infof("创建outline文档成功: rawPath:%s wikiId:%v", path, response.JSON200.Data.Id)
	updateResult, err := dao.FileSyncRecord.WithContext(s.ctx).
		Where(dao.FileSyncRecord.FilePath.Eq(path)).
		Update(dao.FileSyncRecord.OutlineWikiId, response.JSON200.Data.Id.String())
	if err != nil {
		xlog.Log.Errorf("创建outline一般子文件夹 更新db 失败: name:%s  data:%v", stat.Name(), mapping)
		return
	}
	xlog.Log.Infof("更新outline文档数据库成功: rawPath:%s updateResult:%v", path, updateResult.RowsAffected)

}
