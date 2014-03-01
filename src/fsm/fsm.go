package fsm

import (
	"../drivers"
	"../orderMod"
	"time"
	"fmt"
)
const brakeDur = 5		//Duration, in milliseconds, of the braking time when stopping at a floor 
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
var BrakeTimer <-chan time.Time 

func InitElev() int{
	if drivers.ElevInit() ==0 {  //IO init failed
		return 0
	} else {
		drivers.ElevSetSpeed(int(orderMod.Down)*Speed)
		for drivers.ElevGetFloorSensorSignal() ==-1 {
		}
		orderMod.InitOrderMod(drivers.ElevGetFloorSensorSignal())
		drivers.ElevSetSpeed(int(orderMod.Up)*Speed) 
		brake()
		state = Idle
		fmt.Printf("Initialized\n")
		return 1
	}
}

//Reverse the direction to brake
func brake()(){
    fmt.Printf("Begynner \n")
    BrakeTimer = time.After(time.Millisecond*brakeDur)


}

// Checks for events and runs the state machine when some occur
func EventManager() (){
    //syncChan := make(chan bool)
	orderReachedEvent := make(chan bool)
	newOrderEvent:= make (chan bool)
	go orderMod.CheckForEvents(orderReachedEvent,newOrderEvent)
	for {
		select {
		    case <- BrakeTimer:
			    drivers.ElevSetSpeed(int(orderMod.Stop))
			    fmt.Printf("Ferdig \n")
			case newOrder:=<- newOrderEvent:
				if newOrder{
				    fmt.Printf("New order event\n")
				    StateMachine(NewOrder)
				}
			case atOrder:= <-orderReachedEvent:
				if atOrder{
				    fmt.Printf("Order reached event\n")
				    StateMachine(OrderReached)
				 }
			case <- DoorTimer:
				fmt.Printf("Door timer finished\n")
				StateMachine(TimerFinished)
			//case <- syncChan:
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
				    drivers.ElevSetSpeed(-1*int(orderMod.ReturnDirection())*Speed)
					brake()
					DoorTimer = time.After(time.Second*doorOpenDur)
					state=AtFloor
					fmt.Printf("Atfloor \n")
			}
		case AtFloor:
			switch event{
				case TimerFinished:
					if orderMod.IsOrderMatrixEmpty(){
						state=Idle
						fmt.Printf("Idle \n")
						orderMod.GetDir()  //Called jsut to set the the variable direction in orderMod to Stop
					} else {
						state=Running
						fmt.Printf("Runing \n")
						drivers.ElevSetSpeed(int(orderMod.GetDir())*Speed)
					}
			}
		default:
			fmt.Printf("Invalid state. \n")
	}
}
