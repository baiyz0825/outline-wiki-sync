// outline-wiki-sync
//
// @(#)dao.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package dao

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/baiyz0825/outline-wiki-sync/utils/xlog"
)

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

	// 注入db到查询器中
	SetDefault(db)
}
