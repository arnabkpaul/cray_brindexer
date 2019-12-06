package db

import (
	"cray.com/brindexer/fsentity"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type EntityCommitter struct {
	indexNode fsentity.IndexNode
	writerMap map[int]*DBWriter
	waitGroup *sync.WaitGroup
}

func NewEntityCommitter(indexNode fsentity.IndexNode) *EntityCommitter {
	wmap := make(map[int]*DBWriter)
	wcnt := indexNode.WriterCount()
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(wcnt)
	for i := 0; i < wcnt; i++ {
		wmap[i] = NewDBWriter(indexNode, waitGroup, i)
	}

	return &EntityCommitter{indexNode, wmap, waitGroup}

}

func (c *EntityCommitter) AddEntity(ent fsentity.FSEntity) {

	//fmt.Println(ent.PathMd5())
	key := c.indexNode.WriterKeyFromPathMd5(ent.PathMd5())
	//fmt.Println(key)
	c.writerMap[key].Add(ent)
}

func (c *EntityCommitter) Add(filePath string, seq int64) {
	//1. Stat it, create a record, then cache it for commit
	fi, err := os.Lstat(filePath)
	if err != nil {
		return
	}

	if fi.IsDir() {
		pinfo := c.indexNode.PathInfo(filePath)
		c.AddEntity(fsentity.NewPathEntity(fi, pinfo, seq))
	} else {
		pdir := filepath.Dir(filePath)
		pinfo := c.indexNode.PathInfo(pdir)
		c.AddEntity(fsentity.NewFileEntity(fi, pinfo, seq, filePath))
	}

}

func (c *EntityCommitter) flush() {

	for _, w := range c.writerMap {
		w.flush()
	}
}

//Can only be called after WaitForCompletion
func (c *EntityCommitter) TotalCommits() int64 {
	total := int64(0)
	for _, w := range c.writerMap {
		total += w.totalCommits()
	}
	return total
}
func (c *EntityCommitter) WaitForCompletion() {
	fmt.Println("****Flush and wait****")
	c.flush()
	c.waitGroup.Wait()
}
