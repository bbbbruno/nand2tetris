package main

import (
	"bufio"
	"fmt"
	"io"
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

	if err := Assemble(asmFile, hackFile); err != nil {
		log.Println(err)
		return
	}
}

func Assemble(r io.Reader, w io.Writer) error {
	p, t, b := NewParser(), NewTranslator(), bufio.NewWriter(w)

	bcmds, err := t.Translate(p.Parse(r))
	if err != nil {
		return err
	}

	for _, bcmd := range bcmds {
		if _, err := fmt.Fprintln(b, bcmd); err != nil {
			return err
		}
	}
	if err := b.Flush(); err != nil {
		return err
	}

	return nil
}
