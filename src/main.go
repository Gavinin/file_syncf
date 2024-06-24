package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const TargetDir = "/Applications"
const Maxdepth = 2

type DirRecord struct {
	path string
	dept int
}

var dirRecords = make([]DirRecord, 0)

func main() {
	if runtime.GOOS != "darwin" {
		fmt.Println("This program only works on macOS")
		return
	}
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
		return
	}
	exec(currentPath, 1)
	for _, s := range dirRecords {
		if s.dept <= Maxdepth {
			exec(s.path, s.dept)
		}
	}

}

func exec(currentPath string, dep int) {
	fileStream, err := os.Open(currentPath)
	if err != nil {
		return
	}
	defer fileStream.Close()
	filesArray, _ := fileStream.Readdirnames(0)
	for _, fileName := range filesArray {
		abs, err := filepath.Abs(fmt.Sprintf("%s/%s", currentPath, fileName))
		if err != nil {
			fmt.Println(err)
			return
		}
		info, err := os.Stat(abs)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !info.IsDir() {
			continue
		}

		if strings.HasSuffix(fileName, ".app") {
			_, err = os.Stat(fmt.Sprintf("%s/%s", TargetDir, fileName))
			if err != nil {
				if os.IsNotExist(err) {
					ln(abs, fmt.Sprintf("%s/%s", TargetDir, fileName))
				}
			}
		} else {
			if dep <= Maxdepth {
				dirRecords = append(dirRecords, DirRecord{
					path: abs,
					dept: dep + 1,
				})
			}
		}
	}
}

func ln(abs string, target string) {
	err := os.Symlink(abs, target)
	if err != nil {
		fmt.Println(err)
		return
	}
}
