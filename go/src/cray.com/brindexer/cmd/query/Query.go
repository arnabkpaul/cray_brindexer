package main

import (
	"cray.com/brindexer/fsentity"
	"cray.com/brindexer/query"
	"cray.com/brindexer/scan"
	"flag"
	"fmt"
	"os"
	"time"
)

func printUsage() {
	fmt.Println("Execute query to all index DBs\n")
	fmt.Printf("Example: %s -q \"select %%s from %%s where size > 10000\"  TARGET_DIR ... \n\n", os.Args[0])

	fmt.Printf("Usage: %s [OPTIONS] TARGET_DIR ... \n\n", os.Args[0])
	flag.PrintDefaults()

}
func parseCommands() (string, string, bool, *string) {
	help := flag.Bool("h", false, "Print help")
	asis := flag.Bool("a", false, "Executed  query as is, if set the table, columns will not be altered.")
	//threads := flag.Int("t", 32, "Number of threads to run queries")
	querySql := flag.String("q", "", "Sql query in format \"select %s from %s where ...\"")
	dbBase := flag.String("index", "", "Alternative DB dir, if specified, the DB will be stored there.")

	flag.Parse()

	if len(*dbBase) == 0 {
		dbBase = nil
	}

	args := flag.Args()
	if *help || len(args) != 1 {
		printUsage()
		os.Exit(0)
	}

	stat, err := os.Lstat(args[0])
	if err != nil || !stat.IsDir() {
		fmt.Printf("Target directory does not exist: %s\n", args[0])
		os.Exit(0)
	}
	return args[0], *querySql, *asis, dbBase
}

func main() {
	target, querySql, asis, dbBase := parseCommands()
	start := time.Now()
	indexNode := fsentity.NewDefaultIndexNode(target, dbBase)

	level := indexNode.ReadLevel()
	//fmt.Println(level)
	printer := query.NewPrintQueryCallback(indexNode.BaseDir(), asis)
	search := query.NewLeveledSearchExecutor(&indexNode, querySql, level, asis, printer)

	scanner := scan.NewDirectoryScanner(level)
	scanner.Scan(target, search)
	search.Close()
	end := time.Now()
	elapsed := end.Sub(start)

	//fmt.Println("Total records found: ", printer.RecordCount())
	fmt.Println("Total records found: ", search.TotalRecordCnt())
	fmt.Println("\n")
	fmt.Println("Time elapsed: ", elapsed)
}
