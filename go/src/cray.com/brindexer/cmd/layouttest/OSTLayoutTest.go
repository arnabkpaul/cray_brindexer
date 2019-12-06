package main

import (
	"cray.com/brindexer/lustre"
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("File path is required")
		return
	}
	fmt.Println("File path:", args[1])

	lo := lustre.GetLayout(args[1])
	if lo == nil {
		fmt.Println("Failed to get layout")
		return
	}
	lo.Dump()

}
