package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestConvertSimpleAdd(t *testing.T) {
	vmFile, _ := os.ReadFile(filepath.Join("fixtures", "SimpleAdd", "SimpleAdd.vm"))
	in := bytes.NewBuffer(vmFile)
	got := bytes.NewBuffer([]byte{})
	asmFile, _ := os.ReadFile(filepath.Join("fixtures", "SimpleAdd", "SimpleAdd.asm"))
	expected := bytes.NewBuffer(asmFile)
	err := convert(in, got)
	if err != nil {
		t.Errorf("ERROR: %#v", err)
	} else if got.String() != expected.String() {
		t.Errorf("FAILED: expected %#v, got %#v", expected.String(), got.String())
	}
}

func TestConvertStackTest(t *testing.T) {
	vmFile, _ := os.ReadFile(filepath.Join("fixtures", "StackTest", "StackTest.vm"))
	in := bytes.NewBuffer(vmFile)
	got := bytes.NewBuffer([]byte{})
	asmFile, _ := os.ReadFile(filepath.Join("fixtures", "StackTest", "StackTest.asm"))
	expected := bytes.NewBuffer(asmFile)
	err := convert(in, got)
	if err != nil {
		t.Errorf("ERROR: %#v", err)
	} else if got.String() != expected.String() {
		t.Errorf("FAILED: expected %#v, got %#v", expected.String(), got.String())
	}
}

func TestConvertBasicTest(t *testing.T) {
	vmFile, _ := os.ReadFile(filepath.Join("fixtures", "BasicTest", "BasicTest.vm"))
	in := bytes.NewBuffer(vmFile)
	got := bytes.NewBuffer([]byte{})
	asmFile, _ := os.ReadFile(filepath.Join("fixtures", "BasicTest", "BasicTest.asm"))
	expected := bytes.NewBuffer(asmFile)
	err := convert(in, got)
	if err != nil {
		t.Errorf("ERROR: %#v", err)
	} else if got.String() != expected.String() {
		t.Errorf("FAILED: expected %#v, got %#v", expected.String(), got.String())
	}
}

func TestConvertPointerTest(t *testing.T) {
	vmFile, _ := os.ReadFile(filepath.Join("fixtures", "PointerTest", "PointerTest.vm"))
	in := bytes.NewBuffer(vmFile)
	got := bytes.NewBuffer([]byte{})
	asmFile, _ := os.ReadFile(filepath.Join("fixtures", "PointerTest", "PointerTest.asm"))
	expected := bytes.NewBuffer(asmFile)
	err := convert(in, got)
	if err != nil {
		t.Errorf("ERROR: %#v", err)
	} else if got.String() != expected.String() {
		t.Errorf("FAILED: expected %#v, got %#v", expected.String(), got.String())
	}
}
