package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	uuid "github.com/google/uuid"
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

func main() {

	// channel for messages
	// Making a buffer of 100 messages for concurrent messages to be sent
	// messages := make(chan string, 100)
	messages := make(chan Message, 100)

	// setup a room :)
	red := Room{}

	fmt.Println("Hello World!")

	// Setup a listener to listen for TCP connections
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to establish a connection: %v", err)
	}

	// connection setups
	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Printf("Could not establish a connection: %v", err)
		}

		// generate a username for our new connection
		u := uuid.NewString()

		// lets add the connection to our room
		us := User{username: u, connection: conn}
		red.users = append(red.users, us)
		fmt.Printf("Users: %v", red.users)

		go handleConnection(conn, us, messages)

		fmt.Println("exiting handle")

		// continuously loop and check for messages
		// ideally, we want to put our connections in a room :)
		// we use a nameless func here to ensure that we split this
		// off from the main Accept() that needs to happen here
		go func() {

			for msg := range messages {
				fmt.Println(msg.sender.username)
				// we also want to loop through our connections in our room and broadcast.
				// Ideally though, this username is the username we are broadcasting to...
				// we also want to know who is the sender :)
				for _, u := range red.users {
					_, err := io.WriteString(u.connection, msg.sender.username+">"+msg.msg+"\n>")
					if err != nil {
						fmt.Printf("Cannot send to client: %v", err)
					}
				}

			}

		}()

	}
}
