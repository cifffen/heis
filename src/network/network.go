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

const NumbOfBroadcasts = 5
const BroadCastIp = "192.168.1.255"
const NetworkPort = ":2224"

const (
	InvalidMsg ActionType  = iota  //	Only used to check if the message recieved is of type ButtonMsg.
	NewOrder ActionType 		 //
	DeleteOrder
	Tender
	AddOrder
)
type OrderType struct{
	Button 	int			// Holds the button on the floor, Up or Down
	Floor 	int			// Holds the floor
}

type ButtonMsg struct {
	Action    	ActionType   	// Holds what the information of what to do with the message
	Order 		OrderType 		// Holds the floor and button of the order
	TenderVal 	int				// If the action is a Tender, this will hold the cost from the sender, that is, the value from the cost function for this order
}

func BroadcastOnNet(msg ButtonMsg) {
	addr, err := net.ResolveUDPAddr("udp", NetworkIp+NetworkPort)
	if err != nil {
		fmt.Println(err)
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < NumbOfBroadcasts; i++ {  
		_, err = sock.WriteTo(buf, addr)
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
	var msg ButtonMessage
	for {
		buf := make([]byte, 1024)
		rlen, addr, err := sock.ReadFromUDP(buf)
		if addr != sAddr { // Don't handle if it's from the computer
			err = json.Unmarshal(buf[0:rlen], &msg)
			if err != nil {
				log.Printf("Error: %v.", err)
			} else if msg.Action != InvalidMsg{  // If the message received is not of type ButtonMsg, all elements of msg will be zero(msg={0,0,0}), so we can check if it is valid or not
				msgChan <- m
			}
		}
	}

}
