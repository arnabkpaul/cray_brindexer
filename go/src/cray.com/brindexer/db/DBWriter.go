package db

import (
	"cray.com/brindexer/fsentity"
	_ "fmt"
	"path/filepath"
	"sync"
)

const queueSize = 1000

type DBWriter struct {
	indexNode    fsentity.IndexNode
	dbWrapperMap map[int]*SqliteWrapper
	taskQueue    chan *CommitTask
	waitGroup    *sync.WaitGroup
}

func NewDBWriter(indexNode fsentity.IndexNode, waitGroup *sync.WaitGroup, dirIndex int) *DBWriter {
	writerName := indexNode.Number2WriterName(dirIndex)

	dbWrapperMap := make(map[int]*SqliteWrapper)
	for i := 0; i < indexNode.DBCountPerWriter(); i++ {
		dbFileName := indexNode.Number2WriterName(i)
		dbFile := filepath.Join(indexNode.DBDir(), writerName+dbFileName+".db")
		dbWrapperMap[i] = NewSqliteWrapper(dbFile)

	}
	taskQueue := make(chan *CommitTask, queueSize)

	dbw := &DBWriter{indexNode, dbWrapperMap, taskQueue, waitGroup}
	dbw.init(dirIndex)
	return dbw
}

func (w *DBWriter) init(index int) {
	go commitWorker(index, w.taskQueue, w)
}

func (w *DBWriter) flushOne(dbw *SqliteWrapper) {
	recs := dbw.records()
	if len(recs) == 0 {
		return
	}
	w.taskQueue <- &CommitTask{recs, dbw}
	//CommitTask{recs, dbw}.commit()
}

func (w *DBWriter) Add(ent fsentity.FSEntity) {
	key := w.indexNode.DBWrapperKeyFromPathMd5(ent.PathMd5())
	dbw := w.dbWrapperMap[key]
	cached := dbw.Add(ent)
	if cached != nil {
		w.taskQueue <- &CommitTask{cached, dbw}
		//CommitTask{cached, dbw}.commit()
	}

}

func (w *DBWriter) flush() {
	for _, dbw := range w.dbWrapperMap {
		w.flushOne(dbw)
	}
	close(w.taskQueue)
}

func (w *DBWriter) close() {
	for _, dbw := range w.dbWrapperMap {
		dbw.close()
	}
}

func (w *DBWriter) totalCommits() int64 {
	total := int64(0)
	for _, dbw := range w.dbWrapperMap {
		total += dbw.RecordCnt()
	}
	return total
}
