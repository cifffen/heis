package main

import (
	"../drivers"
	"../fsm"
	"../orderHandler"
)
func main() int{
	if fsm.InitElev()==0{
		fmt.Printf("Unable to initialize elevator hardware.\n")
		return 1
	}
	orderReachedEvent := make(chan bool)
	timerOutEvent := make(chan bool)
	newOrderEvent:= make (chan bool)
	go AtOrder(orderReachedEvent)
	go GetOrder(NewOrderEvent)
	for{
		select{
			case <-orderReachedEvent:
				stateMachine(OrderReached)
			case <- timerOutEvent:
				statemachine(TimerFinished)
			case <- NewOrderEvent:
				statemachine(NewOrder)
		}
	}
	
	return 0
}
