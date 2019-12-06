package indexing

import (
	"cray.com/brindexer/fsentity"
	"cray.com/brindexer/query"
	"cray.com/brindexer/scan"
	"cray.com/brindexer/utils"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

const pathQuery = "select path, pmd5, lastmtime from path"
const cacheSize = 10000
const maxCache = 100000

type PathRecord struct {
	path    string
	pathmd5 string
	fi      *os.FileInfo
	mtime   int64
}

//Need this to query for single path
func (r *PathRecord) ProcessRow(rows *sql.Rows, cols []string) error {
	rows.Scan(&r.path, &r.pathmd5, &r.mtime)
	return nil
}

func (r *PathRecord) RecordCount() int64 {
	//won't matter, just satisfy interface
	return 1
}

type IncrementalIndexer struct {
	indexNode fsentity.IndexNode
	indexer   *DirectoryIndexer
	pathCache []*PathRecord
	seq       int64
	lock      *sync.Mutex
}

func NewIncrementalIndexer(indexNode fsentity.IndexNode, threads int) *IncrementalIndexer {
	seq := time.Now().Unix()
	cache := make([]*PathRecord, 0, cacheSize)
	return &IncrementalIndexer{indexNode, NewDirectoryIndexer(indexNode, threads), cache, seq, &sync.Mutex{}}
}
func (ii *IncrementalIndexer) RecordCount() int64 {
	//won't matter, just satisfy interface
	return 0
}
func (ii *IncrementalIndexer) flushCache() error {
	for _, rec := range ii.pathCache {
		if rec.fi == nil {
			pinfo := ii.indexNode.PathInfo(rec.path)
			pent := fsentity.NewPathEntity(nil, pinfo, ii.seq)
			ii.indexer.indexer.committer.AddEntity(pent)
		} else {
			//Scan 1 level, update itself and direct children
			//also pick up the new dirs
			fmt.Println("Update path:", rec.path)
			scanner := scan.NewDirectoryScanner(1)
			scanner.Scan(rec.path, ii)
		}
	}
	ii.pathCache = make([]*PathRecord, 0, cacheSize)
	return nil
}
func (ii *IncrementalIndexer) addChangePath(rec *PathRecord) error {
	ii.lock.Lock()
	defer ii.lock.Unlock()
	ii.pathCache = append(ii.pathCache, rec)
	if len(ii.pathCache) > maxCache {
		return filepath.SkipDir
	}
	return nil

}
func (ii *IncrementalIndexer) ProcessRow(rows *sql.Rows, cols []string) error {

	var rec PathRecord
	rec.fi = nil
	rows.Scan(&rec.path, &rec.pathmd5, &rec.mtime)
	rec.path = filepath.Join(ii.indexNode.BaseDir(), rec.path)
	//fmt.Println(rec.path)
	//Now check if mtime for this path changed.
	fi, err := os.Lstat(rec.path)
	if err != nil {
		if os.IsNotExist(err) {
			//mark for deletion, fi is nil
			return ii.addChangePath(&rec)
		}
		return nil
	}
	stat := fi.Sys().(*syscall.Stat_t)
	mtime, _, _ := utils.TimesFromStat_s(stat)
	if mtime != rec.mtime {
		rec.fi = &fi
		//fmt.Println("***Path changed:", rec.path, mtime, rec.mtime)
		return ii.addChangePath(&rec)
	}
	return nil
}

func (ii *IncrementalIndexer) Execute() {
	querier := query.NewSearchExecutor(ii.indexNode)
	querier.ExecuteQuery(pathQuery, true, ii)
	ii.flushCache()
	ii.indexer.Close()
}

func (ii *IncrementalIndexer) CommittedRecCnt() int64 {
	return ii.indexer.CommittedRecCnt()
}

//Implement call back interface
func (ii *IncrementalIndexer) ProcessRootDir(path string, name string) error {
	return ii.indexer.indexer.processDir(path, name)

}
func (ii *IncrementalIndexer) ProcessLeafDir(path string, name string) error {
	//Check if this is already in DB, if already in DB, skip
	//Otherwise, index this dir recursively
	pinfo := ii.indexNode.PathInfo(path)
	pathSql := fmt.Sprintf("select path, pmd5, lastmtime from path where pmd5='%s'", pinfo.PathMd5())
	dbFile := ii.indexNode.DBFileFromPathMd5(pinfo.PathMd5())
	//fmt.Println("dbFile***", dbFile, pathSql)
	rec := &PathRecord{"", "", nil, 0}
	query.Query(dbFile, pathSql, rec)
	if rec.mtime == 0 {
		//The path is not indexed
		fmt.Println("New Path detected:", path)
		//Update the dir non-recursively. but find the new subdirs
		scanner := scan.NewDirectoryScanner(0)
		scanner.Scan(path, ii.indexer)
	}
	return nil
}
func (ii *IncrementalIndexer) ProcessMidDir(path string, name string) error {
	//we should not have any Mid dir
	return ii.indexer.indexer.processDir(path, name)
}
func (ii *IncrementalIndexer) ProcessMidFile(path string) error {
	ii.indexer.indexer.processFile(path)
	return nil
}
