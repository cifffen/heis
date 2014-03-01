package fsm

import (
	"../drivers"
	"../orderMod"
	"time"
	"fmt"
)
const brakeDur = 5 		//Duration, in milliseconds, of the braking time when stopping at a floor 
const doorOpenDur = 3  	//Duration, in seconds, of the time the door stays open when arriving at a floor
const Speed = 300      	//The speed of the motor
type(
	Event int
	State int
)

const(
	OrderReached Event =iota
    TimerFinished
	NewOrder
)
const (
	Idle State =iota
	Running
	AtFloor
)
var state State
var DoorTimer <-chan time.Time 

func InitElev() int{
	if drivers.ElevInit() ==0 {  //IO init failed
		return 0
	} else {
		drivers.ElevSetSpeed(int(orderMod.Down)*Speed)
		for drivers.ElevGetFloorSensorSignal() ==-1 {
		}
		orderMod.InitOrderMod(drivers.ElevGetFloorSensorSignal())
		drivers.ElevSetSpeed(int(orderMod.Stop))
		state = Idle
		fmt.Printf("Initialized\n")
		return 1
	}
}

//Reverse the direction to brake
func brake()(){
	drivers.ElevSetSpeend(-1*int(orderMod.GetDir())*Speed)  
	time.Sleep(time.Millisecond*brakeDur)
	drivers.ElevSetSpeend(int(orderMod.Stop))
}

// Checks for events and runs the state machine when some occur
func EventManager() (){
	orderReachedEvent := make(chan bool)
	newOrderEvent:= make (chan bool)
	go orderMod.checkForEvents(orderReachedEvent,newOrderEvent)
	for {
		select {
			case <- newOrderEvent:
				fmt.Printf("New order event\n")
				StateMachine(fsm.NewOrder)
			case <-orderReachedEvent:
				fmt.Printf("Order reached event\n")
				StateMachine(fsm.OrderReached)
			case <- DoorTimer:
				fmt.Printf("Door timer finished\n")
				StateMachine(fsm.TimerFinished)
			
		}
	}
}
	
	
	
func StateMachine(event Event)(){
	switch state{
		case Idle:
			switch event{
				case NewOrder:
					drivers.ElevSetSpeed(int(orderMod.GetDir())*Speed)
					state = Running
			}
		case Running:
			switch event{
				case OrderReached:
					go brake()
					DoorTimer = time.After(time.Second*doorOpenDur)
					state=AtFloor
			}
		case AtFloor:
			switch event{
				case TimerFinished:
					if orderMod.isOrderMatrixEmpty(){
						state=Idle
						orderMod.GetDir()  //Called jsut to set the the variable direction in orderMod to Stop
					} else {
						state=Running
						drivers.ElevSetSpeed(int(orderMod.GetDir())*Speed)
					}
			}
		default:
			fmt.Printf("Invalid state. \n")
	}
}
