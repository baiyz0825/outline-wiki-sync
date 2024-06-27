// outline-wiki-sync
//
// @(#)watch.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package service

import (
	"context"
	"os"
	"path/filepath"

	"github.com/baiyz0825/outline-wiki-sync/utils/xlog"
	"github.com/fsnotify/fsnotify"
)

type FileWatch struct {
	WatchRootPath string
	watch         *fsnotify.Watcher
	ctx           context.Context
}

func NewFileWatch(ctx context.Context, rootPath string) *FileWatch {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		xlog.Log.Error("初始化文件监听失败: ", err)
	}
	w := &FileWatch{
		WatchRootPath: rootPath,
		watch:         watcher,
		ctx:           ctx,
	}
	return w
}

// WatchDir 监控目录
func (w *FileWatch) WatchDir() {
	// 通过Walk来遍历目录下的所有子目录
	_ = filepath.Walk(w.WatchRootPath, func(path string, info os.FileInfo, err error) error {
		// 判断是否为目录，监控目录,目录下文件也在监控范围内，不需要加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = w.watch.Add(path)
			if err != nil {
				return err
			}
			xlog.Log.Infof("监控 : %s", path)
		}
		return nil
	})

	// run listener
	go w.watchEvent()
}

func (w *FileWatch) watchEvent() {
	for {
		select {
		case <-w.ctx.Done():
			{
				xlog.Log.Infof("监测到退出信号,文件监听退出")
				return
			}
		case ev := <-w.watch.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					xlog.Log.Infof("创建文件 : %s", ev.Name)
					file, err := os.Stat(ev.Name)
					if err == nil && file.IsDir() {
						_ = w.watch.Add(ev.Name)
						xlog.Log.Infof("添加监控 : %s", ev.Name)
					}
					RunFileEventHandler(ev)

				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					xlog.Log.Infof("文件被写入 : %s", ev.Name)
					RunFileEventHandler(ev)
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					xlog.Log.Infof("删除文件 : %s", ev.Name)
					// 如果删除文件是目录，则移除监控
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						_ = w.watch.Remove(ev.Name)
						xlog.Log.Infof("删除监控 : %s", ev.Name)
					}

					RunFileEventHandler(ev)
				}

				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					// 如果重命名文件是目录，则移除监控 ,注意这里无法使用os.Stat来判断是否是目录了
					// 因为重命名后，go已经无法找到原文件来获取信息了,所以简单粗爆直接remove
					xlog.Log.Infof("重命名文件 : %s", ev.Name)
					_ = w.watch.Remove(ev.Name)

					RunFileEventHandler(ev)
				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					xlog.Log.Infof("修改权限 : %s", ev.Name)

					RunFileEventHandler(ev)
				}
			}
		case err := <-w.watch.Errors:
			{
				xlog.Log.Infof("文件监听出现问题: %v", err)
				return
			}
		}
	}
}
