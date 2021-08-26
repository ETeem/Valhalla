package main

import (
	"fmt"
//	"os"
	"strings"
	"os/exec"
)

type ProcessInfo struct {
	Owner string
	Pid string
	Cpu string
	Mem string
	Program string
}

func Get_CPU_And_Memory(tproc *ProcessInfo) {
	cmd := exec.Command("/usr/bin/top", "-b", "-p", tproc.Pid, "-n", "1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(out))
		return
	}

	output := string(out)

	lines := strings.Split(output, "\n")
	p1 := lines[7]
	p1 = strings.Join(strings.Fields(p1), " ")
	parts := strings.Split(p1, " ")

	tproc.Cpu = parts[8]
	tproc.Mem = parts[9]

	return
}

func main() {
	var proc []ProcessInfo

        cmd := "ps aux | sort -nrk 3,3 | head -n 2"
        out, err := exec.Command("bash", "-c", cmd).Output()
        if err != nil {
		fmt.Println("Failed To Execute Command: " + err.Error())
		return
        }

	output := string(out)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if len(line) < 1 {
			continue
		}

		tmpproc := ProcessInfo{}
		line = strings.Join(strings.Fields(line), " ")
		parts := strings.Split(line, " ")
		tmpproc.Owner = parts[0]
		tmpproc.Pid = parts[1]
		tmpproc.Program = parts[10]

		Get_CPU_And_Memory(&tmpproc)

		proc = append(proc, tmpproc)
	}

	var stitch string
	for _, p := range proc {
		stitch += "Pid: " + p.Pid + " | Owner: " + p.Owner + " | Program: "  + p.Program + " | CPU: " + p.Cpu + " | Memory: " + p.Mem + "&"
	}

	stitch = stitch[:len(stitch) - 1]

	fmt.Println(stitch)
}

