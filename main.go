package main

import (
	"os"
	"fmt"
	"golang.org/x/crypto/ssh"
)

const (
	cParamsNum = 7 // include the executable file, user name, password, ip addr1, volume path on ip addr1, ip addr2, volume path on ip addr2.
	cSSHPort = ":22"
)

func main() {
	if len(os.Args) != cParamsNum {
		fmt.Println("Parameters are invalid. Please input like DeleteGFSVolume.exe user-name password ip1 volume-path ip2 volume-path.")
		os.Exit(-1)
	}

	var host1 = Client{
		userName: os.Args[1],
		password: os.Args[2],
		ipAddr: os.Args[3] + cSSHPort,
		volume: os.Args[4],
	}

	var host2 = Client{
		userName: os.Args[1],
		password: os.Args[2],
		ipAddr: os.Args[5] + cSSHPort,
		volume: os.Args[6],
	}

	if err := host1.deleteVolume(); err != nil {
		fmt.Println(host1.userName, host1.password, host1.ipAddr, host1.volume, " failed to delete the specific volume.", err)
		os.Exit(-1)
	}

	if err := host2.deleteVolume(); err != nil {
		fmt.Println(host2.userName, host2.password, host2.ipAddr, host2.volume, " failed to delete the specific volume.", err)
		os.Exit(-1)
	}
}


type Client struct {
	userName string
	password string
	ipAddr string
	volume string
}

func (c Client) deleteVolume() error {
	pw := []ssh.AuthMethod{ssh.Password(c.password)}
	conf := ssh.ClientConfig{
		User: c.userName,
		Auth: pw,
	}

	client, err := ssh.Dial("tcp", c.ipAddr, &conf)
	if err != nil {
		fmt.Println("failed to dial the ip ", c.ipAddr, " with", c.userName, c.password, c.volume, ".", err)
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Println("failed to create a new session for ip ", c.ipAddr, " with", c.userName, c.password, c.volume, ".", err)
		return err
	}
	defer session.Close()

	var cmd [2]string
	cmd[0] = fmt.Sprintf("setfattr -x trusted.glusterfs.volume-id %s", c.volume)
	cmd[1] = fmt.Sprintf("rm -rf %s", c.volume)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	//fmt.Printf("run cmd 1: \"%s\"\n", cmd[0])
	//fmt.Printf("run cmd 2: \"%s\"\n", cmd[1])
	if err = session.Run(cmd[0] + "; " + cmd[1]); err != nil {
		fmt.Println("failed to run cmd:", cmd[0], "; ", cmd[1], err)
		return err
	}

	return nil
}