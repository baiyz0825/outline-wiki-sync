// outline-wiki-sync
//
// @(#)sync.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package service

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/baiyz0825/outline-wiki-sdk"
	"github.com/baiyz0825/outline-wiki-sync/dao"
	"github.com/baiyz0825/outline-wiki-sync/model"
	"github.com/baiyz0825/outline-wiki-sync/utils"
	cache2 "github.com/baiyz0825/outline-wiki-sync/utils/cache"
	"github.com/baiyz0825/outline-wiki-sync/utils/client"
	"github.com/baiyz0825/outline-wiki-sync/utils/fileutils"
	"github.com/baiyz0825/outline-wiki-sync/utils/jsonutils"
	"github.com/baiyz0825/outline-wiki-sync/utils/xlog"
	"github.com/baiyz0825/outline-wiki-sync/utils/xrandonm"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

type SyncMarkDownFile struct {
	FileRootPaths []string
	filePathLocks *sync.Map
	ctx           context.Context
}

func NewSyncMarkDownFile(ctx context.Context, fileRootPaths []string) *SyncMarkDownFile {
	return &SyncMarkDownFile{
		FileRootPaths: fileRootPaths,
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
	defer rootWg.Wait()
	for _, fileRootPathItem := range s.FileRootPaths {
		rootWg.Add(1)
		go func(fileRootPath string, wg *sync.WaitGroup) {
			defer wg.Done()
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
					Color:       utils.PtrString(xrandonm.GenerateRandomColor()),
					Name:        base,
					Private:     utils.PtrBool(true),
				}
				ok, response := client.OutlineSdk.CreateCollection(s.ctx, request)
				if !ok {
					xlog.Log.Errorf("创建outline文件夹失败: rawPath:%s request:%v response:%v", fileRootPath, request, response)
					return
				}
				xlog.Log.Debugf("创建outline文件夹成功: rawPath:%s collectionId:%v", fileRootPath, jsonutils.ToJsonStr(response))
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
					return
				}
				collectionId = mapping.CollectionId
				xlog.Log.Infof("创建outline集合 成功: name:%s  collectionId:%v", base, mapping.CollectionId)
			}

			// process file dir

			if len(collectionId) == 0 || len(fileRootPath) == 0 {
				return
			}

			parentId := ""
			fileSystem := os.DirFS(fileRootPath)
			err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					xlog.Log.Error("遍历文件路径失败: %v", err)
					return nil
				}

				if path == "." {
					return nil
				}

				absPath := filepath.Join(fileRootPath, path)
				// 深度遍历结束回到最上层时候恢复parentId
				if filepath.Dir(absPath) == fileRootPath {
					parentId = ""
				}

				if d.IsDir() && !fileutils.CheckIsEmptyDir(absPath) {
					// 获取目录的 parent Id
					parentId = s.processDir(absPath, collectionId)
					return nil
				}

				// processFile
				s.processFile(absPath, parentId, collectionId)

				return nil
			})

			if err != nil {
				xlog.Log.Errorf("遍历文件夹结束但是存在错误: %v", err)
				return
			}
			xlog.Log.Infof("遍历文件夹出结束")
		}(fileRootPathItem, &rootWg)
	}
}

func (s *SyncMarkDownFile) getPathDirParentId(absPath string) string {
	// 1. 检查缓存
	cacheKey := cache2.XCache.GenCollectionCacheKey(absPath)
	fromCache := cache2.XCache.GetDataFromCache(cacheKey)
	if fromCache != nil {
		return fromCache.(string)
	}

	// 2. 检查数据库是否创建了这个Id
	wikiCollectionMapping, err := dao.OutlineWikiCollectionMapping.WithContext(s.ctx).
		Where(dao.OutlineWikiCollectionMapping.CollectionPath.Eq(absPath), dao.OutlineWikiCollectionMapping.RealCollection.Is(false)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		xlog.Log.Errorf("查询outline一般子文件夹配置: rawPath:%s", absPath)
		return ""
	}
	if wikiCollectionMapping != nil {
		return wikiCollectionMapping.CurrentId
	}
	return ""
}

// processDir 处理文件夹 这里对文件夹上锁，有一个处理中或者处理成功就不处理了
func (s *SyncMarkDownFile) processDir(absPath, collectionId string) string {

	parentId := s.getPathDirParentId(filepath.Dir(absPath))

	// 3. create new
	// 获取最后一层文件夹名称 数据库存储全路径映射
	lastPathName := filepath.Base(absPath)
	// 正常创建空文档集合
	collectionUUID, _ := uuid.Parse(collectionId)
	request := outline.PostDocumentsCreateJSONRequestBody{
		CollectionId: collectionUUID,
		Publish:      utils.PtrBool(true),
		Text:         utils.PtrString(lastPathName),
		Title:        lastPathName,
	}
	if len(parentId) != 0 {
		uuidParentDocId, _ := uuid.Parse(parentId)
		request.ParentDocumentId = &uuidParentDocId
	}
	ok, response := client.OutlineSdk.CreateDocument(s.ctx, request)
	if !ok {
		xlog.Log.Errorf("创建outline一般子文件夹失败: rawPath:%s request:%s response:%s", absPath, jsonutils.ToJsonStr(request), jsonutils.ToJsonStr(response))
		return ""
	}
	xlog.Log.Infof("创建outline一般子文件夹成功: rawPath:%s Id:%v parentId:%v", absPath, response.JSON200.Data.Id, response.JSON200.Data.ParentDocumentId)
	// 4. save data
	mapping := &model.OutlineWikiCollectionMapping{
		CollectionId:   collectionId,
		CurrentId:      response.JSON200.Data.Id.String(),
		ParentId:       parentId,
		CollectionPath: absPath,
		CollectionName: lastPathName,
		RealCollection: false,
		Sync:           true,
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
	}
	err := dao.OutlineWikiCollectionMapping.WithContext(s.ctx).Create(mapping)
	if err != nil {
		xlog.Log.Errorf("创建outline一般子文件夹 并保存db 失败: name:%s  data:%v", lastPathName, mapping)
		return ""
	}
	// update cache
	cache2.XCache.SetDataToCache(cache2.XCache.GenCollectionCacheKey(absPath), mapping.CurrentId, cache.NoExpiration)
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
	collectionUUID, _ := uuid.Parse(collectionId)
	request := outline.PostDocumentsCreateJSONRequestBody{
		CollectionId: collectionUUID,
		Publish:      utils.PtrBool(true),
		Text:         utils.PtrString(string(fileContent)),
		Title:        stat.Name(),
	}
	// 最一层文件没有父 parentDoc
	if len(parentId) != 0 {
		uuidParentDocId, _ := uuid.Parse(parentId)
		request.ParentDocumentId = &uuidParentDocId
	}
	ok, response := client.OutlineSdk.CreateDocument(s.ctx, request)
	if !ok {
		xlog.Log.Errorf("创建outline文档失败: rawPath:%s request:%s response:%s", path, jsonutils.ToJsonStr(request), jsonutils.ToJsonStr(request))
		return
	}
	xlog.Log.Infof("创建outline文档成功: rawPath:%s wikiId:%v parentId: %s", path, response.JSON200.Data.Id, response.JSON200.Data.ParentDocumentId)
	updateResult, err := dao.FileSyncRecord.WithContext(s.ctx).
		Where(dao.FileSyncRecord.FilePath.Eq(path)).
		UpdateSimple(dao.FileSyncRecord.OutlineWikiId.Value(response.JSON200.Data.Id.String()), dao.FileSyncRecord.Sync.Value(true))
	if err != nil {
		xlog.Log.Errorf("创建outline一般子文件夹 更新db 失败: name:%s  data:%v", stat.Name(), mapping)
		return
	}
	xlog.Log.Infof("更新outline文档数据库成功: rawPath:%s updateResult:%v", path, updateResult.RowsAffected)

}
