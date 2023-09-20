package main

import "time"

type Message struct {
	Action string `json:"action"`
	Order  int    `json:"order"`
}

const NB_PART = 2
const REQUEST_NUM = 100
const ABORT_RATE = 0.1
const LOST_RATE = 0.01
const TIMEOUT = time.Second * 2

type Network struct {
	Host    string
	Port    string
	NetType string
}

var coordinatorNetwork Network = Network{Host: "localhost", Port: "3334", NetType: "tcp"}

var participantNetworks = []Network{
	{Host: "localhost", Port: "3335", NetType: "tcp"},
	{Host: "localhost", Port: "3336", NetType: "tcp"},
}
