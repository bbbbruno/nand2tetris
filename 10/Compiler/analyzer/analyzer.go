package analyzer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"compiler/engine"
	"compiler/tokenizer"
)

type Analyzer interface {
	Run() error
}

type analyzer struct {
	source string
	name   string
}

func New(source string) Analyzer {
	a := &analyzer{}
	a.source = source
	a.name = strings.Replace(filepath.Base(a.source), filepath.Ext(a.source), "", 1)
	return a
}

func (a *analyzer) Run() error {
	paths, err := findPaths(a.source)
	if err != nil {
		return err
	}
	dest, err := getDestPath(a.source)
	if err != nil {
		return err
	}

	outFile, err := os.Create(filepath.Join(dest, a.name+"T"+".xml"))
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		tokenizer := tokenizer.New(file)
		engine := engine.New(tokenizer, outFile)
		if err := engine.Compile(); err != nil {
			return err
		}
	}

	return nil
}

// 与えられたソースから全ての.jackファイルのパスを取得する。
// ディレクトリならその中のファイル全ての、ファイルならそのファイル自身のパスを探索する。
func findPaths(source string) (paths []string, err error) {
	finfo, err := os.Stat(source)
	if err != nil {
		return nil, err
	}

	if finfo.IsDir() {
		files, err := ioutil.ReadDir(source)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".jack" {
				paths = append(paths, filepath.Join(source, file.Name()))
			}
		}
	} else {
		paths = append(paths, source)
	}

	return paths, nil
}

// 与えられたソースがディレクトリかファイルかを判定して出力先ディレクトリパスを返す。
func getDestPath(source string) (dest string, err error) {
	finfo, err := os.Stat(source)
	if err != nil {
		return "", err
	}

	if finfo.IsDir() {
		return source, nil
	} else {
		return filepath.Dir(source), nil
	}
}
