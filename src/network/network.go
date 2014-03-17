package network
/*
Holds the functions used for communication on the network. 
Interface:
func BroadcastOnNet(msgOutChan <-chan types.OrderMsg) - Broadcasts a message on the network
func ListenOnNetwork(msgChan chan<- types.OrderMsg, networkAlive chan<- bool) 
- Listens for messages on the network and send them over the msgChan. 
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"../types"
	"time"
)

var sock *net.UDPConn // Socket used by both the listen- and send functions in the module.
		      // As it's a generic stream-oriented network connection, we wont have
		      // any collision problems (Golang for the win!). 

const NumbOfBroadcasts	= 5 //  To increase the chance that the other computers gets the message. No pun intended.
const BroadCastIp = "129.241.187.255" // Local netowrk IP
const NetworkPort = ":2224" // Port used
const MaxNonJson  = 10	// Holds the given amount of non JSON messages reveiced in a given intervall
const NonJsonInt  = 60	// Number of seconds in the intervall where we check for invalid JSON messages

//Broadcast message on the local network at the given port
func BroadcastOnNet(msgOutChan <-chan types.OrderMsg) {
	timeCheck := time.Second*15 // We will only try to send a new message after 30 seconds if we had an error sending 
	timeSinceErr := time.Now()
	timeSinceErr.Add(timeCheck)
	for{
		if time.Since(timeSinceErr) >= timeCheck{
			select {
			case msg := <-msgOutChan:
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
						fmt.Printf("Wont send messages on network for %d time. \n", timeCheck)
						timeSinceErr = time.Now()
						break
					}
				}
			}
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
// Listen for messages over the network
func ListenOnNetwork(msgInChan chan<- types.OrderMsg, networkAlive chan<- bool) {
	addr, err := net.ResolveUDPAddr("udp", NetworkPort)
	if err != nil {			     // If we have an error here we can't listen on the network, so we tell the order.
		log.Printf("Error: %v", err) // module that we are shutting down before we do
		networkAlive <- false
		return
	}
	sock, err = net.ListenUDP("udp", addr)
	if err != nil {  		     // If we have an error here we can't listen on the network, so we tell the order
		log.Printf("Error: %v", err) //  module that we are shutting down before we do.
		networkAlive <- false
		return
	}
	sAddr, err := net.ResolveUDPAddr("udp", getSelfIP()+NetworkPort) // Get the computer's address on the network so it
	if err != nil {							 // doesn't read its own broadcasts.If we have an error
		log.Printf("Error: %v", err)                             // here we can't listen on the network, so we tell the
		networkAlive <- false					 // order module that we are shutting down before we do.		
		return
	}
	fmt.Println("Listnening on port", addr)
	var msg types.OrderMsg		 // Message variable to hold the received messages
	nonJson := 0 			 // Keeps tracks of the number of received messages that wasn't JSON objects
	intTime := time.Now()  		 // Start of time intervall for maximum amount of non-JSON object received
	buf := make([]byte, 1024)
	for {
		rlen, addr, err := sock.ReadFromUDP(buf)
		if err != nil{
			log.Printf("Error: %v", err)
		} else if addr != sAddr { 	// Don't handle if it's from the computer
			err = json.Unmarshal(buf[0:rlen], &msg)
			if err != nil {
				nonJson++
				log.Printf("Error: %v.", err)
				if time.Since(intTime) > NonJsonInt { // Reset the non-JSON counter if the intervall has run out
					nonJson =0
				}
				intTime = time.Now()	      // Restart time if there is a wrong message
				if nonJson > MaxNonJson{      // If we get to many messages that aren't JSON object,
					networkAlive <- false // we shut down the network mod
					return
				}
			} else if msg.Action != types.InvalidMsg{  // If the message received is not of type OrderMsg,
				msgInChan <- msg                   // all elements of msg will be zero(msg={0,0,0}), 
			}                                          // so we can check if it is valid or not
		}
	}
}
