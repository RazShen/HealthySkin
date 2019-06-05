package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const BUFFERSIZE = 1024

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func connectToServer(conn net.Conn, reader *bufio.Reader) error {
	/*
		Sending user name and password to server
	*/
	fmt.Println("Enter your username")
	usrName, _ := reader.ReadString('\n')
	if usrName == "q\n" {
		os.Exit(0)
	}
	// send to socket
	_, err := fmt.Fprintf(conn, usrName+"\n")

	fmt.Println("Enter your password")
	password, err := reader.ReadString('\n')
	if password == "q\n" {
		os.Exit(0)
	}

	// send to socket
	_, err = fmt.Fprintf(conn, password+"\n")
	if err != nil {
		return err
	}
	/*
	   Finish sending user name and password to server
	*/
	message, _ := bufio.NewReader(conn).ReadString('\n')
	if strings.Contains(message, "0") {

	} else {

	}
	fmt.Print("Message from server: " + message)
	return nil
}

func sendImageToServer(conn net.Conn, reader *bufio.Reader) error {
	fmt.Print("Please insert path of image to send to HealthySkin server: " +
		"\n(to quit insert 'q')\n")
	imagePath, _ := reader.ReadString('\n')
	if imagePath == "q\n" {
		os.Exit(0)
	}
	trimmedFN := strings.Trim(imagePath, "\n")
	image, err := os.Open(trimmedFN)
	if err != nil {
		fmt.Println("Image Path doesn't exist.. exiting", err.Error())
		os.Exit(1)
	}

	fileInfo, err := image.Stat()
	if err != nil {
		return err
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fmt.Println("Sending filesize")
	_, err = conn.Write([]byte(fileSize))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file")
	for {
		_, err = image.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		_, err = conn.Write(sendBuffer)
	}
	err = image.Close()
	if err != nil {
		return err
	}
	return nil
}

func main() {

	// connect to this socket
	conn, err := net.Dial("tcp", "172.18.34.99:8181")
	if err != nil {
		fmt.Println("Error connecting to server:", err.Error())
		os.Exit(1)
	}
	// read in input from stdin
	reader := bufio.NewReader(os.Stdin)

	err = connectToServer(conn, reader)
	if err != nil {
		fmt.Println("Error sending user name or password to server")
	}

	err = sendImageToServer(conn, reader)
	if err != nil {
		fmt.Println("Error sending image to server")
	}

	// listen for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
}