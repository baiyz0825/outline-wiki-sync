// outline-wiki-sync
//
// @(#)fileeventhandler.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package service

import (
	"os"

	"github.com/baiyz0825/outline-wiki-sync/utils"
	"github.com/fsnotify/fsnotify"
)

type FileChangeEventHandler interface {
	// Handle 处理文件夹或者文件变动
	Handle(isDir bool, info os.FileInfo)
}

// BaseFileChangeEventHandler 文件夹或者文件变动事件处理器基类
type BaseFileChangeEventHandler struct {
	FileChangeEventHandler
}

var eventHandler map[fsnotify.Op]FileChangeEventHandler

// RunFileEventHandler 获取对应事件的处理器
func RunFileEventHandler(event fsnotify.Event) {
	handler, ok := eventHandler[event.Op]
	if !ok {
		return
	}
	fi, err := os.Stat(event.Name)
	if err != nil {
		utils.Log.Error("获取对应文件变动事件处理器失败: %v", err)
		handler.Handle(fi.IsDir(), fi)
	}
}

func init() {
	baseFileChangeEventHandler := &BaseFileChangeEventHandler{}
	eventHandler[fsnotify.Create] = &fileCreateEventHandler{
		BaseFileChangeEventHandler: baseFileChangeEventHandler,
	}
	eventHandler[fsnotify.Remove] = &fileRemoveEventHandler{
		BaseFileChangeEventHandler: baseFileChangeEventHandler,
	}
	eventHandler[fsnotify.Rename] = &fileRenameEventHandler{
		BaseFileChangeEventHandler: baseFileChangeEventHandler,
	}
	eventHandler[fsnotify.Write] = &fileWriteEventHandler{
		BaseFileChangeEventHandler: baseFileChangeEventHandler,
	}
}

// fileCreateEventHandler 文件创建事件处理器
type fileCreateEventHandler struct {
	*BaseFileChangeEventHandler
}

func (h *fileCreateEventHandler) Handle(isDir bool, info os.FileInfo) {

}

// fileRemoveEventHandler 文件删除事件处理器
type fileRemoveEventHandler struct {
	*BaseFileChangeEventHandler
}

func (h *fileRemoveEventHandler) Handle(isDir bool, info os.FileInfo) {

}

// fileRenameEventHandler 文件重命名事件处理器
type fileRenameEventHandler struct {
	*BaseFileChangeEventHandler
}

func (h *fileRenameEventHandler) Handle(isDir bool, info os.FileInfo) {

}

// fileWriteEventHandler 文件写入事件处理器
type fileWriteEventHandler struct {
	*BaseFileChangeEventHandler
}

func (h *fileWriteEventHandler) Handle(isDir bool, info os.FileInfo) {

}
