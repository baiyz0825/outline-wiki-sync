// outline-wiki-sync
//
// @(#)dao.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package dao

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/baiyz0825/outline-wiki-sync/utils/xlog"
)

var initSql = `
	CREATE TABLE IF NOT EXISTS file_sync_record
	(
	    id              INTEGER PRIMARY KEY AUTOINCREMENT,
	    outline_wiki_id TEXT,
	    collection_id   TEXT,
	    file_name       TEXT,
	    file_size       TEXT,
	    file_path       TEXT,
	    file_content    TEXT,
	    sync            INTEGER,
	    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    deleted         INTEGER
	);
	
	CREATE INDEX IF NOT EXISTS idx_outline_wiki_id ON file_sync_record(outline_wiki_id);
	CREATE INDEX IF NOT EXISTS idx_collection_id ON file_sync_record(collection_id);
	CREATE INDEX IF NOT EXISTS idx_file_name ON file_sync_record(file_name);
	
	CREATE TABLE IF NOT EXISTS outline_wiki_collection_mapping
	(
	    id              INTEGER PRIMARY KEY AUTOINCREMENT,
	    collection_id   TEXT,
	    current_id      TEXT,
	    parent_id       TEXT,
	    collection_path TEXT,
	    collection_name TEXT,
	    sync            INTEGER,
	    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    deleted         INTEGER
	);
	
	CREATE INDEX IF NOT EXISTS idx_collection_id ON outline_wiki_collection_mapping(collection_id);
	CREATE INDEX IF NOT EXISTS idx_collection_path ON outline_wiki_collection_mapping(collection_path);
	CREATE INDEX IF NOT EXISTS idx_collection_name ON outline_wiki_collection_mapping(collection_name);
`

func Init(dbPath string) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		xlog.Log.Errorf("数据库初始化失败: %v", err)
	}
	// db, err := gorm.Open(mysql.New(mysql.Config{
	// 	DSN:                       "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // DSN data source name
	// 	DefaultStringSize:         256,                                                                        // string 类型字段的默认长度
	// 	DisableDatetimePrecision:  true,                                                                       // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	// 	DontSupportRenameIndex:    true,                                                                       // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	// 	DontSupportRenameColumn:   true,                                                                       // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	// 	SkipInitializeWithVersion: false,                                                                      // 根据当前 MySQL 版本自动配置
	// }), &gorm.Config{})
	// 执行初始化SQL语句
	result := db.Exec(initSql)
	if result.Error != nil {
		xlog.Log.Errorf("数据库初始化失败: %v", result.Error)
		os.Exit(1)
	}
	// 注入db到查询器中
	SetDefault(db)
	xlog.Log.Infof("数据库初始化成功: path:%s", dbPath)
}
