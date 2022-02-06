package engine

import (
	"fmt"
)

func (e *engine) writeArithmethic(cmd string) error {
	if _, err := fmt.Fprintln(e, cmd); err != nil {
		return err
	}

	return nil
}

func (e *engine) writePush(segment string, index int) error {
	text := fmt.Sprintf("push %s %d", segment, index)
	if _, err := fmt.Fprintln(e, text); err != nil {
		return err
	}

	return nil
}

func (e *engine) writePop(segment string, index int) error {
	text := fmt.Sprintf("pop %s %d", segment, index)
	if _, err := fmt.Fprintln(e, text); err != nil {
		return err
	}

	return nil
}

func (e *engine) writeLabel(label string) error {
	text := fmt.Sprintf("label %s", label)
	if _, err := fmt.Fprintln(e, text); err != nil {
		return err
	}

	return nil
}

func (e *engine) writeIf(label string) error {
	text := fmt.Sprintf("if-goto %s", label)
	if _, err := fmt.Fprintln(e, text); err != nil {
		return err
	}

	return nil
}

func (e *engine) writeGoto(label string) error {
	text := fmt.Sprintf("goto %s", label)
	if _, err := fmt.Fprintln(e, text); err != nil {
		return err
	}

	return nil
}

func (e *engine) writeFunction(receiver string, name string, nLocals int) error {
	text := fmt.Sprintf("function %s.%s %d", receiver, name, nLocals)
	if _, err := fmt.Fprintln(e, text); err != nil {
		return err
	}

	return nil
}

func (e *engine) writeCall(receiver string, name string, nArgs int) error {
	text := fmt.Sprintf("call %s.%s %d", receiver, name, nArgs)
	if _, err := fmt.Fprintln(e, text); err != nil {
		return err
	}

	return nil
}

func (e *engine) writeReturn() error {
	if _, err := fmt.Fprintln(e, "return"); err != nil {
		return err
	}

	return nil
}
