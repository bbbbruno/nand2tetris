package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	hackFile, err := os.Create(filename[:len(filename)-len(ext)] + ".hack")
	if err != nil {
		log.Println(err)
		return
	}
	defer hackFile.Close()

	p, t, b, c := NewParser(), NewCode(), bufio.NewWriter(hackFile), bufio.NewScanner(asmFile)
	for c.Scan() {
		cmd := p.Parse(c.Text())
		if cmd == nil {
			continue
		}

		str, err := t.Translate(cmd)
		if err != nil {
			log.Println(err)
			return
		}

		if _, err := fmt.Fprintln(b, str); err != nil {
			log.Println(err)
			return
		}
	}
	if err := c.Err(); err != nil {
		log.Println(err)
		return
	}

	b.Flush()
}
