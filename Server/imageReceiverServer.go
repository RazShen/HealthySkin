package main

import (
	"HealthySkin/DBDAL"
	"HealthySkin/MLAPI"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings" // only needed below for sample processing
)

const (
	CONN_HOST  = "127.0.0.1"
	CONN_PORT  = "8181"
	CONN_TYPE  = "tcp"
	BUFFERSIZE = 1024
)

// main function of the server- it listens to tcp connections and activates handle client for each
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

// this function handles new client connection to the server
func handleClient(conn net.Conn) {
	fmt.Println("New client connected to server!")
	//bufferFileName := make([]byte, 64)
	bufferUserName := make([]byte, 64)
	bufferPassword := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	// Get User name
	_, err := conn.Read(bufferUserName)
	if err != nil {
		handleError(conn, "Error reading user name from user:", err)
		return
	}

	userName := strings.Split(string(bufferUserName), "\n")[0]

	// Get Password
	_, err = conn.Read(bufferPassword)
	if err != nil {
		handleError(conn, "Error reading file name from user:", err)
		return
	}

	password := strings.Split(string(bufferPassword), "\n")[0]

	userInfo := DBDAL.GetUserDetailsById(userName, GetMD5Hash(password))
	if userInfo == nil {
		_, err = fmt.Fprintln(conn, "0")
		userInfo = getNewUserDetailsFromClinet(conn, userName, GetMD5Hash(password))
		DBDAL.SaveUserInfoDetails(*userInfo)
		if err != nil {
			handleError(conn, "Error writing to user:", err)
			return
		}

	} else {
		_, err = fmt.Fprintln(conn, "1")
		if err != nil {
			handleError(conn, "Error writing to user:", err)
			return
		}
	}

	fileName := "image from user_" + userName

	_, err = conn.Read(bufferFileSize)
	if err != nil {
		handleError(conn, "Error reading file size from user:", err)
		return
	}
	// get the file size of the image
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	newFile, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}

	defer newFile.Close()
	defer conn.Close()
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
	prediction := MLAPI.GetIsCancerImage(newFile, userInfo)
	if prediction == true {
		_, err = conn.Write([]byte("Prediction for image: negative.\n"))
	} else {
		_, err = conn.Write([]byte("Prediction for image: positive. Go to doctor urgently!\n"))
	}
}

// handle error
func handleError(conn net.Conn, errMsg string, err error) {
	fmt.Println(errMsg, err.Error())
	conn.Close()
}

// the server doesn't save the password itself but a md5 of the password
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// this function returns struct of user details from the client
func getNewUserDetailsFromClinet(conn net.Conn, userName string, md5Password string) *DBDAL.UserInfo {
	s, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		handleError(conn, "Error reading age from user:", err)
		return nil
	}
	age, _ := strconv.Atoi(s)
	s, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		handleError(conn, "Error reading weight from user:", err)
		return nil
	}
	weight, _ := strconv.Atoi(s)
	return &DBDAL.UserInfo{md5Password, userName, age, weight}
}
