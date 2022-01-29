package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestConvertSimpleAdd(t *testing.T) {
	filename := "SimpleAdd"
	vmFile, _ := os.ReadFile(filepath.Join("fixtures", filename, filename+".vm"))
	in := bytes.NewBuffer(vmFile)
	got := bytes.NewBuffer([]byte{})
	asmFile, _ := os.ReadFile(filepath.Join("fixtures", filename, filename+".asm"))
	expected := bytes.NewBuffer(asmFile)
	err := convert(in, got, filename)
	if err != nil {
		t.Errorf("ERROR: %#v", err)
	} else if got.String() != expected.String() {
		t.Errorf("FAILED: expected %#v, got %#v", expected.String(), got.String())
	}
}

func TestConvertStackTest(t *testing.T) {
	filename := "StackTest"
	vmFile, _ := os.ReadFile(filepath.Join("fixtures", filename, filename+".vm"))
	in := bytes.NewBuffer(vmFile)
	got := bytes.NewBuffer([]byte{})
	asmFile, _ := os.ReadFile(filepath.Join("fixtures", filename, filename+".asm"))
	expected := bytes.NewBuffer(asmFile)
	err := convert(in, got, filename)
	if err != nil {
		t.Errorf("ERROR: %#v", err)
	} else if got.String() != expected.String() {
		t.Errorf("FAILED: expected %#v, got %#v", expected.String(), got.String())
	}
}

func TestConvertBasicTest(t *testing.T) {
	filename := "BasicTest"
	vmFile, _ := os.ReadFile(filepath.Join("fixtures", filename, filename+".vm"))
	in := bytes.NewBuffer(vmFile)
	got := bytes.NewBuffer([]byte{})
	asmFile, _ := os.ReadFile(filepath.Join("fixtures", filename, filename+".asm"))
	expected := bytes.NewBuffer(asmFile)
	err := convert(in, got, filename)
	if err != nil {
		t.Errorf("ERROR: %#v", err)
	} else if got.String() != expected.String() {
		t.Errorf("FAILED: expected %#v, got %#v", expected.String(), got.String())
	}
}

func TestConvertPointerTest(t *testing.T) {
	filename := "PointerTest"
	vmFile, _ := os.ReadFile(filepath.Join("fixtures", filename, filename+".vm"))
	in := bytes.NewBuffer(vmFile)
	got := bytes.NewBuffer([]byte{})
	asmFile, _ := os.ReadFile(filepath.Join("fixtures", filename, filename+".asm"))
	expected := bytes.NewBuffer(asmFile)
	err := convert(in, got, filename)
	if err != nil {
		t.Errorf("ERROR: %#v", err)
	} else if got.String() != expected.String() {
		t.Errorf("FAILED: expected %#v, got %#v", expected.String(), got.String())
	}
}
