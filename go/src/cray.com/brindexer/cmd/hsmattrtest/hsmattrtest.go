package main

import (
	"cray.com/brindexer/lustre"
	_ "flag"
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("File path is required")
		return
	}
	fmt.Println(args[1])

	ha := lustre.GetHsmAttrs(args[1])
	fmt.Println("Hsm attributes")
	ha.Dump()

}
