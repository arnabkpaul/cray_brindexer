package db

import (
	"cray.com/brindexer/fsentity"
	"fmt"
	"strings"
	"time"
)

type CommitTask struct {
	records []fsentity.FSEntity
	dbConn  *SqliteWrapper
}

func (t CommitTask) commit() error {
	//fmt.Println(" CommitTask, commit:", t.records.Len())
	start := time.Now()

	db := t.dbConn.Connection()
	defer db.Close()
	var err, reterr error

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("CommitTask, failed to start transaction")
		return err
	}

	for _, rec := range t.records {

		sqlstr := rec.CommitSqls()
		//fmt.Println(sqlstr)
		sqls := strings.Split(sqlstr, ";")

		for _, sql := range sqls {
			_, err := tx.Exec(sql)
			if err != nil {
				fmt.Println("CommitTask, failed to insert:", err)
				fmt.Println(sql)
				reterr = err
			}
		}
	}

	err = tx.Commit()

	t.dbConn.addRecordCnt(len(t.records))
	end := time.Now()
	elapsed := end.Sub(start)
	if elapsed.Seconds() > 2.0 {
		fmt.Println("***CommitTask takes too long:", elapsed, len(t.records))
	}
	return reterr

}

func commitWorker(index int, jobs <-chan *CommitTask, w *DBWriter) {
	for task := range jobs {
		err := task.commit()
		if err != nil {
			fmt.Println("Commit failed, requeue task")
			w.taskQueue <- task
		}
	}
	w.close()
	w.waitGroup.Done()
	//fmt.Println("*****worker", index, "done")
}
