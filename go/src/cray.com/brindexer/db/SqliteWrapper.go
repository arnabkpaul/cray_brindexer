package db

import (
	"cray.com/brindexer/fsentity"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sync"
)

type SqliteWrapper struct {
	dbFile   string
	recCache []fsentity.FSEntity
	lock     *sync.Mutex
	recCnt   int64
}

func NewSqliteWrapper(dbFile string) *SqliteWrapper {
	recCache := make([]fsentity.FSEntity, 0, fsentity.MaxCacheSize)
	w := &SqliteWrapper{dbFile: dbFile, recCache: recCache, lock: &sync.Mutex{}, recCnt: 0}
	w.init()
	return w
}

func (conn *SqliteWrapper) init() {

	db := conn.Connection()
	defer db.Close()
	db.Exec(fsentity.FileSql)
	db.Exec(fsentity.PathSql)
}
func (conn *SqliteWrapper) Connection() *sql.DB {
	db, err := sql.Open("sqlite3", conn.dbFile+"?_journal_mode=wal&_timeout=10000")
	if err != nil {
		fmt.Println(err)
	}

	return db
}

func (conn *SqliteWrapper) Add(ent fsentity.FSEntity) []fsentity.FSEntity {
	conn.lock.Lock()
	defer conn.lock.Unlock()
	conn.recCache = append(conn.recCache, ent)
	if len(conn.recCache) >= fsentity.MaxCacheSize {
		//fmt.Println("Flush:", conn.recCache)
		cached := conn.recCache
		conn.recCache = make([]fsentity.FSEntity, 0, fsentity.MaxCacheSize)
		return cached
	}
	return nil
}

func (conn *SqliteWrapper) Exec(sql string) {
	db := conn.Connection()
	defer db.Close()
	db.Exec(sql)
}

func (conn *SqliteWrapper) createIndices() {
	if !conn.DbExists() {
		return
	}
	db := conn.Connection()
	defer db.Close()
	db.Exec("create index if not exists idx_f_size on entries_0 (size)")
	db.Exec("create index if not exists idx_f_mtime on entries_0 (mtime)")
	db.Exec("create index if not exists idx_f_seq on entries_0 (seq)")
	db.Exec("create index if not exists idx_f_type on entries_0 (type)")
	db.Exec("create index if not exists idx_path_mtime on path (lastmtime)")
	db.Exec("create index if not exists idx_path_seq on path (pseq)")

}

func (conn *SqliteWrapper) records() []fsentity.FSEntity {
	conn.lock.Lock()
	defer conn.lock.Unlock()
	cached := conn.recCache
	conn.recCache = make([]fsentity.FSEntity, 0, fsentity.MaxCacheSize)
	return cached
}

func (conn *SqliteWrapper) RecordCnt() int64 {
	return conn.recCnt
}

func (conn *SqliteWrapper) addRecordCnt(cnt int) {
	conn.recCnt += int64(cnt)
}

func (conn *SqliteWrapper) close() {
	conn.createIndices()
}

func (conn *SqliteWrapper) DbExists() bool {
	_, err := os.Stat(conn.dbFile)
	return !os.IsNotExist(err)

}
