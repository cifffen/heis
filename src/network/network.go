package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"../types"
)

var sock *net.UDPConn

type ActionType int

const NumbOfBroadcasts = 5
const BroadCastIp = "129.241.187.255"
const NetworkPort = ":2224"



func BroadcastOnNet(msg types.OrderMsg) {
	addr, err := net.ResolveUDPAddr("udp", BroadCastIp+NetworkPort)
	if err != nil {
		fmt.Println(err)
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < NumbOfBroadcasts; i++ {  
		_, err = sock.WriteTo(buf, addr)
		fmt.Printf("Printing \n")
		if err != nil {
			log.Println(err)
		}
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
func ListenOnNetwork(msgChan chan<- types.OrderMsg) {
	addr, err := net.ResolveUDPAddr("udp", NetworkPort)
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
	var msg types.OrderMsg
	for {
		buf := make([]byte, 1024)
		rlen, addr, err := sock.ReadFromUDP(buf)
		if addr != sAddr { // Don't handle if it's from the computer
			err = json.Unmarshal(buf[0:rlen], &msg)
			if err != nil {
				log.Printf("Error: %v.", err)
			} else if msg.Action != types.InvalidMsg{  // If the message received is not of type OrderMsg, all elements of msg will be zero(msg={0,0,0}), so we can check if it is valid or not
				fmt.Printf("Msg in network: %d \n", msg)
				msgChan <- msg
			}
		}
	}
}
