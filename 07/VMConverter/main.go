package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	filename := os.Args[1]
	ext := filepath.Ext(filename)
	if ext != ".vm" {
		log.Println("error: not vm file")
		return
	}

	vmFile, err := os.Open(filename)
	if err != nil {
		log.Println("error: ", err)
		return
	}

	asmFile, err := os.Create(strings.Replace(filename, ext, "", 1) + ".asm")
	if err != nil {
		log.Println("error: ", err)
		return
	}
	defer asmFile.Close()

	if err := convert(vmFile, asmFile); err != nil {
		log.Println("error: ", err)
		return
	}
}

func convert(r io.Reader, w io.Writer) error {
	p, cw := NewParser(&r), NewCodeWriter()

	cw.SetNewFile(w)
	for p.HasMoreCommands() {
		if err := p.Advance(); err != nil {
			return err
		}
		if p.CommandType() == C_ARITHMETIC {
			if err := cw.WriteArithmethic(p.Arg1()); err != nil {
				return err
			}
		} else {
			if err := cw.WritePushPop(p.CommandType(), p.Arg1(), p.Arg2()); err != nil {
				return err
			}
		}
	}
	if err := p.Scanner.Err(); err != nil {
		return err
	}
	if err := cw.Close(); err != nil {
		return err
	}

	return nil
}
