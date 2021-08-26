package main

import (
	"os"
	"fmt"
	"strings"
	"sshclient"
)


func main() {
	if len(os.Args) < 3 {
		fmt.Println(os.Args[0] + " nagios-host other-host")
		fmt.Println(os.Args)
		return
	}

	nagioshost := os.Args[1]
	otherhost := os.Args[2]

	keyrsa, _ := sshclient.GetKeyFile("id_rsa")
	output, err := sshclient.RunOneCommand(nagioshost, "cat /home/nagios/.ssh/id_rsa.pub", 5, keyrsa)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	output = strings.TrimSpace(output)
	sshclient.RunOneCommand(otherhost, "mkdir /home/nagios/.ssh && chmod 700 /home/nagios/.ssh", 5, keyrsa)
	sshclient.RunOneCommand(otherhost, "echo " + output + " >> /home/nagios/.ssh/authorized_keys && chmod 600 /home/nagios/.ssh/authorized_keys", 5, keyrsa)
	sshclient.RunOneCommand(otherhost, "chown -R nagios. /home/nagios/.ssh/", 5, keyrsa)

	fmt.Println("Added Key For Nagios")
}
