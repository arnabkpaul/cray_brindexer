package main

import (
	"cray.com/brindexer/indexing"
	"cray.com/brindexer/scan"
	"flag"
	"fmt"
	"os"
	"time"
)

func printUsage() {
	fmt.Println("Scan a given directory and stat each file entry\n")

	fmt.Printf("Usage: %s [OPTIONS] TARGET_DIR ... \n", os.Args[0])
	flag.PrintDefaults()

}
func parseCommands() (string, int, bool, bool) {
	help := flag.Bool("h", false, "Print help")
	threads := flag.Int("t", 400, "Number of threads for scan")
	statf := flag.Bool("s", false, "Stat each file entry")
	xattr := flag.Bool("x", false, "Fetch extedned attributes file entry")

	flag.Parse()

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

	return args[0], *threads, *statf, *xattr
}

func main() {
	target, threads, statf, xattr := parseCommands()
	start := time.Now()

	fmt.Println("Start Null index:", target)
	indexer := indexing.NewNullIndexer(threads, statf, xattr)
	scanner := scan.NewDirectoryScanner(0)
	fscount := scanner.Scan(target, indexer)
	indexer.Close()

	fmt.Println("File system objects: ", fscount)
	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Time elapsed: ", elapsed)
	totalSec := int64(elapsed.Seconds()) + 1
	fmt.Println("Index rate (obj/sec): ", fscount/totalSec)
}
