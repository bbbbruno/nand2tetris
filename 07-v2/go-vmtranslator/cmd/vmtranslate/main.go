package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"vmtranslate"

	"golang.org/x/xerrors"
)

func main() {
	if err := run(); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	root := os.Args[1]
	finfo, err := os.Stat(root)
	if err != nil {
		return xerrors.Errorf("Error while reading specified file or directory: %w", err)
	}

	var fnames []string
	if finfo.IsDir() {
		files, err := os.ReadDir(root)
		if err != nil {
			return xerrors.Errorf("Error while reading directory %s: %w", root, err)
		}
		for _, file := range files {
			if !file.IsDir() {
				fnames = append(fnames, filepath.Join(root, file.Name()))
			}
		}
	} else {
		fnames = append(fnames, root)
	}

	for _, fname := range fnames {
		ext := filepath.Ext(fname)
		if ext != ".vm" {
			return xerrors.Errorf("%s is not vm file", fname)
		}

		in, err := os.Open(fname)
		if err != nil {
			return xerrors.Errorf("Failed to open the input file: %w", err)
		}

		out, err := os.Create(strings.Replace(fname, ext, ".asm", 1))
		if err != nil {
			return xerrors.Errorf("Failed to create the output file: %w", err)
		}

		if err := vmtranslate.VMTranslate(in, out); err != nil {
			return xerrors.Errorf("Failed to translate: %w", err)
		}

		in.Close()
		out.Close()
	}

	return nil
}
