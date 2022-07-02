package main

import (
	chord "dht/chord"
	"os"
)

func main() {
	ip := os.Args[1]
	port := os.Args[2]
	ipSucc := os.Args[3]
	portSucc := os.Args[4]
	switch os.Args[5] {
	case "1":
		chord.JoinTest(ip, port, ipSucc, portSucc)
	case "2":
		chord.SleepTest(ip, port, ipSucc, portSucc)
	case "3":
		chord.StabilizeTest(ip, port, ipSucc, portSucc)
	case "4":
		chord.FireTest(ip, port, ipSucc, portSucc)
	default:
		chord.SelfTest(ip, port, ipSucc, portSucc)
	}
}
