package model

import (
	"github.com/davyxu/tabtoy/util"
)

type FileMeta struct {
	FileName string
	Mode     string
}

type Globals struct {
	Types *TypeManager  // 定义的表头结构和枚举
	Datas *TableManager // 表数据

	DataFileGetter util.FileGetter // 数据源
	ParaLoading    bool            // 并发加载
	CacheDir       string          // 缓存路径

	InputFiles []*FileMeta // 输入的文件列表

	Version           string // 工具版本号
	PackageName       string // 文件生成时的包名
	CombineStructName string // 包含最终表所有数据的根结构
}

func (self *Globals) AddFile(mode, fileName string) {
	var meta FileMeta
	meta.FileName = fileName
	meta.Mode = mode
	self.InputFiles = append(self.InputFiles, &meta)
}

func NewGlobals() *Globals {
	return &Globals{
		Types: NewTypeManager(),
		Datas: NewTableManager(),
	}
}
