package query

import (
	"cray.com/brindexer/fsentity"
	"fmt"
)

const queryTable string = "entries_0 as f join path as p on f.pathmd5=p.pmd5"
const queryCols string = "path, name, size, mtime, pathmd5"

type SearchExecutor struct {
	indexNode fsentity.IndexNode
}

func NewSearchExecutor(indexNode fsentity.IndexNode) *SearchExecutor {
	return &SearchExecutor{indexNode}
}

func (s *SearchExecutor) ExecuteQuery(query string, asis bool, callback QueryCallBack) {

	if !asis {
		query = fmt.Sprintf(query, queryCols, queryTable)
	}
	//fmt.Println(query)

	dbDir := s.indexNode.DBDir()
	QueryAll(dbDir, query, s.indexNode.DBCount(), callback)
}

//Note: query must be in format "select %s from %s where ...."
func (s *SearchExecutor) Execute(query string, asis bool, callback QueryCallBack) int64 {
	s.ExecuteQuery(query, asis, callback)
	return callback.RecordCount()

}
