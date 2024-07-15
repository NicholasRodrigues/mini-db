package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	var authEnabled bool
	var useTLS bool
	var caCertPath string

	flag.BoolVar(&authEnabled, "auth", false, "Enable authentication with token")
	flag.BoolVar(&useTLS, "tls", false, "Enable TLS")
	flag.StringVar(&caCertPath, "ca-cert", "./client.pem", "Path to CA certificate")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: client -auth=<true|false> -tls=<true|false> -ca-cert=<path/to/client.pem> <address> <port>")
		return
	}

	address := args[0]
	port := args[1]
	addr := fmt.Sprintf("%s:%s", address, port)

	var conn net.Conn
	var err error
	if useTLS {
		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			log.Fatalf("Failed to load CA certificate: %v", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			RootCAs:            caCertPool,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		}
		conn, err = tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			log.Fatalf("Failed to connect to server with TLS: %v", err)
		}
	} else {
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
	}
	defer func() {
		if conn != nil {
			err := conn.Close()
			if err != nil {
				log.Printf("Failed to close connection: %v", err)
			}
		}
	}()

	fmt.Println("Connected to server at", addr)
	if authEnabled {
		fmt.Println("Authentication is enabled. Enter commands with auth token.")
	} else {
		fmt.Println("Authentication is not enabled. Enter commands without auth token.")
	}
	fmt.Println("Enter commands (SET [<auth_token>] <key> <value> or LOOKUP [<auth_token>] <key>)")
	fmt.Println("Type 'EXIT' to terminate connection.")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if cmd == "EXIT" {
			fmt.Println("Exiting...")
			break
		}

		if authEnabled {
			parts := strings.Fields(cmd)
			if (parts[0] == "SET" && len(parts) != 4) || (parts[0] == "LOOKUP" && len(parts) != 3) {
				fmt.Println("Invalid command format. Use SET or LOOKUP with an auth token.")
				continue
			}
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
