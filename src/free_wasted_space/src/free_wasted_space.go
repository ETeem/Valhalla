package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"os/exec"
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

	if strings.HasSuffix(current.FileName, ".log") {
		if err := os.Truncate(current.FileName, 0); err != nil {
			fmt.Println("Failed to truncate file: " + err.Error())
			return
		} 

		fmt.Println("Truncated File: " + current.FileName)
		return
	} else if strings.Contains(current.FileName, "/var/cache/yum") {
		fmt.Println("/var/cache/yum is in the list... running yum clean all")
		out, err := exec.Command("yum", "clean", "all").Output()
		if err != nil {
			fmt.Println("Error cleaning yum cache: " + err.Error())
			return
		} else {
			fmt.Println(out)
			return
		}
	} else {
		if current.FileName == "/var/log/messages" {
			if err := os.Truncate("/var/log/messages", 0); err != nil {
				fmt.Println("Failed to truncate file: " + err.Error())
				return
			}
		}	 

		fmt.Println("Truncated File: " + current.FileName)
		return
	}

	return //fileList, nil
}
