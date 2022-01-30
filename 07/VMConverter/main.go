package main

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"vmconverter/parser"
	"vmconverter/translator"
)

func main() {
	var paths []string
	filepath.WalkDir(os.Args[1], func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".vm" {
			paths = append(paths, path)
		}

		return nil
	})

	for _, path := range paths {
		filename := strings.Replace(filepath.Base(path), ".vm", "", 1)
		vmFile, err := os.Open(path)
		if err != nil {
			log.Println("error: ", err)
			return
		}

		asmFile, err := os.Create(strings.Replace(path, ".vm", "", 1) + ".asm")
		if err != nil {
			log.Println("error: ", err)
			return
		}
		defer asmFile.Close()

		if err := convert(vmFile, asmFile, filename); err != nil {
			log.Println("error: ", err)
			return
		}
	}
}

func convert(r io.Reader, w io.Writer, filename string) error {
	p, cw := parser.NewParser(&r), translator.NewCodeWriter()

	cw.SetNewFile(w, filename)
	for p.HasMoreCommands() {
		if err := p.Advance(); err != nil {
			return err
		}
		if p.IsArithmetic() {
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
