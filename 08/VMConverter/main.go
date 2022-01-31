package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"vmconverter/parser"
	"vmconverter/translator"
	"vmconverter/vmcommand"
)

func main() {
	root := os.Args[1]
	paths, dest, err := findAsmFiles(root)
	if err != nil {
		log.Println("error: ", err)
		return
	}

	if len(paths) == 0 {
		log.Println("no asm file found, skipping...")
		return
	}

	if err := convertEachFile(root, paths, dest); err != nil {
		log.Println("error: ", err)
		return
	}
}

func findAsmFiles(root string) (paths []string, dest string, err error) {
	finfo, err := os.Stat(root)
	if err != nil {
		return nil, "", err
	}

	if finfo.IsDir() {
		files, err := ioutil.ReadDir(root)
		if err != nil {
			return nil, "", err
		}

		dest = root
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".vm" {
				paths = append(paths, filepath.Join(dest, file.Name()))
			}
		}
	} else {
		dest = filepath.Dir(root)
		paths = append(paths, root)
	}

	return paths, dest, nil
}

func convertEachFile(root string, paths []string, dest string) error {
	filename := strings.Replace(filepath.Base(root), filepath.Ext(root), "", 1)
	asmFile, err := os.Create(filepath.Join(dest, filename+".asm"))
	if err != nil {
		return err
	}
	defer asmFile.Close()

	for _, path := range paths {
		vmFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer vmFile.Close()

		if err := convert(vmFile, asmFile, filename); err != nil {
			return err
		}
	}

	return nil
}

var cw = translator.NewCodeWriter()

func convert(r io.Reader, w io.Writer, filename string) error {
	p := parser.NewParser(&r)

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
