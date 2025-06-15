package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("client!")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Could not setup client connection: %v", err)
	}

	go func() {

		ns := bufio.NewScanner(bufio.NewReader(conn))
		for ns.Scan() {
			fmt.Printf("message received: %v", ns.Text())
		}

	}()

}
