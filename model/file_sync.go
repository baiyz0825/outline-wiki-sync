// outline-wiki-sync
//
// @(#)file_sync.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package model

import (
	"time"

	"gorm.io/gorm"
)

type FileSyncRecord struct {
	Id            string         `gorm:"column:id;primaryKey" comment:"记录ID"`
	OutlineWikiId string         `gorm:"column:outline_wiki_id" comment:"大纲Wiki ID"`
	CollectionId  string         `gorm:"column:collection_id" comment:"集合ID"`
	FileName      string         `gorm:"column:file_name" comment:"文件名"`
	FileSize      float64        `gorm:"column:file_size" comment:"文件大小"`
	FilePath      string         `gorm:"column:file_path" comment:"文件路径"`
	FileContent   string         `gorm:"column:file_content" comment:"文件内容"`
	Sync          bool           `gorm:"column:sync" comment:"同步标志"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" comment:"创建时间"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" comment:"更新时间"`
	Deleted       gorm.DeletedAt `gorm:"column:deleted" comment:"删除标志"`
}

// TableName 添加约束
func (FileSyncRecord) TableName() string {
	return "file_sync_record"
}

type OutlineWikiCollectionMapping struct {
	Id             uint           `gorm:"column:id;primaryKey;autoIncrement" comment:"主键ID"`
	CollectionId   string         `gorm:"column:collection_id;index" comment:"集合ID"`
	CurrentId      string         `gorm:"column:current_id;index" comment:"当前子文件夹ID"`
	ParentId       string         `gorm:"column:parent_id;index" comment:"父ID"`
	CollectionPath string         `gorm:"column:collection_path;index" comment:"集合路径"`
	CollectionName string         `gorm:"column:collection_name;index" comment:"集合名称"`
	RealCollection bool           `gorm:"column:real_collection;index" comment:"是否是走createCollection创建的还是一个子文档集合"`
	Sync           bool           `gorm:"column:sync" comment:"同步标志"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" comment:"创建时间"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" comment:"更新时间"`
	Deleted        gorm.DeletedAt `gorm:"column:deleted" comment:"删除标志"`
}

// TableName 添加约束
func (OutlineWikiCollectionMapping) TableName() string {
	return "outline_wiki_collection_mapping"
}
