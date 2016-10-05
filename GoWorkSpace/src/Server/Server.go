package main

import (
    "io/ioutil"
	"fmt"
	"net"
	"os"
	"encoding/json"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8081"
	CONN_TYPE = "tcp"
	LOADSAVE  = "LoadSave"
)

type Block struct {
	ID int `json:"id"`
	X int `json:"x"`
	Y int `json:"y"`
	Width int `json:"width"`
	Height int `json:"height"`
} 

func jsonStuff() []byte {
	fmt.Println("getting file")
	dat, err := ioutil.ReadFile("C:/Independent-Study-Server/GoWorkSpace/src/Server/test.txt")
	
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println("Got file")
	}
	
	bytes := []byte(string(dat))
	var blocks[][] Block
	json.Unmarshal(bytes, &blocks)
	
	
	fmt.Println("getting json")
	b, _ := json.Marshal(blocks)
	fmt.Println("returning json")
	return b
}

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	text := make([]byte, 0, 4096)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	
	text = append(text, buf[:reqLen]...)
	fmt.Print("Message received: " + string(text))
	
	if strings.Contains(string(text), "json") {
		conn.Write(jsonStuff())
		fmt.Println("Sent json to client")
	} else {
		// Send a response back to person contacting us.
		conn.Write([]byte(string(text)))
	}
	
	conn.Close()
}