package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: client <address> <port> [--tls]")
		return
	}

	address := os.Args[1]
	port := os.Args[2]
	useTLS := len(os.Args) == 4 && os.Args[3] == "--tls"
	addr := fmt.Sprintf("%s:%s", address, port)

	var conn net.Conn
	var err error
	if useTLS {
		conn, err = tls.Dial("tcp", addr, &tls.Config{
			InsecureSkipVerify: true,
		})
	} else {
		conn, err = net.Dial("tcp", addr)
	}
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server at", addr)
	fmt.Println("Enter commands (SET [<auth_token>] <key> <value> or LOOKUP [<auth_token>] <key>)")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if cmd == "EXIT" {
			fmt.Println("Exiting...")
			break
		}

		_, err := conn.Write([]byte(cmd + "\n"))
		if err != nil {
			log.Printf("Failed to send command: %v", err)
			continue
		}

		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Failed to read response: %v", err)
			continue
		}
		fmt.Println(strings.TrimSpace(resp))
	}
}
