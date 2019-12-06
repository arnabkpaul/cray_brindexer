package fsentity

import ()

type PathInfo struct {
	rpath   string
	pathmd5 string
}

func (pi *PathInfo) PathMd5() string {
	return pi.pathmd5
}
func (pi *PathInfo) RPath() string {
	return pi.rpath
}

type IndexNode interface {
	DBDir() string
	DBBase() *string
	BaseDir() string
	RPath(path string) string
	DBCount() int
	WriterCount() int
	DBCountPerWriter() int
	PathInfo(path string) *PathInfo
	Number2FileName(num int) string
	Number2WriterName(num int) string
	WriterKeyFromPathMd5(md5 string) int
	DBWrapperKeyFromPathMd5(md5 string) int
	DBFileFromPathMd5(md5 string) string
	IsSubIndexNode(path string) bool
	SaveLevel(level int)
	ReadLevel() int
	SaveLastScan()
	ReadLastScan() int
}
