package main

import (
	"fmt"
	"log"
	"main.go/cmd/internal/request"
	"net"
)

const port = ":42069"

func main() {

	tcpListener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer tcpListener.Close()
	for {
		con, err := tcpListener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		//fmt.Printf("Connection to %s successful\n", con.RemoteAddr())
		req, err := request.RequestFromReader(con)
		fmt.Println("Request line:")
		fmt.Println("- Method: " + req.RequestLine.Method)
		fmt.Println("- Target: " + req.RequestLine.RequestTarget)
		fmt.Println("- Version: " + req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Printf("- Body:\n%s", string(req.Body))
	}
}
