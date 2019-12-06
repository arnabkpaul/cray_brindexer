package DebugTools

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"
)

func printUsage() {
	fmt.Println("Scan a given directory and indexing in DBs\n")

	fmt.Printf("Usage: %s [OPTIONS] TARGET_DIR ... \n", os.Args[0])
	flag.PrintDefaults()

}
func parseCommands() (string, int, bool) {
	threads := flag.Int("t", 100, "Number of threads for scan")
	full := flag.Bool("f", false, "Force to do a full scan, default false")
	help := flag.Bool("h", false, "Print help")
	flag.Parse()

	args := flag.Args()
	if *help || len(args) != 1 {
		printUsage()
		os.Exit(0)
	}

	stat, _ := os.Lstat(args[0])
	if !stat.IsDir() {
		fmt.Printf("Target directory does not exist: %s\n", args[0])
		os.Exit(0)
	}
	return args[0], *threads, *full
}

func main() {

	start := time.Now()
	target := os.Args[0]
	fmt.Printf("read dir: %s\n", target)

	var stat syscall.Stat_t
	syscall.Lstat(target, &stat)

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Time eclpsed: ", elapsed)
}
