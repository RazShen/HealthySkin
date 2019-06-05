package main

import "net"
import "fmt"
import "bufio"
import "os"

func main() {

	// connect to this socket
	conn, err := net.Dial("tcp", "127.0.0.1:8181")
	if err != nil {
		fmt.Println("Error connecting to server:", err.Error())
		os.Exit(1)
	}
	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)

		usrName, _ := reader.ReadString('\n')
		if usrName == "q\n" {
			os.Exit(0)
		}
		// send to socket
		fmt.Fprintf(conn, strlen(usrName)+"\n")
		fmt.Fprintf(conn, usrName+"\n")
		password, _ := reader.ReadString('\n')
		if password == "q\n" {
			os.Exit(0)
		}
		// send to socket
		fmt.Fprintf(conn, strlen(password)+"\n")
		fmt.Fprintf(conn, password+"\n")
		fmt.Print("Please insert path of image to send to HealthySkin server: " +
			"\n(to quit insert 'q')\n")
		password, _ := reader.ReadString('\n')
		if password == "q\n" {
			os.Exit(0)
		}
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
