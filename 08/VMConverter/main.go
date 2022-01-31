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
	"vmconverter/vmcommand"
)

func main() {
	paths, err := findAsmFiles(os.Args[1])
	if err != nil {
		log.Println("error: ", err)
		return
	}

	if err := convertEachFile(paths); err != nil {
		log.Println("error: ", err)
		return
	}
}

func findAsmFiles(root string) (paths []string, err error) {
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
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
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func convertEachFile(paths []string) error {
	for _, path := range paths {
		vmFile, err := os.Open(path)
		if err != nil {
			return err
		}

		asmFile, err := os.Create(strings.Replace(path, ".vm", "", 1) + ".asm")
		if err != nil {
			return err
		}
		defer asmFile.Close()

		filename := strings.Replace(filepath.Base(path), filepath.Ext(path), "", 1)
		if err := convert(vmFile, asmFile, filename); err != nil {
			return err
		}
	}

	return nil
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
			switch p.CommandType() {
			case vmcommand.C_POP, vmcommand.C_PUSH:
				if err := cw.WritePushPop(p.CommandType(), p.Arg1(), p.Arg2()); err != nil {
					return err
				}
			case vmcommand.C_LABEL:
				if err := cw.WriteLabel(p.Arg1()); err != nil {
					return err
				}
			case vmcommand.C_GOTO:
				if err := cw.WriteGoto(p.Arg1()); err != nil {
					return err
				}
			case vmcommand.C_IF:
				if err := cw.WriteIf(p.Arg1()); err != nil {
					return err
				}
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
