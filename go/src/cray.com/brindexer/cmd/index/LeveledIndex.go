package main

import (
	"cray.com/brindexer/fsentity"
	"cray.com/brindexer/indexing"
	"cray.com/brindexer/scan"
	"flag"
	"fmt"
	"os"
	"time"
)

func printUsage() {
	fmt.Println("Scan a given directory and indexing in DBs\n")

	fmt.Printf("Usage: %s [OPTIONS] TARGET_DIR ... \n", os.Args[0])
	flag.PrintDefaults()

}
func parseCommands() (string, int, *string, bool, int) {
	help := flag.Bool("h", false, "Print help")
	threads := flag.Int("t", 400, "Number of threads for scan")
	full := flag.Bool("f", false, "Force to do a full scan, default false")
	dbBase := flag.String("index", "", "Alternative DB dir, if specified, the DB will be stored there.")
	level := flag.Int("l", 1, "Level of directory to scan in separate process.")

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
	return args[0], *threads, dbBase, *full, *level
}

func main() {
	target, threads, dbBase, full, level := parseCommands()
	start := time.Now()
	indexNode := fsentity.NewDefaultIndexNode(target, dbBase)

	//We scan each leaf directory within a separate go-routine
	lastScan := indexNode.ReadLastScan()
	dbcount := int64(0)
	savedLevel := indexNode.ReadLevel()
	if savedLevel > 0 {
		//Note: This in general should not be mised with parallel scan
		//This is to make sure that in case of user error, we don't create
		//Confusion. the default 8 processes is used.
		fmt.Println("The directory is indexed  with level:", level)
		fmt.Println("The last saved level will be used for the index")
		pindexer := indexing.NewParallelIndexer(&indexNode, threads, 8, full)
		//We scan each directory at (level) within a separate go-routine
		scanner := scan.NewDirectoryScanner(level)

		scanner.Scan(target, pindexer)
		pindexer.Close()
	} else if full || lastScan < 0 {
		lindexer := indexing.NewLeveledIndexer(&indexNode, threads)
		scanner := scan.NewDirectoryScanner(level)
		fmt.Println("Start full index:", target)
		scanner.Scan(target, lindexer)
		lindexer.Close()
		dbcount = lindexer.CommitedRecCnt()

	} else {
		fmt.Println("Start incremental index:", target)
		indexer := indexing.NewIncrementalIndexer(&indexNode, threads)
		indexer.Execute()
		dbcount = indexer.CommittedRecCnt()
	}
	indexNode.SaveLevel(0)
	indexNode.SaveLastScan()
	end := time.Now()
	elapsed := end.Sub(start)
	//fmt.Println("File system objects: ", fscount)
	fmt.Println("Commited db records: ", dbcount)
	totalSec := int64(elapsed.Seconds()) + 1
	fmt.Println("Index rate (obj/sec): ", dbcount/totalSec)
	fmt.Println("Time elapsed: ", elapsed)
}
