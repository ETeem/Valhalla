package main

import (
	"os"
	"fmt"
	"time"
	"strings"
	"strconv"
	"sftpclient"
	"sshclient"
)


func main() {
	if len(os.Args) < 4 {
		fmt.Println(os.Args[0] + " nagios-host sudoers-host destination-host")
		fmt.Println(os.Args)
		return
	}

	nagioshost := os.Args[1]
	sudoershost := os.Args[2]
	desthost := os.Args[3]

	keyrsa, _ := sshclient.GetKeyFile("id_rsa")
	output, err := sshclient.RunOneCommand(nagioshost, "cat /home/nagios/.ssh/id_rsa.pub", 5, keyrsa)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	output = strings.TrimSpace(output)
	sshclient.RunOneCommand(desthost, "mkdir /home/nagios/.ssh && chmod 700 /home/nagios/.ssh", 5, keyrsa)
	sshclient.RunOneCommand(desthost, "echo " + output + " >> /home/nagios/.ssh/authorized_keys && chmod 600 /home/nagios/.ssh/authorized_keys", 5, keyrsa)
	sshclient.RunOneCommand(desthost, "chown -R nagios. /home/nagios/.ssh/", 5, keyrsa)

	fmt.Println("Added Key For Nagios")

	_, err = sftpclient.CopyFrom("/etc/sudoers", sudoershost, "/tmp/sudoers", 5, keyrsa)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	now := time.Now().Unix()
	strnow := strconv.FormatInt(now, 10)

	sshclient.RunOneCommand(desthost, "mv /etc/sudoers /etc/sudoers.valk." + strnow, 5, keyrsa)
	sftpclient.CopyFile("/tmp/sudoers", desthost, "/etc/sudoers", 5, keyrsa)

	fmt.Println("success - sudoers copied from " + sudoershost + " to " + desthost)

	output, _ = sshclient.RunOneCommand(desthost, "sed -i 's/inet_protocols = all/inet_protocols = ipv4/' /etc/postfix/main.cf && systemctl restart postfix", 5, keyrsa)
	fmt.Println("fixed postfix config: " + output)
}
