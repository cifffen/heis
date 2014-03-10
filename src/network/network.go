package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"../types"
)

var sock *net.UDPConn // Socket used by both the listen- and send functions in the module. As it's a generic stream-oriented network connection, we wont have
					  // any collision problems (Golang for the win!).
const NumbOfBroadcasts 	= 5  				// Number of broadcasts per messages. To increase the chance that the other computers gets the message. No pun intended.
const BroadCastIp 		= "129.241.187.255" // Local netowrk IP
const NetworkPort 		= ":2224"			// Port used

//Broadcast message on the local network at the given port
func BroadcastOnNet(msg types.OrderMsg) {
	addr, err := net.ResolveUDPAddr("udp", BroadCastIp+NetworkPort)
	if err != nil {
		log.Printf("Error: %v",err)
	}
	buf, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error: %v",err)
	}
	for i := 0; i < NumbOfBroadcasts; i++ {  
		_, err = sock.WriteTo(buf, addr)
		if err != nil {
			log.Printf("Error: %v",err)
		}
	}
}

// Return self-IP 
func getSelfIP() string {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		log.Printf("Error: %v. Runing without self-address checking.", err)
		return "localhost"
	} else {
		return strings.Split(string(conn.LocalAddr().String()), ":")[0]
	}
}
func ListenOnNetwork(msgChan chan<- types.OrderMsg, networkAlive chan<- bool) {
	addr, err := net.ResolveUDPAddr("udp", NetworkPort)
	if err != nil {
		log.Printf("Error: %v", err)
		networkAlive <- false
		return
	}
	sock, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("Error: %v", err)
		networkAlive <- false
		return
	}
	sAddr, err := net.ResolveUDPAddr("udp", getSelfIP()+NetworkPort) // Get the computer's address on the network so it doesn't read its own broadcasts.
	if err != nil {
		log.Printf("Error: %v", err)
		networkAlive <- false
		return
	}
	fmt.Println("Listnening on port", addr)
	var msg types.OrderMsg
	tooManyNonJson := 0 
	buf := make([]byte, 1024)
	for {
		rlen, addr, err := sock.ReadFromUDP(buf)
		if err != nill{
			log.Printf("Error: %v", err)
		} else if addr != sAddr { // Don't handle if it's from the computer
			err = json.Unmarshal(buf[0:rlen], &msg)
			if err != nil {
				tooManyNonJson++
				log.Printf("Error: %v.", err)
				if tooManyNonJson > 100{
					networkAlive <- false
					return
				}
			} else if msg.Action != types.InvalidMsg{  // If the message received is not of type OrderMsg, all elements of msg will be zero(msg={0,0,0}), so we can check if it is valid or not
				msgChan <- msg
			}
		}
	}
}
