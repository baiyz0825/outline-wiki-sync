// outline-wiki-sync
//
// @(#)ratelimit_util.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package ratelimit

import (
	"context"
	"time"

	"golang.org/x/time/rate"

	"github.com/baiyz0825/outline-wiki-sync/utils/xlog"
)

var limiter *rate.Limiter

func Init() {
	limiter = rate.NewLimiter(15*rate.Every(time.Minute), 1)
}

// LimitRunRequest 限流运行
func LimitRunRequest[T any, R any](ctx context.Context, request *T, run func(ctx context.Context, request *T) *R) *R {
	err := limiter.Wait(ctx)
	if err != nil {
		xlog.Log.Errorf("限流控制失败: %v", err)
		return nil
	}
	return run(ctx, request)
}
