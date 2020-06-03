package main

import (
	"fmt"
	"net"
)

func Process(conn net.Conn) {
	defer conn.Close()
	for {
		var buf [4096]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("read from conn failed", err)
			return
		}

		fmt.Println(string(buf[:n]))
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Printf("listen failed")
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go Process(conn)
	}
}
