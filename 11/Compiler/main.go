package main

import (
	"fmt"
	"os"

	"jackcompiler/compiler"
)

func main() {
	c := compiler.New(os.Args[1])
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %+v\n", err)
		return
	}
}
