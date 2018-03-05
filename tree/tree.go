package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// функция фильтрации директорий
func onlyDirs(files []os.FileInfo) []os.FileInfo {
	result := make([]os.FileInfo, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			result = append(result, file)
		}
	}
	return result
}

func recursiveTree(out io.Writer, path string, printFiles bool, offset string, lastOffset string) error {
	// Получаем полный путь
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// выводим название текущей папки с учетом текущего отступа
	fmt.Println(offset + lastOffset + path[len(filepath.Dir(path))+len(string(os.PathSeparator)):])

	// наращиваем отступ для подкаталогов в зависимости от последнего отступа
	// (последний элемент в списке или нет)
	if lastOffset == "├───" {
		offset += "│   "
	} else if lastOffset == "└───" {
		offset += "    "
	}

	// получаем все каталоги в директории
	fileSlice, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// фильтруем их если не установлен флаг показывать файлы
	if !printFiles {
		fileSlice = onlyDirs(fileSlice)
	}

	numberOfLastFile := len(fileSlice) - 1

	// запускаем цикл по файлам
	for i, file := range fileSlice {

		// устанавливаем следующий последний отступ в зависимости от того
		// последний файл в директории или нет
		if i == numberOfLastFile {
			lastOffset = "└───"
		} else {
			lastOffset = "├───"
		}

		// если файл является директорией то
		// рекурсивно вызываем функцию со следующим отступом
		// если нет то печатаем имя файла
		if file.IsDir() {
			recursiveTree(out, path+string(os.PathSeparator)+file.Name(), printFiles, offset, lastOffset)
		} else if printFiles {
			out.Write([]byte(fmt.Sprintf("%s (%dB)\n", offset+lastOffset+file.Name(), file.Size())))
		}

	}
	return nil

}

func dirTree(out io.Writer, path string, printFiles bool) error {
	// находим полный путь до заданной директории
	// и если нет ошибок то вызываем recursiveTree
	_, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	return recursiveTree(out, path, printFiles, "", "")
}

func main() {
	// создаем переменные для флагов и парсим флаги
	targetPath := ""
	printFiles := false
	flag.StringVar(&targetPath, "p", ".", "target path")
	flag.BoolVar(&printFiles, "f", false, "print files")
	flag.Parse()

	// выводим ошибку если она была
	err := dirTree(os.Stdout, targetPath, printFiles)
	if err != nil {
		fmt.Println(err.Error())
	}
}
