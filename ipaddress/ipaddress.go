package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	conn, err := net.Dial("udp", "50.18.192.251:80") // duckduckgo.com IP address
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	fmt.Print(localAddr[0:idx])
}
