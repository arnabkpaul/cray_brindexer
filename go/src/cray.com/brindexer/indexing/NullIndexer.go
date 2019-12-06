package indexing

import (
	"fmt"
	"github.com/pkg/xattr"
	"os"
	"sync"
)

type NullIndexer struct {
	stat      bool
	xattr     bool
	taskQueue chan string
	waitGroup *sync.WaitGroup
}

func NewNullIndexer(threads int, stat bool, xattr bool) *NullIndexer {

	taskQueue := make(chan string, 1000)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(threads)
	indexer := &NullIndexer{stat, xattr, taskQueue, &waitGroup}
	indexer.init(threads)
	return indexer
}

func statWorker(jobs <-chan string, wg *sync.WaitGroup, ni *NullIndexer) {
	for path := range jobs {
		//Do stat ahead, so MDS can cache the stat
		if ni.stat {
			fi, err := os.Lstat(path)
			if err != nil {
				fmt.Println(err, fi, path)
			}

		}
		if ni.xattr {
			data, err := xattr.LGet(path, "lustre.lov")
			if err == nil {
				lsttrs := string(data)
				if len(lsttrs) > 0 {
					fmt.Print(lsttrs)
				}
			}
		}
	}
	wg.Done()
}

func (ni *NullIndexer) init(threads int) {

	for w := 0; w < threads; w++ {
		go statWorker(ni.taskQueue, ni.waitGroup, ni)
	}
}

func (ni *NullIndexer) queueEntry(path string, name string) error {
	if ni.stat || ni.xattr {
		ni.taskQueue <- path
	}
	return nil
}

func (ni *NullIndexer) ProcessRootDir(path string, name string) error {
	return ni.queueEntry(path, name)
}
func (ni *NullIndexer) ProcessLeafDir(path string, name string) error {
	return ni.queueEntry(path, name)
}
func (ni *NullIndexer) ProcessMidDir(path string, name string) error {
	return ni.queueEntry(path, name)
}
func (ni *NullIndexer) ProcessMidFile(path string) error {
	return ni.queueEntry(path, "")
}

func (ni *NullIndexer) Close() {
	close(ni.taskQueue)
	ni.waitGroup.Wait()
}
