package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var allMessage = make([]Message, REQUEST_NUM)

func send(conn net.Conn, message Message) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message:", err.Error())
		os.Exit(1)
	}
	_, err = conn.Write(jsonMessage)
	if err != nil {
		fmt.Println("Error sending message:", err.Error())
		return
	}
	log.Default().Println("Message sent.", message)
}

func receive(conn net.Conn) Message {
	// bufio.NewReader(conn).ReadString('\n')

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		os.Exit(1)
	}
	if n > 0 {
		log.Default().Println("Message received from", conn.RemoteAddr().String(), string(buf))
	}
	var message Message
	err = json.Unmarshal(buf[:n], &message)

	if err != nil {
		fmt.Println("Error reading JSON:", err.Error())
	}

	return message
}

func receiveTW(conn net.Conn) Message {
	// bufio.NewReader(conn).ReadString('\n')

	buf := make([]byte, 2048)
	conn.SetReadDeadline(time.Now().Add(TIMEOUT))
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		conn.SetReadDeadline(time.Time{})
		return Message{}
	}
	conn.SetReadDeadline(time.Time{})
	if n > 0 {
		log.Default().Println("Message received from", conn.RemoteAddr().String(), string(buf))
	}
	var message Message
	err = json.Unmarshal(buf[:n], &message)

	if err != nil {
		fmt.Println("Error reading JSON:", err.Error())
	}

	return message
}

func writeToFile(fileName, content string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(content + "\n")
	if err != nil {
		log.Fatal(err)
	}
}

func serve(host string, connType string) net.Listener {
	listener, err := net.Listen("tcp", host+":"+connType)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on " + host + ":" + connType)
	return listener
}
