package indexing

import (
	"cray.com/brindexer/db"
	"cray.com/brindexer/fsentity"
	_ "fmt"
	"path/filepath"
	"sync"
	"time"
)

const queueSize = 1000

//Implement ScanCallBack
//This is the simple indexer that just add everything to commiter
type BaseIndexer struct {
	indexNode fsentity.IndexNode
	committer *db.EntityCommitter
	taskQueue chan string
	waitGroup *sync.WaitGroup
	threads   int
	//Artificial sequence number for each run. The purpose is
	//to support  1. query for latest version, 2. delete versions older than latest.
	//3. No explicite sql deletion is requred to support file deletion.
	seq int64
}

func indexWorker(jobs <-chan string, commiter *db.EntityCommitter, wg *sync.WaitGroup, seq int64) {
	for path := range jobs {
		//fmt.Println(path)
		commiter.Add(path, seq)
	}
	wg.Done()
}

func NewBaseIndexer(indexNode fsentity.IndexNode, threads int) *BaseIndexer {
	committer := db.NewEntityCommitter(indexNode)
	taskQueue := make(chan string, queueSize)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(threads)
	seq := time.Now().Unix()
	indexer := &BaseIndexer{indexNode: indexNode, committer: committer,
		taskQueue: taskQueue, waitGroup: &waitGroup, threads: threads, seq: seq}
	indexer.init(threads)
	return indexer
}
func (di *BaseIndexer) init(threads int) {

	for w := 0; w < threads; w++ {
		go indexWorker(di.taskQueue, di.committer, di.waitGroup, di.seq)
	}
}

func (di *BaseIndexer) processDir(path string, name string) error {
	if name == fsentity.IndexDirName() ||
		di.indexNode.IsSubIndexNode(path) {
		//fmt.Printf("Skip subdir : %s\n", path)
		return filepath.SkipDir
	}
	di.taskQueue <- path
	return nil
}

func (di *BaseIndexer) processFile(path string) error {
	di.taskQueue <- path
	return nil
}

func (di *BaseIndexer) CommittedRecCnt() int64 {
	return di.committer.TotalCommits()
}
func (di *BaseIndexer) close() {
	close(di.taskQueue)
	di.waitGroup.Wait()
	di.committer.WaitForCompletion()
	//fmt.Println("Total DB commited records", di.committer.TotalCommits())
}
