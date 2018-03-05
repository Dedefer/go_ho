package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func init() {

}

func recursiveTree(out io.Writer, path string, printFiles bool, offset string) error {
	path, err := filepath.Abs(path)
	fmt.Printf("%s%s\n", offset, path[len(filepath.Dir(path))+1:])
	if len(offset) == 0 {
		offset = "├── "
	} else {
		offset = "│   " + offset
	}
	if err != nil {
		return err
	}

	fileSlice, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range fileSlice {
		if file.IsDir() {
			recursiveTree(out, path+string(os.PathSeparator)+file.Name(), printFiles, offset)
		} else {
			if printFiles {
				out.Write([]byte(fmt.Sprintf("%s%s (%dB)\n", offset, file.Name(), file.Size())))
			}
		}
	}
	return nil

}

func dirTree(out io.Writer, path string, printFiles bool) error {
	_, err := ioutil.ReadDir(path)
	if err == nil {
		return recursiveTree(out, path, printFiles, "")
	}
	return err
}

func main() {

	var targetPath string
	var printFiles bool
	flag.StringVar(&targetPath, "p", ".", "target path")
	flag.BoolVar(&printFiles, "f", false, "print files")
	flag.Parse()

	err := dirTree(os.Stdout, targetPath, printFiles)
	if err != nil {
		fmt.Println(err.Error())
	}

}
