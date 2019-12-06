package query

import (
	"cray.com/brindexer/fsentity"
	_ "fmt"
	"path/filepath"
	"sync"
)

const summaryWorkerCnt = 16

type LeveledSummaryExecutor struct {
	indexNode fsentity.IndexNode
	grandSum  *Summary
	level     int
	lock      *sync.Mutex
	taskQueue chan string
	wg        *sync.WaitGroup
}

func summaryWorker(jobs <-chan string, le *LeveledSummaryExecutor) {
	total := NewSummary()
	for path := range jobs {
		//sum := exe.Execute()
		indexNode := fsentity.NewDefaultIndexNode(path, le.indexNode.DBBase())
		sumx := NewSummaryExecutor(&indexNode)
		sum := sumx.Execute()
		total.Add(sum)
	}
	le.addSummary(total)
	le.wg.Done()
}

func NewLeveledSummaryExecutor(indexNode fsentity.IndexNode, level int) *LeveledSummaryExecutor {

	taskQueue := make(chan string, 1000)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(summaryWorkerCnt)

	le := &LeveledSummaryExecutor{indexNode, &Summary{0, 0, 0, 0, 0, &sync.Mutex{}},
		level, &sync.Mutex{}, taskQueue, waitGroup}
	le.init()
	return le

}
func (le *LeveledSummaryExecutor) init() {
	for w := 0; w < summaryWorkerCnt; w++ {
		go summaryWorker(le.taskQueue, le)
	}
}

func (le *LeveledSummaryExecutor) processIndexNode(path string) error {
	//indexNode := fsentity.NewDefaultIndexNode(path, le.indexNode.DBBase())
	//sumx := NewSummaryExecutor(&indexNode)
	le.taskQueue <- path

	return nil
}

func (le *LeveledSummaryExecutor) addSummary(sum *Summary) {
	le.lock.Lock()
	defer le.lock.Unlock()
	le.grandSum.dirCnt += sum.dirCnt
	le.grandSum.fileCnt += sum.fileCnt
	le.grandSum.fileSize += sum.fileSize
	le.grandSum.lnkCnt += sum.lnkCnt
}
func (le *LeveledSummaryExecutor) ProcessRootDir(path string, name string) error {
	le.processIndexNode(path)
	if le.level <= 0 {
		// meaning it is not leveled
		return filepath.SkipDir
	}
	return nil
}
func (le *LeveledSummaryExecutor) ProcessLeafDir(path string, name string) error {
	le.processIndexNode(path)
	return filepath.SkipDir
}
func (le *LeveledSummaryExecutor) ProcessMidDir(path string, name string) error {
	return nil
}
func (le *LeveledSummaryExecutor) ProcessMidFile(path string) error {
	return nil
}

func (le *LeveledSummaryExecutor) GrandSummary() *Summary {
	return le.grandSum
}

func (le *LeveledSummaryExecutor) Close() {
	close(le.taskQueue)
	le.wg.Wait()
}
