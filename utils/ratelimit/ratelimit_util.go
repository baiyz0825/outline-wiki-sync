// outline-wiki-sync
//
// @(#)ratelimit_util.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package ratelimit

import (
	"context"
	"time"

	"golang.org/x/time/rate"

	"github.com/baiyz0825/outline-wiki-sync/utils"
)

var limiter *rate.Limiter

func Init() {
	limiter = rate.NewLimiter(15*rate.Every(time.Second), 1)
}

// LimitRun 限流运行
func LimitRun[T any](ctx context.Context, response *T, fuc func(*T) bool) bool {
	err := limiter.Wait(ctx)
	if err != nil {
		utils.Log.Errorf("限流控制失败: %v", err)
		return false
	}
	return fuc(response)
}
