package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	filename := os.Args[1]
	ext := filepath.Ext(filename)
	if ext != ".asm" {
		log.Println("error: not asm file")
		return
	}

	asmFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer asmFile.Close()

	hackFile, err := os.Create(strings.Replace(filename, ext, "", 1) + ".hack")
	if err != nil {
		log.Println(err)
		return
	}
	defer hackFile.Close()

	p, t, b := NewParser(), NewTranslator(), bufio.NewWriter(hackFile)
	cmds := p.Parse(asmFile)
	bcmds, err := t.Translate(cmds)
	if err != nil {
		log.Println(err)
		return
	}

	for _, bcmd := range bcmds {
		if _, err := fmt.Fprintln(b, bcmd); err != nil {
			log.Println(err)
			return
		}
	}
	b.Flush()
}
