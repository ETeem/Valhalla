package main

import (
	"fmt"
	"os"
	"time"
	"strings"
	"os/exec"
)

func main() {
	os.Chdir("/usr/local/nagiosxi/scripts")

        if len(os.Args) < 2 {
                fmt.Println("Missing Path")
                return
        }

	host := strings.TrimSpace(os.Args[1])

	// Delete Services
	output, err := exec.Command("/usr/bin/php", "./nagiosql_delete_service.php", "--config=" + host).Output()
	if err != nil {
		fmt.Println("error deleting services: - " + host + " - " + string(output) + " - " + err.Error())
	}

	// Reconfigure
	time.Sleep(1 * time.Minute)
	output, err = exec.Command("/usr/bin/php", "./reconfigure_nagios.sh").Output()
	if err != nil {
		fmt.Println("error reconfiguring Nagios: " + string(output) + " - " + err.Error())
	}

	// Delete Host
	time.Sleep(1 * time.Minute)
	output, err = exec.Command("/usr/bin/php", "./nagiosql_delete_host.php", "--host=" + host).Output()
	if err != nil {
		fmt.Println("error deleting host: - " + host + " - " + string(output) + " - " + err.Error())
	}

	// Reconfigure
	time.Sleep(1 * time.Minute)
	output, err = exec.Command("/usr/bin/php", "./reconfigure_nagios.sh").Output()
	if err != nil {
		fmt.Println("error reconfiguring Nagios: " + string(output) + " - " + err.Error())
	}

	fmt.Println("success")
	return 
}
