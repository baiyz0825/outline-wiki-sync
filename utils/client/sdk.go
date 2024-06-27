// outline-wiki-sync
//
// @(#)outline.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package client

import (
	"context"
	"net/http"
	"os"

	"github.com/baiyz0825/outline-wiki-sdk"
	"github.com/baiyz0825/outline-wiki-sync/utils"
	"github.com/baiyz0825/outline-wiki-sync/utils/ratelimit"
)

type OutLineSdk struct {
	OutlineClientWithResponses *outline.ClientWithResponses
	OutlineClient              *outline.Client
}

var OutlineSdk *OutLineSdk

func Init(host string) {
	Client, err := outline.NewClientWithResponses(host)
	if err != nil {
		utils.Log.Errorf("初始化outLine客户端失败: %s", err)
		os.Exit(1)
	}
	OutlineSdk = &OutLineSdk{
		OutlineClientWithResponses: Client,
	}
}

// CreateCollection 创建集合
func (s *OutLineSdk) CreateCollection(ctx context.Context, request outline.PostCollectionsCreateJSONRequestBody) (bool, outline.PostCollectionsCreateResponse) {
	response := outline.PostCollectionsCreateResponse{}

	f := func(responseP *outline.PostCollectionsCreateResponse) bool {
		response, err := s.OutlineClientWithResponses.PostCollectionsCreateWithResponse(context.Background(), request)
		if err != nil || response.StatusCode() != http.StatusOK || response.JSON200 == nil {
			utils.Log.Errorf("创建outline文件夹失: %v", response)
			return false
		}
		return true
	}

	if !ratelimit.LimitRun[outline.PostCollectionsCreateResponse](ctx, &response, f) {
		return false, response
	}
	return true, response
}
