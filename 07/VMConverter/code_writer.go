package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type Translator interface {
	SetNewFile(io.Writer)
	WriteArithmethic(string) error
	WritePushPop(commandType, string, int) error
	Close() error
}

type codewriter struct {
	*bufio.Writer
	closed bool
}

// 新しいTranslatorのインスタンスを生成する。
func NewCodeWriter() Translator {
	return &codewriter{}
}

// 新しいファイルに書き込む準備をして自身をオープンにする。
func (c *codewriter) SetNewFile(w io.Writer) {
	c.closed = false
	c.Writer = bufio.NewWriter(w)
}

// バッファに溜めた内容をファイルに書き込んで自身をクローズする。
func (c *codewriter) Close() error {
	if err := c.Flush(); err != nil {
		return err
	}
	if c.closed {
		return errors.New("already closed")
	}

	c.closed = true
	return nil
}

// 算術コマンドを書き込む。
func (c *codewriter) WriteArithmethic(instruction string) error {
	var txt string
	if instruction == "add" {
		txt = `// add
@SP
M=M-1
A=M
D=M
M=0
@SP
M=M-1
A=M
D=D+M
M=0
@SP
A=M
M=D
@SP
M=M+1`
	}
	_, err := fmt.Fprintln(c, txt)
	if err != nil {
		return err
	}

	return nil
}

// PUSHまたはPOPコマンドを書き込む。
func (c *codewriter) WritePushPop(instruction commandType, segment string, index int) error {
	return nil
}
