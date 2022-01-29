package main

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type commandType int

const (
	C_ARITHMETIC commandType = iota + 1
	C_PUSH
	C_POP
)

type Parser interface {
	HasMoreCommands() bool
	Advance() error
	CommandType() commandType
	Arg1() string
	Arg2() int
}

type parser struct {
	*bufio.Scanner
	currentCommand *command
	nextCommand    *command
}

type command struct {
	instruction, arg1, arg2 string
}

// 新しいParserのインスタンスを生成する。
func NewParser(file *io.Reader) *parser {
	scanner := bufio.NewScanner(*file)
	return &parser{scanner, nil, nil}
}

// 次のVMコマンドが存在するかどうかを判定する。
// WARNING: 怪しい
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

	p.currentCommand = p.nextCommand
	return nil
}

var instructionMap = map[string]commandType{
	"add": C_ARITHMETIC,
}

// 現在のVMコマンドの種類を返す。
func (p parser) CommandType() commandType {
	return instructionMap[p.currentCommand.instruction]
}

// 現在のVMコマンドの第一引数を返す。
// VMコマンドの種類が算術コマンド（C_ARITHMETIC）である場合はコマンド自体を返す。
func (p parser) Arg1() string {
	if p.CommandType() == C_ARITHMETIC {
		return p.currentCommand.instruction
	}

	return p.currentCommand.arg1
}

// 現在のVMコマンドの第二引数を返す。
func (p parser) Arg2() int {
	i, _ := strconv.Atoi(p.currentCommand.arg2)
	return i
}

// 与えられた文字列を解析し、commandオブジェクトを返す
func (p parser) parse(s string) *command {
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

	return &command{instruction, arg1, arg2}
}
