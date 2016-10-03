package main

import (
    "io/ioutil"
	"fmt"
	"net"
	"os"
	"encoding/json"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8081"
	CONN_TYPE = "tcp"
	LOADSAVE  = "LoadSave"
)

type Field struct {
	Blocks [][]struct {
		Height int `json:"height"`
		ID     int `json:"id"`
		Width  int `json:"width"`
		X      int `json:"x"`
		Y      int `json:"y"`
	} `json:"blocks"`
	Border struct {
		Height int `json:"height"`
		ID     int `json:"id"`
		Width  int `json:"width"`
		X      int `json:"x"`
		Y      int `json:"y"`
	} `json:"border"`
}

func jsonStuff() {
	dat, _ := ioutil.ReadFile("test.txt")
	bytes := []byte(string(dat))
	var field Field
	json.Unmarshal(bytes, &field)
	
	b, _ := json.Marshal(field)
	s := string(b)
	fmt.Println(s)
}


func main() {
	jsonStuff()
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
	
	// Send a response back to person contacting us.
	conn.Write([]byte(string(text)))
	
	conn.Close()
}