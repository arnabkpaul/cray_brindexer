package query

import (
	"cray.com/brindexer/fsentity"
	"path/filepath"
	"sync"
)

const searchWorkerCnt = 16

type LeveledSearchExecutor struct {
	indexNode fsentity.IndexNode
	querySql  string
	level     int
	asis      bool
	lock      *sync.Mutex
	taskQueue chan *SearchExecutor
	wg        *sync.WaitGroup
	callback  QueryCallBack
}

func searchWorker(jobs <-chan *SearchExecutor, ls *LeveledSearchExecutor) {
	for exe := range jobs {
		exe.Execute(ls.querySql, ls.asis, ls.callback)
	}
	ls.wg.Done()
}

func NewLeveledSearchExecutor(indexNode fsentity.IndexNode, querySql string,
	level int, asis bool, callback QueryCallBack) *LeveledSearchExecutor {

	taskQueue := make(chan *SearchExecutor, 1000)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(searchWorkerCnt)
	ls := &LeveledSearchExecutor{indexNode, querySql,
		level, asis, &sync.Mutex{}, taskQueue, waitGroup, callback}
	ls.init()
	return ls

}
func (ls *LeveledSearchExecutor) init() {
	for w := 0; w < searchWorkerCnt; w++ {
		go searchWorker(ls.taskQueue, ls)
	}
}

func (ls *LeveledSearchExecutor) processIndexNode(path string) error {
	indexNode := fsentity.NewDefaultIndexNode(path, ls.indexNode.DBBase())
	search := NewSearchExecutor(&indexNode)
	ls.taskQueue <- search
	return nil
}
func (ls *LeveledSearchExecutor) ProcessRootDir(path string, name string) error {
	ls.processIndexNode(path)
	if ls.level <= 0 {
		//meaning it is not leveled
		return filepath.SkipDir
	}
	return nil
}
func (ls *LeveledSearchExecutor) ProcessLeafDir(path string, name string) error {
	ls.processIndexNode(path)
	return filepath.SkipDir
}
func (ls *LeveledSearchExecutor) ProcessMidDir(path string, name string) error {
	return nil
}
func (ls *LeveledSearchExecutor) ProcessMidFile(path string) error {
	return nil
}

func (ls *LeveledSearchExecutor) TotalRecordCnt() int64 {
	return ls.callback.RecordCount()
}

func (ls *LeveledSearchExecutor) Close() {
	close(ls.taskQueue)
	ls.wg.Wait()
}
