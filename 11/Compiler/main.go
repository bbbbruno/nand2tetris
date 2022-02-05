package main

import (
	"fmt"
	"os"

	"compiler/analyzer"
)

func main() {
	analyzer := analyzer.New(os.Args[1])
	if err := analyzer.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
		return
	}
}
