package query

import (
	"database/sql"
	"sync"
)

type RecordCounter struct {
	cnt  int64
	lock *sync.Mutex
}

func (r *RecordCounter) ProcessRow(rows *sql.Rows, cols []string) error {

	var cnt int64 = 0
	r.lock.Lock()
	defer r.lock.Unlock()
	rows.Scan(&cnt)
	r.cnt += cnt
	return nil
}

func (r *RecordCounter) reset() {
	r.cnt = int64(0)
}

func (r *RecordCounter) RecordCount() int64 {
	return r.cnt
}
