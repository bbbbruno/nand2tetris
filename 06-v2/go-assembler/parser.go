package assemble

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type parseError struct {
	line int
	err  error
}

func (e parseError) Error() string {
	return fmt.Sprintf("Parse error on line %d: %s", e.line, e.err.Error())
}

func (e parseError) UnWrap() error {
	return e.err
}

func ParseLabels(in io.Reader) ([]string, *SymTable, error) {
	lines := []string{}
	s := bufio.NewScanner(in)
	st := NewSymTable()
	assembledLineNumber := 0

	for currentLineNumber := 1; s.Scan(); currentLineNumber++ {
		str := removeCommentAndSpaces(s.Text())
		lines = append(lines, str)

		// when the line is comment or newline, don't count the line
		if str == "" {
			continue
		}
		// When (XXX), the next line number is added to the symbol table, and don't count the line.
		if isLCmd(str) {
			sym := str[1 : len(str)-1]
			if !IsValidSymbol(sym) {
				return nil, nil, &parseError{line: currentLineNumber, err: &InvalidSymbolError{symbol: sym}}
			}
			st.AddSymbolWithAddr(sym, assembledLineNumber)
			continue
		}

		assembledLineNumber++
	}

	return lines, st, nil
}

func ParseLines(lines []string, st *SymTable) (cmds []Cmd, err error) {
	for i, line := range lines {
		// don't parse empty line or L command.
		if line == "" || isLCmd(line) {
			continue
		}

		// parse A command
		if isACmd(line) {
			cmd, err := NewACmd(line, st)
			if err != nil {
				return nil, &parseError{line: i + 1, err: err}
			}
			cmds = append(cmds, cmd)
			continue
		}

		// parse C command
		cmd, err := NewCCmd(line)
		if err != nil {
			return nil, &parseError{line: i + 1, err: err}
		}
		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

func removeCommentAndSpaces(s string) string {
	if i := strings.LastIndex(s, "//"); i != -1 { // remove commentouts
		s = s[:i]
	}
	s = strings.TrimSpace(s) // remove trailing spaces
	return s
}

func isACmd(s string) bool {
	return s[0] == '@'
}

func isLCmd(s string) bool {
	return s[0] == '(' && s[len(s)-1] == ')'
}
