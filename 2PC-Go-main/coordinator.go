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

var randomNum = rand.New(rand.NewSource(0))

func main() {
	// Listen for incoming connections.
	listener := serve(coordinatorNetwork.Host, coordinatorNetwork.Port)

	// Close the listener when the application closes.
	defer listener.Close()
	//create log
	_, err := os.Create("log_coordinator.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	_, err = os.Create("time_duration.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Listen for two participant
	participants := make([]net.Conn, 0)

	for i := 0; i < NB_PART; i++ {
		participant, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		participants = append(participants, participant)
	}

	s := rand.NewSource(time.Now().Unix())
	randomNum = rand.New(s)

	fmt.Println(participants)
	start := time.Now()
	//start a request source
	//order := xxx()
	for i := 0; i < REQUEST_NUM; i++ {
		order := i
		//2pc
		twoPC(participants, order)
	}
	t := time.Now().Sub(start)
	writeToFile("time_duration.txt", t.String())
	sendDone(participants)
}

func twoPC(participants []net.Conn, order int) {
	// send prepare message
	prepare_message := Message{Action: "PREPARE", Order: order}
	writeToFile("log_coordinator.txt", "begin_transaction "+strconv.Itoa(order))
	sendPrepare(participants, prepare_message)
	// handle prepare response
	isAParticipantAbort := receivePrepare(participants, order)
	responseToAll(participants, order, isAParticipantAbort) // it may be lost
	receiveACK(participants)
	writeToFile("log_coordinator.txt", "end_of_transaction")
}

func sendPrepare(participants []net.Conn, message Message) {
	for _, participant := range participants {
		send(participant, message)
	}
}

func sendDone(participants []net.Conn) {
	for _, participant := range participants {
		send(participant, Message{Action: "DONE"})
	}
}

func receivePrepare(participants []net.Conn, order int) bool {
	isAParticipantAbort := false

	for _, participant := range participants {
		prepareResponse := receive(participant)

		if prepareResponse.Action == "ABORT" {
			isAParticipantAbort = true
		}
	}

	return isAParticipantAbort
}

func responseToAll(participants []net.Conn, order int, isAParticipantAbort bool) {
	for _, participant := range participants {
		a := randomNum.Float64()
		if a <= LOST_RATE {
			log.Default().Println("Message loss.")
			continue
		}
		if isAParticipantAbort {
			abort(participant, order)
		} else {
			Commit(participant, order)
		}
	}
	if isAParticipantAbort {
		writeToFile("log_coordinator.txt", "abort")
	} else {
		writeToFile("log_coordinator.txt", "commit")
	}
}

func receiveACK(participants []net.Conn) {
	for _, participant := range participants {
		ACK := receive(participant)

		if ACK.Action != "ACK" {
			writeToFile("log_coordinator.txt", strconv.Itoa(ACK.Order)+" ACK error")
		}
	}

}
func abort(conn net.Conn, order int) {
	message := Message{Action: "ABORT", Order: order}
	send(conn, message)
}

func Commit(conn net.Conn, order int) {
	message := Message{Action: "COMMIT", Order: order}
	send(conn, message)
}
