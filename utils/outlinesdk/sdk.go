// outline-wiki-sync
//
// @(#)outline.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package outlinesdk

import (
	"os"

	"github.com/baiyz0825/outline-wiki-sdk"
	"github.com/baiyz0825/outline-wiki-sync/utils"
)

var OutlineClient *outline.Client

func Init(host string) {
	Client, err := outline.NewClient(host)
	if err != nil {
		utils.Log.Errorf("初始化outLine客户端失败: %s", err)
		os.Exit(1)
	}
	OutlineClient = Client
}
