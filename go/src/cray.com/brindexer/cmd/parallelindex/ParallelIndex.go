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
func parseCommands() (string, int, int, int, *string, bool) {
	help := flag.Bool("h", false, "Print help")
	threads := flag.Int("t", 300, "Number of threads per process")
	procs := flag.Int("p", 8, "Number of scan processes")
	level := flag.Int("l", 1, "Level of directory to build index[0..4], 0 means 1 index at root")
	full := flag.Bool("f", false, "Force to perform a full scan, default false")
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
	if *level < 0 || *level > 4 {
		fmt.Printf("Level must be between 0 and 4: %s\n", args[0])
		os.Exit(0)
	}
	stat, err := os.Lstat(args[0])
	if err != nil || !stat.IsDir() {
		fmt.Printf("Target directory does not exist: %s\n", args[0])
		os.Exit(0)
	}
	return args[0], *threads, *procs, *level, dbBase, *full
}

func main() {
	target, threads, procs, level, dbBase, full := parseCommands()
	start := time.Now()
	indexNode := fsentity.NewDefaultIndexNode(target, dbBase)

	lastScan := indexNode.ReadLastScan()
	if lastScan != -1 {
		//There is already existing index, the level has to stay the same.
		savedLevel := indexNode.ReadLevel()
		if savedLevel != -1 && level != savedLevel {
			fmt.Println("Level specified is different from saved level, last saved level will be used")
			level = savedLevel
		}
	}
	indexNode.SaveLevel(level)
	pindexer := indexing.NewParallelIndexer(&indexNode, threads, procs, full)
	//We scan each directory at (level) within a separate go-routine
	scanner := scan.NewDirectoryScanner(level)
	fmt.Println("Start to index:", target)
	scanner.Scan(target, pindexer)
	pindexer.Close()
	indexNode.SaveLastScan()
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("All processes completed, time elapsed: ", elapsed)
}
