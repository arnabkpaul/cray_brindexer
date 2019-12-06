package query

import (
	"cray.com/brindexer/fsentity"
	"sync"
)

type SummaryExecutor struct {
	indexNode fsentity.IndexNode
	summary   *Summary
}

func NewSummaryExecutor(indexNode fsentity.IndexNode) *SummaryExecutor {
	return &SummaryExecutor{indexNode, &Summary{0, 0, 0, 0, 0, &sync.Mutex{}}}
}

func (s *SummaryExecutor) Execute() *Summary {
	sum := NewSummary()
	dbDir := s.indexNode.DBDir()
	counter := RecordCounter{int64(0), &sync.Mutex{}}
	query := "select count(*) from entries_0 where type ='f'"
	QueryAll(dbDir, query, s.indexNode.DBCount(), &counter)
	sum.fileCnt = counter.RecordCount()

	counter.reset()
	query = "select sum(size) from entries_0"
	QueryAll(dbDir, query, s.indexNode.DBCount(), &counter)
	sum.fileSize = counter.RecordCount()
	//fmt.Println("Total File Size:", sum.fileSize)

	counter.reset()
	query = "select count(*) from path"
	QueryAll(dbDir, query, s.indexNode.DBCount(), &counter)
	sum.dirCnt = counter.RecordCount()

	counter.reset()
	query = "select count(*) from entries_0 where type ='l'"
	QueryAll(dbDir, query, s.indexNode.DBCount(), &counter)
	sum.lnkCnt = counter.RecordCount()

	return sum

}
