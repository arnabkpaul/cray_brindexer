package fsentity

import (
	"cray.com/brindexer/utils"
	"fmt"
	"os"
	"syscall"
)

type PathEntity struct {
	finfo os.FileInfo
	pinfo *PathInfo
	pseq  int64
}

func NewPathEntity(finfo os.FileInfo, pinfo *PathInfo, seq int64) *PathEntity {
	return &PathEntity{finfo: finfo, pinfo: pinfo, pseq: seq}
}

func (ent *PathEntity) CommitSqls() string {
	if ent.finfo == nil {
		return fmt.Sprintf(deletePathSql, ent.PathMd5(), ent.PathMd5())
	}
	stat := ent.finfo.Sys().(*syscall.Stat_t)
	mtime, _, _ := utils.TimesFromStat_s(stat)
	return fmt.Sprintf(insertPathSql, ent.RPath(), ent.PathMd5(), mtime, 0, 0, ent.pseq)
}

func (ent *PathEntity) PathMd5() string {
	return ent.pinfo.pathmd5
}

func (ent *PathEntity) RPath() string {
	return ent.pinfo.rpath
}

func (ent *PathEntity) Name() string {
	return ent.finfo.Name()
}
