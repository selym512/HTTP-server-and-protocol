package main

import (
	"bufio"
	"net"
	"os"
)

func main() {

	serveraddr := "localhost:42069"
	udpAdd, err := net.ResolveUDPAddr("udp", serveraddr)
	if err != nil {
		os.Exit(1)
	}
	conn, connErr := net.DialUDP("udp", nil, udpAdd)
	if connErr != nil {
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		print("> ")
		readStd, err := reader.ReadString('\n')
		if err != nil {
			print("stdin nil")
		}
		conn.Write([]byte(readStd))
	}

}
