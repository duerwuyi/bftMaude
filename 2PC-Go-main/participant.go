package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

var netConfig Network
var n int
var randomNumber = rand.New(rand.NewSource(0))
var order int

func main() {
	// Listen for incoming connections.
	//Connect to the coordinator
	conn, err := net.Dial("tcp", coordinatorNetwork.Host+":"+coordinatorNetwork.Port)
	if err != nil {
		fmt.Println("Error connection:", err.Error())
		os.Exit(1)
	}

	n, err = strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("excepted an index of participant network, actual: ", err.Error())
		os.Exit(1)
	}

	_, err = os.Create("log_participant_" + strconv.Itoa(n) + ".txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	netConfig = participantNetworks[n]

	//tcp server for ctp
	listener := serve(netConfig.Host, netConfig.Port)
	defer listener.Close()

	connToOther := make([]net.Conn, 0)
	for i, participant := range participantNetworks {
		if i == n {
			continue
		}
		c, err := net.Dial("tcp", participant.Host+":"+participant.Port)
		if err != nil {
			fmt.Println("Error connection:", err.Error())
			os.Exit(1)
		}
		connToOther = append(connToOther, c)
	}

	connFromOther := make([]net.Conn, 0)
	for i, _ := range participantNetworks {
		if i == n {
			continue
		}
		c, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		connFromOther = append(connFromOther, c)
		go func(c net.Conn) {
			//listen for ctp
			c.SetReadDeadline(time.Time{})
			message := receive(c)
			if message.Action == "CTP" {
				send(c, allMessage[message.Order])
			}
		}(c)

	}

	s := rand.NewSource(time.Now().Unix())
	randomNumber = rand.New(s)

	message := Message{}
	for message.Action != "DONE" {
		message = receive(conn)
		if message.Action == "DONE" {
			break
		}
		if message.Action != "PREPARE" {
			fmt.Println(strconv.Itoa(message.Order), ", ", message.Action, " not a prepare message")
			break
		}
		order = message.Order
		handlePrepare(conn, message)
		//phase 2
		message = receiveTW(conn) //timeout caused by msg loss
		if message.Action == "" {
			//ctp
			for _, c := range connToOther {
				send(c, Message{Action: "CTP", Order: order})
				m := receive(c)
				if m.Action != "" { //not unknown
					message = m
					break
				}
			}
		}
		if message.Action == "" {
			log.Default().Println("msg ", strconv.Itoa(message.Order), " remains unknown, ctp failed")
		}
		handleRequest(conn, message)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, message Message) {
	switch message.Action {
	case "ABORT":
		handleAbort(conn, message)
	case "COMMIT":
		handleCommit(conn, message)
	case "DONE":
		os.Exit(0)
	}
	sendACK(conn, message)
}

func getNameFile() string {
	return "log_participant_" + strconv.Itoa(n) + ".txt"
}

func handleAbort(conn net.Conn, message Message) {
	//Write the abort in log file
	writeToFile(getNameFile(), "p"+strconv.Itoa(n)+" abort, order "+strconv.Itoa(message.Order))
	allMessage[message.Order] = Message{"ABORT", message.Order}
}

func handleCommit(conn net.Conn, message Message) {
	writeToFile(getNameFile(), "p"+strconv.Itoa(n)+" commit, order "+strconv.Itoa(message.Order))
	allMessage[message.Order] = Message{"COMMIT", message.Order}
}

func sendACK(conn net.Conn, message Message) {
	send(conn, Message{Action: "ACK", Order: message.Order})
}

func handlePrepare(conn net.Conn, message Message) {
	writeToFile(getNameFile(), "ready")
	//random number between 0 and 1
	a := randomNumber.Float64()
	action := "COMMIT"
	if a <= ABORT_RATE {
		action = "ABORT"
		writeToFile(getNameFile(), "p"+strconv.Itoa(n)+" prepare to abort, order "+strconv.Itoa(message.Order))
	} else {
		writeToFile(getNameFile(), "p"+strconv.Itoa(n)+" prepare to commit, order "+strconv.Itoa(message.Order))
	}
	send(conn, Message{Action: action, Order: message.Order})
}
