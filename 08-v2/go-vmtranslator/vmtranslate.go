package vmtranslate

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

func Run(root string) error {
	finfo, err := os.Stat(root)
	if err != nil {
		return xerrors.Errorf("Error while reading specified file or directory: %w", err)
	}

	var fnames []string
	var basename string
	if finfo.IsDir() {
		basename = finfo.Name()
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
		basename = strings.Replace(finfo.Name(), filepath.Ext(finfo.Name()), "", 1)
		fnames = append(fnames, root)
	}

	for _, fname := range fnames {
		ext := filepath.Ext(fname)
		if ext != ".vm" {
			continue
		}

		in, err := os.Open(fname)
		if err != nil {
			return xerrors.Errorf("Failed to open the input file: %w", err)
		}

		out, err := os.Create(strings.Replace(fname, ext, ".asm", 1))
		if err != nil {
			return xerrors.Errorf("Failed to create the output file: %w", err)
		}

		tr := NewTranslator(basename)

		if err := VMTranslate(in, out, tr); err != nil {
			return xerrors.Errorf("Failed to translate: %w", err)
		}

		in.Close()
		out.Close()
	}

	return nil
}

func VMTranslate(in io.Reader, out io.Writer, tr *Translator) error {
	s := bufio.NewScanner(in)
	w := bufio.NewWriter(out)

	var cmds []*Cmd
	for currentLine := 1; s.Scan(); currentLine++ {
		str := removeCommentAndSpaces(s.Text())
		if str == "" {
			continue
		}

		cmd, err := Parse(str)
		if err != nil {
			return xerrors.Errorf("Parse error on line %d: %w", currentLine, err)
		}

		cmds = append(cmds, cmd)
	}

	for _, cmd := range cmds {
		if _, err := w.WriteString(tr.Translate(cmd)); err != nil {
			return xerrors.Errorf("Error while writing to the file: %w", err)
		}
	}
	if _, err := w.WriteString(END); err != nil {
		return xerrors.Errorf("Error while writing to the file: %w", err)
	}
	if err := w.Flush(); err != nil {
		return xerrors.Errorf("Error on flushing the buffer: %w", err)
	}

	return nil
}

func removeCommentAndSpaces(s string) string {
	if i := strings.LastIndex(s, "//"); i != -1 { // remove commentouts
		s = s[:i]
	}
	s = strings.TrimSpace(s) // remove trailing spaces
	return s
}
