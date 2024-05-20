package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	recursiveTree(out, path, 0, printFiles, "")
	return nil
}

func recursiveTree(out io.Writer, path string, depth int, printFiles bool, prefix string) error {
	files, _ := getFileList(path, printFiles)
	lastElementIndex := len(files) - 1

	for idx, file := range files {
		fmt.Fprintln(out, prefix+getFileString(file, idx == lastElementIndex))
		if file.IsDir() {
			pref := prefix
			if idx == lastElementIndex {
				pref += "\t"
			} else {
				pref += "│\t"
			}
			recursiveTree(out, path+string(os.PathSeparator)+file.Name(), depth+1, printFiles, pref)
		}
	}

	return nil
}

func getFileList(path string, printFiles bool) ([]os.FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err.Error())
		}
	}()

	filesInfo, _ := file.Readdir(-1)

	files := make([]os.FileInfo, 0)
	for _, element := range filesInfo {
		if element.IsDir() {
			files = append(files, element)
		}
		if !(element.IsDir()) && printFiles {
			files = append(files, element)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	return files, nil
}

func getFileString(file os.FileInfo, last bool) string {
	symbol := '├'
	if last {
		symbol = '└'
	}
	result := string(symbol) + "───" + file.Name()
	if !file.IsDir() {
		if file.Size() == 0 {
			result += " (empty)"
		} else {
			result += " (" + strconv.FormatInt(file.Size(), 10) + "b)"
		}

	}
	return result
}
