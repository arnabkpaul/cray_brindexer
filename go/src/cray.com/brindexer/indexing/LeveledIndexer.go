package indexing

import (
	"cray.com/brindexer/fsentity"
	"cray.com/brindexer/scan"
	_ "fmt"
	"path/filepath"
	"sync"
)

const subThreads int = 8
const subQueueLen int = 1000

type LeveledIndexer struct {
	indexer   *BaseIndexer
	taskQueue chan string
	waitGroup *sync.WaitGroup
}

func subIndexWorker(jobs <-chan string, li *LeveledIndexer) {
	for path := range jobs {
		scanner := scan.NewDirectoryScanner(0)
		//fmt.Println("Start to index leaf:", path)
		scanner.Scan(path, li)
	}
	//fmt.Println("subIndexWorker exits")
	li.waitGroup.Done()
}
func NewLeveledIndexer(indexNode fsentity.IndexNode, threads int) *LeveledIndexer {
	indexer := NewBaseIndexer(indexNode, threads)
	taskQueue := make(chan string, subQueueLen)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(subThreads)
	lindexer := &LeveledIndexer{indexer, taskQueue, &waitGroup}
	lindexer.init(subThreads)
	return lindexer
}
func (li *LeveledIndexer) init(threads int) {

	for w := 0; w < threads; w++ {
		go subIndexWorker(li.taskQueue, li)
	}
}

func (li *LeveledIndexer) ProcessRootDir(path string, name string) error {
	return li.indexer.processDir(path, name)
}
func (li *LeveledIndexer) ProcessLeafDir(path string, name string) error {

	li.taskQueue <- path
	return filepath.SkipDir
}
func (li *LeveledIndexer) ProcessMidDir(path string, name string) error {
	return li.indexer.processDir(path, name)
}
func (li *LeveledIndexer) ProcessMidFile(path string) error {
	li.indexer.processFile(path)
	return nil
}
func (di *LeveledIndexer) CommitedRecCnt() int64 {
	return di.indexer.CommittedRecCnt()
}
func (li *LeveledIndexer) Close() {
	close(li.taskQueue)
	li.waitGroup.Wait()
	li.indexer.close()
}
