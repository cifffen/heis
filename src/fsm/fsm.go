package fsm

import (
	"../drivers"
	"../orders"
	"fmt"
	"time"
)

const brakeDur = 10   //Duration, in milliseconds, of the braking time when stopping at a floor
const doorOpenDur = 3 //Duration, in seconds, of the time the door stays open when arriving at a floor
const Speed = 300     //The speed of the motor

type (
	Event int
	State int
)

const (
	OrderReached Event = iota
	TimerFinished
	NewOrder
	SwitchDirection
)
const (
	Idle State = iota
	Running
	AtFloor
)

var state State
var doorTimer <-chan time.Time
var brakeTimer <-chan time.Time
var direction orders.Direction
var noOrders bool

func InitElev() int {
	if drivers.ElevInit() == 0 { //IO init failed
		return 0
	} else {
		direction = orders.Down
		if drivers.ElevGetFloorSensorSignal() != -1 { //Check if the elevator is at a floor
		} else { //else, run downwards until one is found
				drivers.ElevSetSpeed(int(direction) * Speed)
				floor := drivers.ElevGetFloorSensorSignal()
				for floor == -1 {
					floor = drivers.ElevGetFloorSensorSignal()
				}
				drivers.ElevSetSpeed(int(-1*direction) * Speed)
				brake()
		}
		state = Idle
		noOrders = true
		fmt.Printf("Initialized\n")
		return 1
	}
}

//Reverse the direction to brake
func brake() {
	brakeTimer = time.After(time.Millisecond * brakeDur)
}

// Checks for events and runs the state machine when some occur
func EventManager() {
	orderReachedEvent := make(chan bool)
	newOrderEvent 	  := make(chan bool)
	switchDirEvent 	  := make(chan orders.Direction)
	noOrdersEvent	  := make(chan bool)
	go orders.OrderHandler(orderReachedEvent, newOrderEvent, switchDirEvent,  noOrdersEvent)
	for {
		select {
		case <-brakeTimer:    // Brake finished. Set speed to 0
		   fmt.Printf("Brake event\n")
			drivers.ElevSetSpeed(int(orders.Stop))
		case <-newOrderEvent:	// We got a new order, so noOrders must be set to false
			fmt.Printf("New order event\n")
			noOrders = false
			stateMachine(NewOrder)
		case direction = <-switchDirEvent:  // A direction change must happen, so direction is changed for the next time we set elevSetSpeed()
			stateMachine(SwitchDirection)
			fmt.Printf("Switch direction event %d\n", direction)
		case <-orderReachedEvent:			// Reached a floor where there is an order
			fmt.Printf("Order reached event\n")
			stateMachine(OrderReached)
		case <-doorTimer:					// Door timer is finished and we can close the doors
			fmt.Printf("Door timer finished\n")
			stateMachine(TimerFinished)
		case noOrders = <-noOrdersEvent:   // We now have no orders left. No orders i therefore set to true so we can go to Idle
			fmt.Printf("Door timer finished\n")
		}
	}
}

func stateMachine(event Event) {
	switch state {
	case Idle:
		switch event {
		case NewOrder:
				drivers.ElevSetSpeed(int(direction)*Speed)
				state = Running
				fmt.Printf("Running  \n")
		}
	case Running:
		switch event {
		case SwitchDirection:
		   drivers.ElevSetSpeed(int(direction)*Speed)
		case OrderReached:
			drivers.ElevSetSpeed(-1*int(direction)*Speed)
			brake()
			doorTimer = time.After(time.Second*doorOpenDur)
			drivers.ElevSetDoorOpenLamp(1)
			state = AtFloor
			fmt.Printf("Atfloor \n")
		}
	case AtFloor:
		switch event {
		case TimerFinished:
			drivers.ElevSetDoorOpenLamp(0)
			if noOrders {
				state = Idle
				fmt.Printf("Idle \n")
			} else {
				state = Running
				fmt.Printf("Runing \n")
				drivers.ElevSetSpeed(int(direction)*Speed)
			}
		}
	}
}
