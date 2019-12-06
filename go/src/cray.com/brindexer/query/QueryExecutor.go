package query

import (
	"database/sql"
	"fmt"
	"github.com/karrick/godirwalk"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"sync"
)

const queueSize = 1000
const maxQueryWorker = 32

type QueryCallBack interface {
	ProcessRow(rows *sql.Rows, cols []string) error
	RecordCount() int64
}

func Query(dbFile string, query string, callback QueryCallBack) error {

	db, err := sql.Open("sqlite3", dbFile+"?_timeout=10000")
	if err != nil {
		fmt.Println(err, dbFile)
		return err
	}
	defer db.Close()

	//fmt.Println("Executing sql:", query, dbFile)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err, dbFile)
		return err
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	for rows.Next() {
		//fmt.Println("row:", rows)
		if callback.ProcessRow(rows, cols) != nil {
			break
		}
	}
	return rows.Err()

}

func queryWorker(jobs <-chan string, query string, callback QueryCallBack, wg *sync.WaitGroup) {
	for dbFile := range jobs {
		//fmt.Println(dbFile, query)
		Query(dbFile, query, callback)
	}
	wg.Done()
}

func QueryAll(dbDir string, query string, threads int, callback QueryCallBack) error {
	//Cap at 32, some systems default file descriptor is too low, 32 is a reasonable
	//Number in our case. Unless we can set higher FD number from somewhere.
	if threads > maxQueryWorker {
		threads = maxQueryWorker
	}
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(threads)
	jobs := make(chan string, queueSize)
	for w := 1; w <= threads; w++ {
		go queryWorker(jobs, query, callback, &waitGroup)
	}

	godirwalk.Walk(dbDir, &godirwalk.Options{
		Unsorted:            true,
		FollowSymbolicLinks: false,
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			//fmt.Println("Query Entry:", osPathname)
			if de.IsRegular() && strings.HasSuffix(de.Name(), ".db") {
				//Query(osPathname, query, callback)
				jobs <- osPathname
			}

			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})

	close(jobs)
	waitGroup.Wait()
	return nil

}
