package main

import (
	"fmt"
	"../fsm"
	"../orders"
	"../types"
	"os"
    "../network"
	//"../pp"
	"time"
)

func main(){
	/*
	ppSuccess := pp.ProcessPairs(os.Args)    // Launch process pairs
	if ppSuccess==0{   // If the elevator program has crahsed too many times, process pairs will shut down and the program stops.
		fmt.Printf("Too many reboots. Elevator shutting down. \n") 
		go fsm.InitElev()         // Try to init the system so we can stop at a floor in case the elevator was runing during the last crash.
		time.Sleep(time.Second*4) // Sleep for 4 seconds so we can get the init done.
		os.Exit(1)	
	} 
	*/
	if fsm.InitElev()==0{ // If we fail to init the IO we exit the program. Process pairs will eventually start the program back up again
		fmt.Printf("Error: Unable to initialize elevator hardware. Shuting down.\n")
		time.Sleep(time.Second*4) // Sleep for a few seconds so we can read the error message.
		os.Exit(1) 
	}
	// Event channels
	orderReachedEvent 	:= make(chan bool)
	newOrderEvent 		:= make(chan bool)
	newDirEvent 		:= make(chan int)
	noOrdersEvent 		:= make(chan bool)
	doorOpen			:= make(chan bool)
	// Network  and message channels
	msgInChan 		:= make(chan types.OrderMsg) // Channel used to send messages from the network module
	msgOutChan 		:= make(chan types.OrderMsg) // Channel used to send messages to the network module
	netAliveChan	:= make(chan bool)			 // Channel used to tell if the network module has shut downs
	
	go orders.OrderHandler(orderReachedEvent, newOrderEvent, newDirEvent, noOrdersEvent, msgInChan, msgOutChan, netAliveChan, doorOpen)
	go network.ListenOnNetwork(msgInChan, netAliveChan)
    go network.BroadcastOnNet(msgOutChan)
	fsm.EventManager(orderReachedEvent, newOrderEvent, newDirEvent, noOrdersEvent, doorOpen) 
}
