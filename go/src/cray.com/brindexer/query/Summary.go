package query

import (
	"fmt"
	"sync"
)

type Summary struct {
	fileCnt  int64
	dirCnt   int64
	fileSize int64
	lnkCnt   int64
	total    int64
	lock     *sync.Mutex
}

type SummaryCallBack interface {
	ProcessSummary(sum *Summary)
}

func NewSummary() *Summary {
	return &Summary{0, 0, 0, 0, 0, &sync.Mutex{}}

}
func (s *Summary) FileCnt() int64 {
	return s.fileCnt
}

func (s *Summary) LnkCnt() int64 {
	return s.lnkCnt
}
func (s *Summary) DirCnt() int64 {
	return s.dirCnt
}

func (s *Summary) FileSize() int64 {
	return s.fileSize
}

func (s *Summary) Total() int64 {
	return s.fileCnt + s.DirCnt() + s.LnkCnt()
}

func (s *Summary) String() string {
	return fmt.Sprintf("File: %d, Dir: %d, Link: %d, File size: %d, Total count:%d",
		s.fileCnt, s.dirCnt, s.lnkCnt, s.fileSize, s.Total())
}

func (s *Summary) Add(sum *Summary) *Summary {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.dirCnt += sum.dirCnt
	s.fileCnt += sum.fileCnt
	s.fileSize += sum.fileSize
	s.lnkCnt += sum.lnkCnt
	return s
}
