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
	var dir, basename string
	if finfo.IsDir() {
		dir = root
		basename = finfo.Name()
		files, err := os.ReadDir(root)
		if err != nil {
			return xerrors.Errorf("Error while reading directory %s: %w", root, err)
		}
		for _, file := range files {
			ext := filepath.Ext(file.Name())
			if !file.IsDir() && ext == ".vm" {
				fnames = append(fnames, filepath.Join(root, file.Name()))
			}
		}
	} else {
		dir = filepath.Dir(root)
		basename = strings.Replace(finfo.Name(), filepath.Ext(finfo.Name()), "", 1)
		fnames = append(fnames, root)
	}

	out, err := os.Create(filepath.Join(dir, basename+".asm"))
	if err != nil {
		return xerrors.Errorf("Failed to create the output file: %w", err)
	}
	defer out.Close()

	if err := VMTranslate(fnames, out); err != nil {
		return xerrors.Errorf("Failed to translate: %w", err)
	}

	return nil
}

func VMTranslate(fnames []string, out io.Writer) error {
	w := bufio.NewWriter(out)
	if _, err := w.WriteString(initAssembly()); err != nil {
		return xerrors.Errorf("Error while writing to the file: %w", err)
	}

	for _, fname := range fnames {
		tr := NewTranslator(strings.Replace(filepath.Base(fname), filepath.Ext(fname), "", 1))

		in, err := os.Open(fname)
		if err != nil {
			return xerrors.Errorf("Failed to open the input file: %w", err)
		}

		output, err := ParseTranslate(in, tr)
		if err != nil {
			return xerrors.Errorf("Failed to parse and translate the file: %w", err)
		}

		if _, err := w.WriteString(output); err != nil {
			return xerrors.Errorf("Error while writing to the file: %w", err)
		}
		if err := w.Flush(); err != nil {
			return xerrors.Errorf("Error on flushing the buffer: %w", err)
		}

		in.Close()
	}

	return nil
}

func ParseTranslate(in io.Reader, tr *Translator) (output string, err error) {
	s := bufio.NewScanner(in)

	var cmds []*Cmd
	for currentLine := 1; s.Scan(); currentLine++ {
		line := removeCommentAndSpaces(s.Text())
		if line == "" {
			continue
		}

		cmd, err := Parse(line)
		if err != nil {
			return "", xerrors.Errorf("Parse error on line %d: %w", currentLine, err)
		}

		cmds = append(cmds, cmd)
	}

	for _, cmd := range cmds {
		output += tr.Translate(cmd)
	}

	return output, nil
}

func removeCommentAndSpaces(s string) string {
	if i := strings.LastIndex(s, "//"); i != -1 { // remove commentouts
		s = s[:i]
	}
	s = strings.TrimSpace(s) // remove trailing spaces
	return s
}
