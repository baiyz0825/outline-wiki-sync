// outline-wiki-sync
//
// @(#)cmd.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package cmd

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/baiyz0825/outline-wiki-sync/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "outline",
	Short: "同步本地markdown文档到outline",
	Long:  `同步本地markdown文档到outline,并且支持实时监听文件状态进行自动更新同,输出更新信息到db`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var (
	watchFilePath string
	dbPath        string
	syncWatch     bool
	syncCorn      string
)

func init() {
	// 初始化命令行参数
	rootCmd.PersistentFlags().StringVar(&watchFilePath, "watchFilePath", "", "要监视的文件路径")
	defaultWorkDir, _ := os.Getwd()
	rootCmd.PersistentFlags().StringVar(&dbPath, "dbPath", filepath.Join(defaultWorkDir, "outline.db"),
		"要监视的文件完整路径: 默认工作目录下的 outline.db")
	rootCmd.PersistentFlags().BoolVar(&syncWatch, "syncWatch", false, "是否需要进行实时监听")
	rootCmd.PersistentFlags().StringVar(&syncCorn, "syncCorn", "",
		"实时监听 同步时间 corn 默认: 每10min一次 */10 * * * *")
	_ = rootCmd.MarkPersistentFlagRequired("watchFilePath")

	// init db

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Log.Infof("Start run ....")
		// func
		if syncWatch {
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				// run service
			}()

			// 实时监听
			wg.Wait()
		} else {
			// run service
		}
		utils.Log.Infof("执行结束, exit ... ")
		os.Exit(1)
	}
}
