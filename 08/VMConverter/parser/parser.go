package parser

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"vmconverter/vmcommand"
)

type Parser interface {
	HasMoreCommands() bool
	Advance() error
}

type currentCommand vmcommand.VMCommand
type nextCommand vmcommand.VMCommand

type parser struct {
	*bufio.Scanner
	currentCommand
	nextCommand     nextCommand
	CurrentFuncName string
}

func NewParser(file *io.Reader) *parser {
	scanner := bufio.NewScanner(*file)
	return &parser{scanner, nil, nil, ""}
}

// 次のVMコマンドが存在するかどうかを判定する。
func (p *parser) HasMoreCommands() bool {
	if !p.Scan() {
		p.nextCommand = nil
		return false
	}

	command := p.parse(p.Text())
	if command == nil {
		return p.HasMoreCommands()
	}

	p.nextCommand = command
	return true
}

// 次のVMコマンドを読み、それを現在のVMコマンドとする。
func (p *parser) Advance() error {
	if p.nextCommand == nil {
		return errors.New("no more commands")
	}
	if p.nextCommand.CommandType() == vmcommand.C_FUNCTION {
		p.CurrentFuncName = p.nextCommand.Arg1()
	}

	p.currentCommand = p.nextCommand
	return nil
}

// 与えられた文字列を解析してVMCommand型のインスタンスを生成する。
func (p *parser) parse(s string) vmcommand.VMCommand {
	if i := strings.LastIndex(s, "//"); i != -1 { // コメントアウトを除去
		s = s[:i]
	}
	if s == "" {
		return nil
	}

	ss := strings.Split(s, " ")
	var instruction, arg1, arg2 string
	instruction = ss[0]
	if len(ss) > 1 {
		arg1 = ss[1]
	}
	if len(ss) > 2 {
		arg2 = ss[2]
	}

	return vmcommand.NewVMCommand(instruction, arg1, arg2)
}
