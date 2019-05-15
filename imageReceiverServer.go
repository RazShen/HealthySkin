package main

import "net"
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing

const (
    CONN_HOST = "127.0.0.1"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)


func main() {
	fmt.Println("Launching server...")
	// listen on all interfaces
	ln, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	// accept connection on port
	conn, err := ln.Accept()
	handleClient(conn)
	if err != nil {
		fmt.Println("Error connecting to user:", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn *TCPConn) {
	// run loop forever (or until ctrl-c)
	for {
	// will listen for message to process ending in newline (\n)
	message, _ := bufio.NewReader(conn).ReadString('\n')
	// output message received
	fmt.Print("Message Received:", string(message))
	// sample process for string received
	newmessage := strings.ToUpper(message)
	// send new string back to client
	conn.Write([]byte(newmessage + "\n"))
	}
}