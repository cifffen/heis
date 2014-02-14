package main

import (
	//"../drivers"
	"../fsm"
	"../orderMod"
	"fmt"
)



func main(){
	if fsm.InitElev()==0{
		fmt.Printf("Unable to initialize elevator hardware.\n")
	}
	orderReachedEvent := make(chan bool)
	newOrderEvent:= make (chan bool)
	go orderMod.AtOrder(orderReachedEvent)
	go orderMod.GetOrders(newOrderEvent)
	for {
		select {
			case <- newOrderEvent:
				fmt.Printf("New order event\n")
				fsm.StateMachine(fsm.NewOrder)
			case <-orderReachedEvent:
				fmt.Printf("Order reached event\n")
				fsm.StateMachine(fsm.OrderReached)
			case <- fsm.DoorTimer:
				fmt.Printf("Door timer finished\n")
				fsm.StateMachine(fsm.TimerFinished)
			
		}
	}
	
}
