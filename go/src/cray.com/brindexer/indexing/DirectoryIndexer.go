package indexing

import (
	"cray.com/brindexer/fsentity"
)

type DirectoryIndexer struct {
	indexer *BaseIndexer
}

func NewDirectoryIndexer(indexNode fsentity.IndexNode, threads int) *DirectoryIndexer {

	indexer := NewBaseIndexer(indexNode, threads)
	return &DirectoryIndexer{indexer}

}

func (di *DirectoryIndexer) ProcessRootDir(path string, name string) error {
	return di.indexer.processDir(path, name)
}
func (di *DirectoryIndexer) ProcessLeafDir(path string, name string) error {
	return di.indexer.processDir(path, name)
}
func (di *DirectoryIndexer) ProcessMidDir(path string, name string) error {
	return di.indexer.processDir(path, name)
}
func (di *DirectoryIndexer) ProcessMidFile(path string) error {
	di.indexer.processFile(path)
	return nil
}

func (di *DirectoryIndexer) Close() {
	di.indexer.close()
}

func (di *DirectoryIndexer) CommittedRecCnt() int64 {
	return di.indexer.CommittedRecCnt()
}
