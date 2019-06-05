package main

import (
	"io"
	"net"
	"strconv"
)
import "fmt"
import "strings" // only needed below for sample processing
import "os"

const (
	CONN_HOST  = "172.18.34.99"
	CONN_PORT  = "8181"
	CONN_TYPE  = "tcp"
	BUFFERSIZE = 1024
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
	defer ln.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	// accept connection on port
	for {
		fmt.Println("Waiting for new connection...")
		conn, err := ln.Accept()
		handleClient(conn)
		if err != nil {
			fmt.Println("Error connecting to user:", err.Error())
			os.Exit(1)
		}
	}

}

func handleClient(conn net.Conn) {
	//// run loop forever (or until ctrl-c)
	//for {
	//	// will listen for message to process ending in newline (\n)
	//	message, err := bufio.NewReader(conn).ReadString('\n')
	//	if err != nil {
	//		fmt.Println("Error reading from user:", err.Error())
	//		conn.Close()
	//		return
	//	}
	//
	//	// output message received
	//	fmt.Print("Message Received:", string(message))
	//	// sample process for string received
	//
	//	newmessage := strings.ToUpper(message)
	//	// send new string back to client
	//	_, err = conn.Write([]byte(newmessage + "\n"))
	//	if err != nil {
	//		fmt.Println("Error writing to user:", err.Error())
	//		conn.Close()
	//		return
	//	}
	//
	//}

	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	_, err := conn.Read(bufferFileSize)
	if err != nil {
		handleError(conn, "Error reading file size from user:", err)
		return
	}

	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	_, err = conn.Read(bufferFileName)
	if err != nil {
		handleError(conn, "Error reading file name from user:", err)
		return
	}

	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}

	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, conn, (fileSize - receivedBytes))
			_, err = conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			if err != nil {
				handleError(conn, "Error reading file bytes from user:", err)
				return
			}
			break
		}
		io.CopyN(newFile, conn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	fmt.Println("Received file completely!")
}

func handleError(conn net.Conn, errMsg string, err error) {
	fmt.Println(errMsg, err.Error())
	conn.Close()
}
