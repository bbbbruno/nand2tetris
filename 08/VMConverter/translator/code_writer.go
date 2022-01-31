package translator

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"vmconverter/vmcommand"
)

type Translator interface {
	SetNewFile(io.Writer, string)
	WriteInit() error
	WriteArithmethic(string) error
	WritePushPop(vmcommand.CommandType, string, int) error
	WriteLabel(string, string) error
	WriteGoto(string, string) error
	WriteIf(string, string) error
	WriteFunction(string, int) error
	WriteReturn() error
	WriteCall(string, int) error
	Close() error
}

type codewriter struct {
	*bufio.Writer
	filename string
	closed   bool
}

func NewCodeWriter() Translator {
	return &codewriter{}
}

func (c *codewriter) SetNewFile(w io.Writer, filename string) {
	c.filename = filename
	c.closed = false
	c.Writer = bufio.NewWriter(w)
	c.WriteInit()
}

func (c *codewriter) WriteInit() error {
	text := setSP() + callFunc("Sys.init", 0)
	_, err := fmt.Fprint(c, text)
	return err
}

var operationMap = map[string]string{
	"add": operateDouble("M+D"),
	"sub": operateDouble("M-D"),
	"neg": operateSingle("-M"),
	"and": operateDouble("M&D"),
	"or":  operateDouble("M|D"),
	"not": operateSingle("!M"),
}
var comparisonMap = map[string]string{
	"eq": "JEQ",
	"gt": "JGT",
	"lt": "JLT",
}

// 算術コマンドを書き込む。
func (c *codewriter) WriteArithmethic(instruction string) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	var text string
	switch instruction {
	case "eq", "gt", "lt":
		text = operateCompare(comparisonMap[instruction])
	default:
		text = operationMap[instruction]
	}

	_, err := fmt.Fprint(c, text)
	return err
}

// PUSHまたはPOPコマンドを書き込む。
func (c *codewriter) WritePushPop(instruction vmcommand.CommandType, segment string, index int) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	switch instruction {
	case vmcommand.C_PUSH:
		return c.writePush(segment, index)
	case vmcommand.C_POP:
		return c.writePop(segment, index)
	default:
		return nil
	}
}

func (c *codewriter) WriteLabel(label string, currentFuncName string) error {
	text := defineLabel(label, currentFuncName)
	_, err := fmt.Fprint(c, text)
	return err
}

func (c *codewriter) WriteGoto(label string, currentFuncName string) error {
	text := goTo(label, currentFuncName)
	_, err := fmt.Fprint(c, text)
	return err
}

func (c *codewriter) WriteIf(label string, currentFuncName string) error {
	text := pop("M") + ifGoTo(label, currentFuncName)
	_, err := fmt.Fprint(c, text)
	return err
}

func (c *codewriter) WriteFunction(name string, numlocals int) error {
	text := defineFunc(name) + strings.Repeat(push("0"), numlocals)
	_, err := fmt.Fprint(c, text)
	return err
}

func (c *codewriter) WriteReturn() error {
	text := returnFunc()
	_, err := fmt.Fprint(c, text)
	return err
}

func (c *codewriter) WriteCall(name string, numargs int) error {
	text := callFunc(name, numargs)
	_, err := fmt.Fprint(c, text)
	return err
}

func (c *codewriter) checkClosed() error {
	if c.closed {
		return errors.New("writer is closed")
	} else {
		return nil
	}
}

func (c *codewriter) writePush(segment string, index int) error {
	var text string
	switch segment {
	case "constant":
		text = getConstant(index) + push("D")
	case "static":
		text = getStatic(c.filename, index) + push("D")
	default:
		text = memoryPush(segment, index)
	}
	_, err := fmt.Fprint(c, text)
	return err
}

func (c *codewriter) writePop(segment string, index int) error {
	var text string
	switch segment {
	case "static":
		text = pop("M") + setStatic(c.filename, index)
	default:
		text = memoryPop(segment, index)
	}
	_, err := fmt.Fprint(c, text)
	return err
}

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
