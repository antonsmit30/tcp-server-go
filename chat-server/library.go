package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type User struct {
	connection net.Conn
	username   string
}

type Room struct {
	users []User
}

type Message struct {
	sender User
	msg    string
}

func handleConnection(conn net.Conn, user User, messages chan Message) {
	defer conn.Close()
	fmt.Printf("Handling the connection %v\n", conn)
	_, err := io.WriteString(conn, "client>welcome to the room user: "+user.username+" \n>")
	if err != nil {
		fmt.Printf("Cannot write to client: %v", err)

	}
	for {
		fmt.Println("Start")

		// read in data returned to us and output to screen
		nr := bufio.NewReader(conn)
		ns := bufio.NewScanner(nr)
		for ns.Scan() {
			// messages <- ns.Text()
			// instead of a string channel we now have a Message channel
			messages <- Message{
				sender: user,
				msg:    ns.Text(),
			}
		}
		fmt.Println("End")
		{
			break
		}
	}
}
