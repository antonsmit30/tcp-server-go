package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// Function to return our local date and time,
// useful for our messages to have some form of
// date and timestamp
// (Taken from docs: https://pkg.go.dev/time#Time.Format)
// The layout string used by the Parse function and Format method
// shows by example how the reference time should be represented.
// We stress that one must show how the reference time is formatted,
// not a time of the user's choosing. Thus each layout string is a
// representation of the time stamp,
//
//	Jan 2 15:04:05 2006 MST
//
// An easy way to remember this value is that it holds, when presented
// in this order, the values (lined up with the elements above):
//
//	1 2  3  4  5    6  -7
func returnCurrentTime() string {
	now := time.Now().Local()
	log.Printf("date time: %v", now.String())
	return string(now.Day())
}

func main() {

	// channel for messages
	// Making a buffer of 100 messages for concurrent messages to be sent
	// messages := make(chan string, 100)
	messages := make(chan Message, 100)

	// setup a room :)
	red := Room{}

	// Setup a listener to listen for TCP connections
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to establish a connection: %v", err)
	}

	// Utilizing logger for stdout
	logger := log.Default()
	logger.Println("Weave online")

	// connection setups
	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Printf("Could not establish a connection: %v", err)
		}

		// generate a username for our new connection
		// u := uuid.NewString()
		// New logic, ask user for username :p
		go func() {

			var u string

			_, err = io.WriteString(conn, returnCurrentTime()+" : Give yourself a username: ")
			if err != nil {
				log.Println("Could not send message to gather username")
			}
			// setup a temp reader to collect input from user
			ns := bufio.NewScanner(bufio.NewReader(conn))
			for ns.Scan() {
				u = ns.Text()
				log.Printf("User has selected: %v\n", u)
				break
			}

			// ok so we might need to make the logic to add the user details
			// part of this goroutine as well
			// lets add the connection to our room
			us := User{username: u, connection: conn}
			red.users = append(red.users, us)
			fmt.Printf("Users: %v", red.users)

			go handleConnection(conn, us, messages)
			fmt.Println("exiting handle")

		}()
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
					_, err := io.WriteString(u.connection, msg.sender.username+" > "+msg.msg+"\n>")
					if err != nil {
						fmt.Printf("Cannot send to client: %v", err)
					}
				}

			}

		}()

	}
}
