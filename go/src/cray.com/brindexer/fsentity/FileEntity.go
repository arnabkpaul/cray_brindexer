package fsentity

import (
	"cray.com/brindexer/lustre"
	"cray.com/brindexer/utils"
	"fmt"
	"os"
	"strconv"
	_ "strings"
	"syscall"
)

type FileEntity struct {
	finfo    os.FileInfo
	pinfo    *PathInfo
	seq      int64
	layout   *lustre.OSTLayout
	hamAttrs *lustre.HSMAttrs
}

func NewFileEntity(finfo os.FileInfo, pinfo *PathInfo, seq int64, absPath string) *FileEntity {

	ha := lustre.GetHsmAttrs(absPath)
	layout := lustre.GetLayout(absPath)
	return &FileEntity{finfo: finfo, pinfo: pinfo, seq: seq, layout: layout, hamAttrs: ha}

}

func (ent *FileEntity) PathMd5() string {
	return ent.pinfo.pathmd5
}

func (ent *FileEntity) RPath() string {
	return ent.pinfo.rpath
}

func (ent *FileEntity) Name() string {
	return ent.finfo.Name()
}

func (ent *FileEntity) ftype() string {
	var ftype string = "l"
	mode := ent.finfo.Mode()
	if mode.IsRegular() {
		ftype = "f"
	} else if mode.IsDir() {
		ftype = "d"
	}
	return ftype
}

func (ent *FileEntity) CommitSqls() string {
	ftype := ent.ftype()
	stat := ent.finfo.Sys().(*syscall.Stat_t)
	mtime, atime, ctime := utils.TimesFromStat_s(stat)

	osts := ent.layout.OstIndice()
	ostStrs := ":"
	for i := range osts {
		ostStrs += strconv.Itoa(int(osts[i]))
	}
	ostStrs += ":"

	return fmt.Sprintf(InsertFileSql, ent.finfo.Name(), ftype, ent.pinfo.pathmd5,
		stat.Ino, stat.Mode, stat.Nlink, stat.Uid, stat.Gid, stat.Size, stat.Blksize, stat.Blocks, atime, mtime, ctime,
		ent.layout.PoolName(), ostStrs,
		ent.layout.MirrorState(), ent.hamAttrs.HsmArchId, ent.hamAttrs.HsmCompat, ent.hamAttrs.HsmFlags, ent.hamAttrs.HsmArchVer, ent.seq)
}
