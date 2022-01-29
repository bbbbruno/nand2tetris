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

func NewCodeWriter() Translator {
	return &codewriter{}
}

func (c *codewriter) SetNewFile(w io.Writer) {
	c.closed = false
	c.Writer = bufio.NewWriter(w)
}

func (c *codewriter) Close() error {
	if _, err := fmt.Fprint(c, end()); err != nil {
		return err
	}
	if err := c.Flush(); err != nil {
		return err
	}
	if c.closed {
		return errors.New("already closed")
	}

	c.closed = true
	return nil
}

var instructionAsmMap = map[string]string{
	"add": pop("M") + pop("D+M") + push(),
}

// 算術コマンドを書き込む。
func (c *codewriter) WriteArithmethic(instruction string) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	text := instructionAsmMap[instruction]
	_, err := fmt.Fprint(c, text)
	return err
}

// PUSHまたはPOPコマンドを書き込む。
func (c *codewriter) WritePushPop(instruction commandType, segment string, index int) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	switch instruction {
	case C_PUSH:
		return c.writePush(segment, index)
	case C_POP:
		return c.writePop(segment, index)
	default:
		return nil
	}
}

func (c *codewriter) checkClosed() error {
	if c.closed {
		return errors.New("writer is closed")
	} else {
		return nil
	}
}

var segmentMap = map[string]func(int) string{
	"constant": func(index int) string { return constant(index) + push() },
}

func (c *codewriter) writePush(segment string, index int) error {
	text := segmentMap[segment](index)
	_, err := fmt.Fprint(c, text)
	return err
}

func (c *codewriter) writePop(segment string, index int) error {
	return nil
}

func push() string {
	return `@SP
A=M
M=D
@SP
M=M+1
`
}

func pop(comp string) string {
	return fmt.Sprintf(`@SP
M=M-1
A=M
D=%s
M=0
`, comp)
}

func constant(index int) string {
	return fmt.Sprintf(`@%d
D=A
`, index)
}

func end() string {
	return `(END)
@END
0;JMP
`
}
