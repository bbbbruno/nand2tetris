package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"vmconverter/parser"
	"vmconverter/translator"
)

// TODO：ディレクトリを指定できるようにする。ディレクトリ内の全てのファイルを変換する。
func main() {
	path := os.Args[1]
	base := filepath.Base(path)
	ext := filepath.Ext(path)
	filename := strings.Replace(base, ext, "", 1)
	if ext != ".vm" {
		log.Println("error: not vm file")
		return
	}

	vmFile, err := os.Open(path)
	if err != nil {
		log.Println("error: ", err)
		return
	}

	asmFile, err := os.Create(strings.Replace(path, ext, "", 1) + ".asm")
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
