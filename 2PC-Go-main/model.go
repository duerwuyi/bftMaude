package main

import "time"

type Message struct {
	Action string `json:"action"`
	Order  int    `json:"order"`
}

const NB_PART = 2
const REQUEST_NUM = 10000
const ABORT_RATE = 0.1
const LOST_RATE = 0.1
const TIMEOUT = time.Millisecond

type Network struct {
	Host    string
	Port    string
	NetType string
}

var coordinatorNetwork Network = Network{Host: "10.10.1.1", Port: "3334", NetType: "tcp"}

var participantNetworks = []Network{
	{Host: "10.10.1.2", Port: "3334", NetType: "tcp"},
	{Host: "10.10.1.3", Port: "3334", NetType: "tcp"},
}

/*
var coordinatorNetwork Network = Network{Host: "localhost", Port: "3334", NetType: "tcp"}
var participantNetworks = []Network{
	{Host: "localhost", Port: "3334", NetType: "tcp"},
	{Host: "localhost", Port: "3334", NetType: "tcp"},
}
*/
