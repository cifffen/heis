package main

import(
	"net"
	"fmt"
	"encoding/json"
	"log"
	"strings"
)

var sock *net.UDPConn
type ActionType int
const BroadCastIp = "192.168.1.255"
const NetworkPort = ":2224"

const(
	NewOrder ActionType = iota
	DeleteOrder
	Tender
)

type ButtonMessage struct {
	Action ActionType
	Floor int
	Button int
	TenderVal int
}

func BroadcastOnNet(msg ButtonMessage)(){
	addr, _ := net.ResolveUDPAddr("udp", NetworkIp+NetworkPort)
	buf,err1 := json.Marshal(msg)
	if err1 != nil{
		fmt.Println(err1)
	}
	//rAddr, _ := net.ResolveUDPAddr("udp", "192.168.1.255:2224")
	_,err1 =sock.WriteTo(buf,addr)
	if err1 != nil{
		log.Println(err1)
	}
}

// Return IP of own computer
func getSelfIP() string {
	conn,err :=net.Dial("udp", "google.com:80") 
	if err !=nil {
		log.Printf("Error: %v. Runing without self address checking.", err)
		return "localhost"
	} else{
		return strings.Split(string(conn.LocalAddr().String()), ":" )[0] 
	} 
}
func listenOnNetwork()(){
	addr, err := net.ResolveUDPAddr("udp", BroadCastPort)
	if err != nil {
		log.Printf("Error: %v. Runing without network connetion", err)
		return
	}
	sock, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("Error: %v. Runing without network connetion", err)
		return
	}
	fmt.Println("Listnening on port", addr)
	go func(){
		for {
			buf := make([]byte, 1024)
			rlen, addr, err := sock.ReadFromUDP(buf)
			sAddr, _ := net.ResolveUDPAddr("udp", "localhost:2224")
			if addr != sAddr{
				err = json.Unmarshal(buf[0:rlen], &m)
				if err != nil {
				  log.Printf("Error: %v. Sending aborted", err)
				}
			}
		}
	}
}