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

var reqAuth outline.RequestEditorFn

func Init(host, sdkAuth string) {
	client, err := outline.NewClientWithResponses(
		host+"/api",
		outline.WithHTTPClient(&http.Client{
			Transport: &LoggingTransport{Transport: http.DefaultTransport},
		}),
	)
	if err != nil {
		xlog.Log.Errorf("初始化outLine客户端失败: %s", err)
		os.Exit(1)
	}
	reqAuth = func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sdkAuth))
		xlog.Log.Debugf("增加请求key: %+v", req.Header)
		return nil
	}
	OutlineSdk = &OutLineSdk{
		OutlineClientWithResponses: client,
	}
}

// CreateCollection 创建集合
func (s *OutLineSdk) CreateCollection(ctx context.Context, request outline.PostCollectionsCreateJSONRequestBody) (bool, outline.PostCollectionsCreateResponse) {
	f := func(ctx context.Context, request *outline.PostCollectionsCreateJSONRequestBody) *outline.PostCollectionsCreateResponse {
		respClient, err := s.OutlineClientWithResponses.PostCollectionsCreateWithResponse(context.Background(), *request, reqAuth)
		if err != nil || respClient.StatusCode() != http.StatusOK || respClient.JSON200 == nil {
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
		response, err := s.OutlineClientWithResponses.PostDocumentsCreateWithResponse(context.Background(), *request, reqAuth)
		if err != nil || response.StatusCode() != http.StatusOK || response.JSON200 == nil {
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
