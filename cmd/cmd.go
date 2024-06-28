// outline-wiki-sync
//
// @(#)cmd.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package cmd

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/baiyz0825/outline-wiki-sync/dao"
	"github.com/baiyz0825/outline-wiki-sync/service"
	"github.com/baiyz0825/outline-wiki-sync/utils/cache"
	"github.com/baiyz0825/outline-wiki-sync/utils/client"
	"github.com/baiyz0825/outline-wiki-sync/utils/ratelimit"
	"github.com/baiyz0825/outline-wiki-sync/utils/xlog"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "outline",
	Short: "同步本地markdown文档到outline",
	Long:  `同步本地markdown文档到outline,并且支持实时监听文件状态进行自动更新同,输出更新信息到db`,
	Run: func(cmd *cobra.Command, args []string) {
		Execute(args)
	},
}

var (
	watchFilePath string
	sdkAuth       string
	outlineHost   string
	dbPath        string
	syncWatch     bool
	runOnceSync   bool
	syncCorn      string
	deleteDb      bool
)

func init() {
	// init cmd params
	rootCmd.PersistentFlags().StringVar(&watchFilePath, "watchFilePath", "", "outline服务host")
	rootCmd.PersistentFlags().StringVar(&sdkAuth, "sdkAuth", "", "outline服务 api key")
	rootCmd.PersistentFlags().StringVar(&outlineHost, "outlineHost", "", "要监视的文件路径")
	defaultWorkDir, _ := os.Getwd()
	rootCmd.PersistentFlags().StringVar(&dbPath, "dbPath", filepath.Join(defaultWorkDir, "outline.db"),
		"要监视的文件完整路径: 默认工作目录下的 outline.db")
	rootCmd.PersistentFlags().BoolVar(&syncWatch, "syncWatch", false, "是否需要进行实时监听")
	rootCmd.PersistentFlags().BoolVar(&runOnceSync, "runOnceSync", true, "只同步一次")
	rootCmd.PersistentFlags().BoolVar(&deleteDb, "deleteDb", false, "是否删除db")
	rootCmd.PersistentFlags().StringVar(&syncCorn, "syncCorn", "",
		"实时监听 同步时间 corn 默认: 每10min一次 */10 * * * *")
	_ = rootCmd.MarkPersistentFlagRequired("watchFilePath")
	_ = rootCmd.MarkPersistentFlagRequired("outlineHost")
	_ = rootCmd.MarkPersistentFlagRequired("sdkAuth")
}

func check() {
	if runOnceSync && syncWatch {
		xlog.Log.Errorf("runOnceSync 和 syncWatch 不能同时设置")
		os.Exit(1)
	}
	// file path check
	if len(watchFilePath) == 0 {
		xlog.Log.Errorf("watchFilePath 不能为空")
		os.Exit(1)
	}
	if info, err := os.Stat(watchFilePath); err != nil {
		if os.IsNotExist(err) {
			xlog.Log.Errorf("文件路径不存在: %s", watchFilePath)
		} else {
			xlog.Log.Errorf("文件路径无效: %s", watchFilePath)
		}
		os.Exit(1)
	} else {
		if !info.IsDir() {
			xlog.Log.Errorf("文件路径不是目录: %s", watchFilePath)
			os.Exit(1)
		}
	}
}

func RunRootCmd() {
	_ = rootCmd.Execute()
}

func Execute(args []string) {

	// check
	check()
	// rateLimit
	ratelimit.Init()
	// init db
	dao.Init(dbPath, deleteDb)
	// init outline client
	client.Init(outlineHost, sdkAuth)
	// init cache
	cache.Init()

	// 创建一个上下文对象
	ctx, cancel := context.WithCancel(context.Background())

	// 捕获中断信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// 启动goroutine监听信号
	go func() {
		// 等待中断信号
		<-signalChan

		// 收到中断信号后，调用cancel函数通知goroutine停止运行
		cancel()
	}()

	cmdMainFunc(ctx)

	// 等待任务完成或收到中断信号
	<-ctx.Done()
}

func cmdMainFunc(ctx context.Context) {
	xlog.Log.Infof("Start run ....")
	fileRootPath := make([]string, 0)
	fileRootPath = append(fileRootPath, watchFilePath)
	// func
	var mainWg sync.WaitGroup

	// run sync markDown
	go func() {
		mainWg.Add(1)
		defer mainWg.Done()
		if runOnceSync {
			service.NewSyncMarkDownFile(ctx, fileRootPath).SyncMarkdownFile()
		}
	}()

	// run sync watchDir
	go func() {
		mainWg.Add(1)
		defer mainWg.Done()
		if syncWatch {
			var wg sync.WaitGroup
			for _, pathItem := range fileRootPath {
				wg.Add(1)
				go func(path string) {
					defer wg.Done()
					// run service
					service.NewFileWatch(ctx, path).WatchDir()
				}(pathItem)
			}
			// 实时监听
			wg.Wait()
		}
	}()
	mainWg.Wait()
	xlog.Log.Infof("执行结束, exit ... ")
}
