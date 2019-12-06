package query

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

type DefaultRecord struct {
	osPath  string
	name    string
	size    int64
	mtime   int64
	pathmd5 string
}

type PrintQueryCallback struct {
	lock     *sync.Mutex
	rootPath string
	cnt      int64
	asis     bool
}

func NewPrintQueryCallback(rootPath string, asis bool) *PrintQueryCallback {
	return &PrintQueryCallback{&sync.Mutex{}, rootPath, 0, asis}
}
func (pc *PrintQueryCallback) exatractRow(rows *sql.Rows, cols []string) []string {
	colValues := make([]string, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i, _ := range colValues {
		columnPointers[i] = &colValues[i]
	}

	rows.Scan(columnPointers...)
	return colValues
}
func (pc *PrintQueryCallback) ProcessRow(rows *sql.Rows, cols []string) error {
	pc.lock.Lock()
	defer pc.lock.Unlock()
	if pc.asis {

		//fmt.Println(cols)
		colValues := pc.exatractRow(rows, cols)

		fmt.Println(strings.Join(colValues, ", "))
		return nil
	}
	var rec DefaultRecord

	rows.Scan(&rec.osPath, &rec.name, &rec.size, &rec.mtime, &rec.pathmd5)
	rec.osPath = filepath.Join(pc.rootPath, rec.osPath, rec.name)
	fmt.Println(rec.osPath, rec.size, rec.mtime)
	pc.cnt++
	return nil
}
func (pc *PrintQueryCallback) RecordCount() int64 {
	return pc.cnt
}
