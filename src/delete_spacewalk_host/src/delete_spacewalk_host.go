package main

import (
	"fmt"
	"os"
	"strings"
	"os/exec"
)

func main() {
	os.Chdir("/usr/local/nagiosxi/scripts")

        if len(os.Args) < 4 {
                fmt.Println("Missing Parameter: <username> <password> <host_to_delete>")
                return
        }

	user := strings.TrimSpace(os.Args[1])
	pass := strings.TrimSpace(os.Args[2])
	host := strings.TrimSpace(os.Args[3])

	output, err := exec.Command("/usr/bin/spacecmd", "-y", "-u", user, "-p", pass, "--", "system_delete", host).Output()
	if err != nil {
		fmt.Println("error deleting Host: - " + host + " - " + string(output) + " - " + err.Error())
		return
	}


	fmt.Println("success")
	return 
}

