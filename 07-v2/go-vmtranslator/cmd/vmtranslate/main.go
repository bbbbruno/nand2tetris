package main

import (
	"log"
	"os"
	"vmtranslate"
)

func main() {
	if err := vmtranslate.Run(os.Args[1]); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
