package main

import (
	"io"
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
	paths, dest, err := FindVmFiles(root)
	if err != nil {
		log.Println("error: ", err)
		return
	}

	if len(paths) == 0 {
		log.Println("no vm file found, skipping...")
		return
	}

	dirname := strings.Replace(filepath.Base(root), filepath.Ext(root), "", 1)
	destFile, err := os.Create(filepath.Join(dest, dirname+".asm"))
	if err != nil {
		log.Println("error: ", err)
		return
	}
	defer destFile.Close()

	if err := VMConvert(paths, destFile); err != nil {
		log.Println("error: ", err)
		return
	}
}

// 指定されたディレクトリ以下の.vmファイルを見つけて、そのパスの配列を返す。
// 指定されたものがファイルであれば、そのファイルのパスを一つだけ格納した配列を返す。
func FindVmFiles(root string) (paths []string, dest string, err error) {
	finfo, err := os.Stat(root)
	if err != nil {
		return nil, "", err
	}

	if finfo.IsDir() { // ディレクトリなら.asmファイルのパスを全て格納
		files, err := os.ReadDir(root)
		if err != nil {
			return nil, "", err
		}

		dest = root
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".vm" {
				if err != nil {
					return nil, "", err
				}

				paths = append(paths, filepath.Join(root, file.Name()))
			}
		}
	} else { // ファイルならそのファイルのパスのみ格納する
		if filepath.Ext(root) == ".vm" {
			dest = filepath.Dir(root)
			paths = append(paths, root)
		}
	}

	return paths, dest, nil
}

// 与えられたパスのファイルそれぞれに対して変換を行い、一つのファイルにまとめて指定のディレクトリに.asmファイルとして出力する。
func VMConvert(paths []string, w io.Writer) error {
	cw := translator.NewCodeWriter(w)

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		filename := strings.Replace(filepath.Base(file.Name()), filepath.Ext(file.Name()), "", 1)
		p := parser.NewParser(file)
		cw.SetFileName(filename)

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
					if err := cw.WriteLabel(p.Arg1(), p.CurrentFuncName); err != nil {
						return err
					}
				case vmcommand.C_GOTO:
					if err := cw.WriteGoto(p.Arg1(), p.CurrentFuncName); err != nil {
						return err
					}
				case vmcommand.C_IF:
					if err := cw.WriteIf(p.Arg1(), p.CurrentFuncName); err != nil {
						return err
					}
				case vmcommand.C_FUNCTION:
					if err := cw.WriteFunction(p.Arg1(), p.Arg2()); err != nil {
						return err
					}
				case vmcommand.C_RETURN:
					if err := cw.WriteReturn(); err != nil {
						return err
					}
				case vmcommand.C_CALL:
					if err := cw.WriteCall(p.Arg1(), p.Arg2()); err != nil {
						return err
					}
				}
			}
		}
		if err := p.Scanner.Err(); err != nil {
			return err
		}
	}

	if err := cw.Close(); err != nil {
		return err
	}

	return nil
}
