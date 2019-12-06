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
	fmt.Println("Summarize a given directory indexed in DBs\n")

	fmt.Printf("Usage: %s [OPTIONS] TARGET_DIR ... \n", os.Args[0])
	flag.PrintDefaults()

}
func parseCommands() (string, *string) {
	help := flag.Bool("h", false, "Print help")
	//threads := flag.Int("t", 32, "Number of threads to run queries")
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
	return args[0], dbBase
}

func main() {
	target, dbBase := parseCommands()
	start := time.Now()
	indexNode := fsentity.NewDefaultIndexNode(target, dbBase)
	level := indexNode.ReadLevel()
	sumx := query.NewLeveledSummaryExecutor(&indexNode, level)
	//Level-1 is for python side compatibility
	scanner := scan.NewDirectoryScanner(level)
	scanner.Scan(target, sumx)
	sumx.Close()
	sum := sumx.GrandSummary()

	end := time.Now()
	elapsed := end.Sub(start)
	//fmt.Println(sum)
	fmt.Printf("Total file count:   %d\n", sum.FileCnt())
	fmt.Printf("Total link count:   %d\n", sum.LnkCnt())
	fmt.Printf("Total Dir count:    %d\n", sum.DirCnt())
	fmt.Printf("Total file size:    %d\n", sum.FileSize())
	fmt.Printf("Total File Objects: %d\n", sum.Total())
	fmt.Println("Time elapsed: ", elapsed)
}
