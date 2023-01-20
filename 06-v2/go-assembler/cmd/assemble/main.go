package main

import (
	"assemble"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	fname := os.Args[1]
	ext := filepath.Ext(fname)
	if ext != ".asm" {
		return xerrors.Errorf("%s is not asm file", fname)
	}

	in, err := os.Open(fname)
	if err != nil {
		return xerrors.Errorf("failed to open the input file: %w", err)
	}
	defer in.Close()

	out, err := os.Create(strings.Replace(fname, ext, ".hack", 1))
	if err != nil {
		return xerrors.Errorf("failed to create the output file: %w", err)
	}
	defer out.Close()

	if err := assemble.Assemble(in, out); err != nil {
		return xerrors.Errorf("failed to assemble: %w", err)
	}

	return nil
}
