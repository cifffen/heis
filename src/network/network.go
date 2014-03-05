package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

var sock *net.UDPConn

type ActionType int

const BroadCastIp = "192.168.1.255"
const NetworkPort = ":2224"

const (
	NewOrder ActionType = iota
	DeleteOrder
	Tender
)

type ButtonMsg struct {
	Action    ActionType
	Floor     int
	Button    int
	TenderVal int
}

func BroadcastOnNet(msg ButtonMsg) {
	addr, _ := net.ResolveUDPAddr("udp", NetworkIp+NetworkPort)
	buf, err1 := json.Marshal(msg)
	if err1 != nil {
		fmt.Println(err1)
	}
	//rAddr, _ := net.ResolveUDPAddr("udp", "192.168.1.255:2224")
	_, err1 = sock.WriteTo(buf, addr)
	if err1 != nil {
		log.Println(err1)
	}
}

// Return IP of own computer
func getSelfIP() string {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		log.Printf("Error: %v. Runing without self address checking.", err)
		return "localhost"
	} else {
		return strings.Split(string(conn.LocalAddr().String()), ":")[0]
	}
}
func ListenOnNetwork(msgChan chan<- ButtonMsg) {
	addr, err := net.ResolveUDPAddr("udp", BroadCastPort)
	if err != nil {
		log.Printf("Error: %v. Running without network connetion", err)
		return
	}
	sock, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("Error: %v. Running without network connetion", err)
		return
	}
	sAddr, err := net.ResolveUDPAddr("udp", getSelfIP()+NetworkPort) // Get the computer's address on the network so it doesn't read its own broadcasts.
	if err != nil {
		log.Printf("Error: %v. Sending aborted", err)
	}
	fmt.Println("Listnening on port", addr)

	for {
		buf := make([]byte, 1024)
		rlen, addr, err := sock.ReadFromUDP(buf)
		if addr != sAddr { // Don't handle if it's from the computer
			err = json.Unmarshal(buf[0:rlen], &m)
			if err != nil {
				log.Printf("Error: %v.", err)
			} else {
				msgChan <- m
			}
		}
	}

}
