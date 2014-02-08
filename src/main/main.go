package main

import (
	//"../drivers"
	"../fsm"
	"../orderMod"
)



func main() int{
	if fsm.InitElev()==0{
		fmt.Printf("Unable to initialize elevator hardware.\n")
		return 1
	}
	orderReachedEvent := make(chan bool)
	newOrderEvent:= make (chan bool)
	go orderMod.AtOrder(orderReachedEvent)
	go orderMod.GetOrders(NewOrderEvent)
	for {
		select {
			case <-orderReachedEvent:
				stateMachine(fsm.OrderReached)
			case <- fsm.DoorTimer:
				statemachine(fsm.TimerFinished)
			case <- NewOrderEvent:
				statemachine(fsm.NewOrder)
		}
	}
	
	return 0
}
