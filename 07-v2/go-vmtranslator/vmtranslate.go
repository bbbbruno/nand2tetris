package vmtranslate

import (
	"bufio"
	"io"
	"strings"

	"golang.org/x/xerrors"
)

func VMTranslate(in io.Reader, out io.Writer) error {
	s := bufio.NewScanner(in)
	w := bufio.NewWriter(out)

	var cmds []Cmd
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
		if _, err := w.WriteString(cmd.Translate()); err != nil {
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
