package indexing

import (
	"cray.com/brindexer/fsentity"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
)

const queueLen = 1000

type ParallelIndexer struct {
	indexer   *BaseIndexer
	taskQueue chan string
	waitGroup *sync.WaitGroup
	full      bool
}

func parallelIndexWorker(jobs <-chan string, indexer *ParallelIndexer) {
	for path := range jobs {
		//Launch in separate process or thread and wait. Assume the index binay is in same dir
		cmd := filepath.Join(filepath.Dir(os.Args[0]), "index")
		var indexDir string = ""
		pindexDir := indexer.indexer.indexNode.DBBase()
		if pindexDir != nil {
			indexDir = *pindexDir
		}
		args := []string{"-t", strconv.Itoa(indexer.indexer.threads), "-index", indexDir, path}
		if indexer.full {
			args = []string{"-f", "-t", strconv.Itoa(indexer.indexer.threads), "-index", indexDir, path}
		}
		//fmt.Println("ParallelIndexer, Process Leaf:", cmd, strings.Join(args, " "))
		indexCmd := exec.Command(cmd, args...)
		indexCmd.Stdout = os.Stdout
		indexCmd.Stderr = os.Stderr
		indexCmd.Start()
		indexCmd.Wait()
	}
	//fmt.Println("subIndexWorker exits")
	indexer.waitGroup.Done()
}
func NewParallelIndexer(indexNode fsentity.IndexNode, threads int, procs int, full bool) *ParallelIndexer {
	indexer := NewBaseIndexer(indexNode, threads)
	taskQueue := make(chan string, queueLen)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(procs)
	pindexer := &ParallelIndexer{indexer, taskQueue, &waitGroup, full}
	pindexer.init(procs)
	return pindexer
}
func (pi *ParallelIndexer) init(procs int) {

	for w := 0; w < procs; w++ {
		go parallelIndexWorker(pi.taskQueue, pi)
	}
}

func (pi *ParallelIndexer) ProcessRootDir(path string, name string) error {
	return pi.indexer.processDir(path, name)
}
func (pi *ParallelIndexer) ProcessLeafDir(path string, name string) error {

	if name == fsentity.IndexDirName() {
		//fmt.Printf("Skip subdir : %s\n", path)
		return filepath.SkipDir
	}
	//fmt.Printf("ParallelIndexer, Process Leaf:", path)
	pi.taskQueue <- path
	return filepath.SkipDir
}
func (pi *ParallelIndexer) ProcessMidDir(path string, name string) error {
	return pi.indexer.processDir(path, name)
}
func (pi *ParallelIndexer) ProcessMidFile(path string) error {
	pi.indexer.processFile(path)
	return nil
}

func (pi *ParallelIndexer) Close() {
	pi.indexer.close()
	close(pi.taskQueue)
	pi.waitGroup.Wait()
}
