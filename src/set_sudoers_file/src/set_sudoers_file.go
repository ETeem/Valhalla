package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"sftpclient"
	"sshclient"
)


func main() {
	if len(os.Args) < 3 {
		fmt.Println(os.Args[0] + "sudoers-source destination-host")
		return
	}

	srchst := os.Args[1]
	dsthst := os.Args[2]

	keyrsa, _ := sshclient.GetKeyFile("id_rsa")
	_, err := sftpclient.CopyFrom("/etc/sudoers", srchst, "/tmp/sudoers", 5, keyrsa)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	now := time.Now().Unix()
	strnow := strconv.FormatInt(now, 10)

	sshclient.RunOneCommand(dsthst, "mv /etc/sudoers /etc/sudoers.valk." + strnow, 5, keyrsa)
	sftpclient.CopyFile("/tmp/sudoers", dsthst, "/etc/sudoers", 5, keyrsa)

	fmt.Println("success - sudoers copied from " + srchst + " to " + dsthst)
}
