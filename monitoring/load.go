package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

func main() {
	address := "150.165.202.101:8080"
	numClients := 100
	numRequests := 1000

	var wg sync.WaitGroup
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := net.Dial("tcp", address)
			if err != nil {
				fmt.Printf("Failed to connect: %v\n", err)
				return
			}
			defer conn.Close()

			for j := 0; j < numRequests; j++ {
				fmt.Fprintf(conn, "SET test_key%d test_value%d\n", j, j)
				response, _ := bufio.NewReader(conn).ReadString('\n')
				if response != "OK\n" {
					fmt.Printf("Unexpected response: %s\n", response)
				}
			}
		}()
	}
	wg.Wait()
}
