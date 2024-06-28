// outline-wiki-sync
//
// @(#)outline.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package client

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/baiyz0825/outline-wiki-sdk"
	"github.com/baiyz0825/outline-wiki-sync/utils/ratelimit"
	"github.com/baiyz0825/outline-wiki-sync/utils/xlog"
)

type OutLineSdk struct {
	OutlineClientWithResponses *outline.ClientWithResponses
	OutlineClient              *outline.Client
}

var OutlineSdk *OutLineSdk

func Init(host, sdkAuth string) {
	Client, err := outline.NewClientWithResponses(
		host+"/api",
		outline.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sdkAuth))
			return nil
		}),
	)
	if err != nil {
		xlog.Log.Errorf("初始化outLine客户端失败: %s", err)
		os.Exit(1)
	}
	OutlineSdk = &OutLineSdk{
		OutlineClientWithResponses: Client,
	}
}

// CreateCollection 创建集合
func (s *OutLineSdk) CreateCollection(ctx context.Context, request outline.PostCollectionsCreateJSONRequestBody) (bool, outline.PostCollectionsCreateResponse) {
	f := func(ctx context.Context, request *outline.PostCollectionsCreateJSONRequestBody) *outline.PostCollectionsCreateResponse {
		respClient, err := s.OutlineClientWithResponses.PostCollectionsCreateWithResponse(context.Background(), *request)
		if err != nil || respClient.StatusCode() != http.StatusOK || respClient.JSON200 == nil {
			xlog.Log.Errorf("创建outline文件夹失: %v", respClient)
			return nil
		}
		return respClient
	}
	response := ratelimit.LimitRunRequest(ctx, &request, f)
	if response == nil {
		return false, outline.PostCollectionsCreateResponse{}
	}
	return true, *response
}

// CreateDocument 创建文档
func (s *OutLineSdk) CreateDocument(ctx context.Context, request outline.PostDocumentsCreateJSONRequestBody) (bool,
	outline.PostDocumentsCreateResponse) {
	f := func(ctx context.Context, request *outline.PostDocumentsCreateJSONRequestBody) *outline.PostDocumentsCreateResponse {
		response, err := s.OutlineClientWithResponses.PostDocumentsCreateWithResponse(context.Background(), *request)
		if err != nil || response.StatusCode() != http.StatusOK || response.JSON200 == nil {
			xlog.Log.Errorf("创建outline文件失败: %v", response)
			return nil
		}
		return response
	}

	response := ratelimit.LimitRunRequest(ctx, &request, f)
	if response == nil {
		return false, outline.PostDocumentsCreateResponse{}
	}
	return true, *response
}
