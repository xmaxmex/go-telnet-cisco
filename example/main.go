package main

import (
	"fmt"
	"os"
	"github.com/xmaxmex/go-telnet-cisco"
)

func main() {
	client := new(telnet.Client)

	err := client.Connect("192.168.10.1:23")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		return
	}

	err = client.Login("admin","admin")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		return
	}


	text,err := client.Cmd("show port")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		return
	}

	fmt.Println(text)

	//text,err = client.Cmd("show access-lists")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
	//	return
	//}
	//
	//fmt.Println(text)

	text,err = client.Cmd("show interface")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		return
	}

	fmt.Println(text)
}
