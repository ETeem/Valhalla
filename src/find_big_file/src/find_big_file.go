package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"path/filepath"
)

type FileSize struct {
	FileName string
	FileSize int64
	IsDir bool
}

func main() {
	var files []FileSize
	var current FileSize
	var second FileSize
	var third FileSize

	var excludelist []string

	if len(os.Args) < 2 {
		fmt.Println("Missing Path")
		return
	}

	mnts, err := os.Open("/proc/mounts")
	if err != nil {
		fmt.Println("Unable To Check Mount Points")
		return
	}
	defer mnts.Close()

	scanner := bufio.NewScanner(mnts)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		expath := parts[1]
		mtype := parts[2]
		if mtype == "nfs" {
			excludelist = append(excludelist, expath)
		}
	}
	
	searchDir := os.Args[1] 
	fileList := make([]string, 0)

	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		flag := false
		for _, ex := range excludelist {
			if strings.Contains(path, ex) {
				flag = true
			}
		}

		if ! flag {
			fileList = append(fileList, path)
		}

		return err
	})
	
	if e != nil {
		panic(e)
	}

	for _, file := range fileList {
		tmpfsize := FileSize{}
		fi, err := os.Lstat(file)
		if err != nil {
			fmt.Println(err.Error())
		}
		if fi.Mode() & os.ModeSymlink == os.ModeSymlink {
			continue
		}

		tmpfsize.FileName = file
		tmpfsize.FileSize = fi.Size()
		tmpfsize.IsDir = fi.IsDir()
		files = append(files, tmpfsize)
	}

	for _, file := range files {
		if file.FileSize > current.FileSize {
			current = file
		}
	}

	for _, file := range files {
		if file.FileName != current.FileName {
			if file.FileSize > second.FileSize {
				second = file
			}
		}
	}

	for _, file := range files {
		if file.FileName != current.FileName && file.FileName != second.FileName {
			if file.FileSize > third.FileSize {
				third = file
			}
		}
	}

	csize := current.FileSize / 1048576
	ssize := second.FileSize / 1048576
	tsize := third.FileSize / 1048576

	if csize > 1024 {
		csize = csize / 1024
		fmt.Printf("%s - %dGB\n", current.FileName, csize)
	} else {
		fmt.Printf("%s - %dMB\n", current.FileName, csize)
	}

	if ssize > 1024 {
		ssize = ssize / 1024
		fmt.Printf("%s - %dGB\n", second.FileName, ssize)
	} else {
		fmt.Printf("%s - %dMB\n", second.FileName, ssize)
	}

	if tsize > 1024 {
		tsize = tsize / 1024
		fmt.Printf("%s - %dGB\n", third.FileName, tsize)
	} else {
		fmt.Printf("%s - %dMB\n", third.FileName, tsize)
	}

	return //fileList, nil
}
