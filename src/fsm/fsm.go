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
	AtEndFloor
)
const (
	Idle State = iota
	Running
	AtFloor
)

var state State
var DoorTimer <-chan time.Time
var BrakeTimer <-chan time.Time

func InitElev() int {
	if drivers.ElevInit() == 0 { //IO init failed
		return 0
	} else {
		if drivers.ElevGetFloorSensorSignal() != -1 { //Check if the elevator is at a floor
		} else { //else, run downwards until one is found
				drivers.ElevSetSpeed(int(orders.Down) * Speed)
				floor := drivers.ElevGetFloorSensorSignal()
				for floor == -1 {
					floor = drivers.ElevGetFloorSensorSignal()
				}
				drivers.ElevSetSpeed(int(orders.Up) * Speed)
				brake()
		}
		orders.InitOrderMod(floor)
		state = Idle
		fmt.Printf("Initialized\n")
		return 1
	}
}

//Reverse the direction to brake
func brake() {
	BrakeTimer = time.After(time.Millisecond * brakeDur)
}

// Checks for events and runs the state machine when some occur
func EventManager() {
	orderReachedEvent := make(chan bool)
	newOrderEvent := make(chan bool)
	atEndEvent := make(chan bool)
	go orders.CheckForEvents(orderReachedEvent, newOrderEvent, atEndEvent)
	for {
		select {
		case <-BrakeTimer:
			drivers.ElevSetSpeed(int(orders.Stop))
			fmt.Printf("Ferdig \n")
		case <-newOrderEvent:
			fmt.Printf("New order event\n")
			stateMachine(NewOrder)
		case <-atEndEvent:
			stateMachine(AtEndFloor)
		case atOrder := <-orderReachedEvent:
			if atOrder {
				fmt.Printf("Order reached event\n")
				stateMachine(OrderReached)
			}
		case <-DoorTimer:
			fmt.Printf("Door timer finished\n")
			stateMachine(TimerFinished)

		}
		/*
				if state== Running{
			        drivers.ElevSetSpeed(int(orders.ReturnDirection())*Speed)
			    }
		*/
	}
}

func stateMachine(event Event) {
	switch state {
	case Idle:
		switch event {
		case NewOrder:
			if orders.GetDir() != 0 {
				drivers.ElevSetSpeed(int(orders.GetDir()) * Speed)
				state = Running
			} else {
				DoorTimer = time.After(time.Second * doorOpenDur)
				drivers.ElevSetDoorOpenLamp(1)
				state = AtFloor
			}
		}
	case Running:
		switch event {
		case AtEndFloor:
			drivers.ElevSetSpeed(int(orders.ReturnDirection()) * Speed)
		case OrderReached:
			drivers.ElevSetSpeed(-1 * int(orders.ReturnDirection()) * Speed)
			brake()
			DoorTimer = time.After(time.Second * doorOpenDur)
			drivers.ElevSetDoorOpenLamp(1)
			state = AtFloor
			fmt.Printf("Atfloor \n")
		}
	case AtFloor:
		switch event {

		case OrderReached:
			fmt.Printf("At floor again\n")
			orders.GetDir()
			state = AtFloor
		case TimerFinished:
			if orders.IsLocOrdMatEmpty() {
				drivers.ElevSetDoorOpenLamp(0)
				state = Idle
				fmt.Printf("Idle \n")
			} else if orders.GetDir() == orders.Stop {
				DoorTimer = time.After(time.Second * doorOpenDur)
			} else {
				drivers.ElevSetDoorOpenLamp(0)
				state = Running
				fmt.Printf("Runing \n")
				drivers.ElevSetSpeed(int(orders.GetDir()) * Speed)
			}

		}
	}
}
