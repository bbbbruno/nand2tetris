package main

import (
	"fmt"
	"os"

	"jackcompiler/compiler"
)

func main() {
	c := compiler.New(os.Args[1])
	if err := c.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
		return
	}
}
