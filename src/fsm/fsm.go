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
	Event int         // The event type
	State func(Event) // A state is a function that takes in an event and acts based on that event
)
type FSM struct { // The state machine type
	state     State // State holds the current state
	direction int   // Holds the direction of travel
	noOrders  bool  // True if there are no orders in the order list in the orders module
}

const ( // Events
	OrderReached Event = iota
	TimerFinished
	NewOrder
	SwitchDirection
)

//-----State diagram-----------------------------//
// States: Idle, Running, AtFloor
// Idle 	if NewOrder -> Running
// Running	if AtOrder	-> 	AtFloor
// AtFloor	if DoorTimer ? !noOrders -> Running
// AtFloor	if DoorTimer ? noOrders  -> Idle
//-----------------------------------------------//

var doorTimer <-chan time.Time  // The timer for the door (door open duration)
var brakeTimer <-chan time.Time // The timer for the brake (brake duration)

func InitElev() int {
	if drivers.ElevInit() == 0 { //IO init failed
		return 0 // Return 0 for failure
	} else {
		if drivers.ElevGetFloorSensorSignal() != -1 { //Check if the elevator is at a floor
		} else { //else, run downwards until one is found
			drivers.ElevSetSpeed(int(order.Down) * Speed)
			floor := drivers.ElevGetFloorSensorSignal()
			for floor == -1 {
				floor = drivers.ElevGetFloorSensorSignal()
			}
			drivers.ElevSetSpeed(int(order.Up) * Speed)
			brake()
		}
		fmt.Printf("Initialized\n")
		return 1 // Return 1 for success
	}
}

//Set the brake timer
func brake() {
	brakeTimer = time.After(time.Millisecond * brakeDur)
}

// Checks for events and runs the state machine when events occur
func EventManager() {
	var fsm FSM                      // Make a state machine
	fsm.state = fsm.idleState        // Set initial state to idle
	fsm.noOrders = true              // We have no orders at the start
	fsm.direction = int(orders.Down) // Set inital direction down (as our init runs downwards)

	// Event channels
	orderReachedEvent := make(chan bool)
	newOrderEvent := make(chan bool)
	switchDirEvent := make(chan orders.Direction)
	noOrdersEvent := make(chan bool)
	// Start Orderhandler in the orders module as a gorutine with the event channels as arguments to communicate
	go orders.OrderHandler(orderReachedEvent, newOrderEvent, switchDirEvent, noOrdersEvent)
	for {
		select {
		case <-brakeTimer: // Brake finished. Set speed to 0
			drivers.ElevSetSpeed(int(orders.Stop))
		case <-newOrderEvent: // New order so noOrders must be set to false
			fsm.noOrders = false
			fsm.state(NewOrder)
		case dir := <-switchDirEvent: // A direction change must happen, so direction is changed for the next time we set elevSetSpeed()
			fsm.direction = int(dir) // Converting from type orders.Direction to int to simplify
			fsm.state(SwitchDirection)
		case <-orderReachedEvent: // Reached a floor where there is an order
			fsm.state(OrderReached)
		case <-doorTimer: // Door timer is finished and we can close the doors
			fsm.state(TimerFinished)
		case fsm.noOrders = <-noOrdersEvent: // We now have no orders left. No orders i therefore set to true so we can go to Idle
		}
	}
}

// Idle state
func (fsm *FSM) idleState(event Event) {
	switch event {
	case NewOrder: // If there is a new order, set speed and go to running state
		drivers.ElevSetSpeed(fsm.direction * Speed)
		fsm.state = fsm.runningState
	}
}

// Running state
func (fsm *FSM) runningState(event Event) {
	switch event {
	case SwitchDirection: // If there is a change in direction we set the direction again
		drivers.ElevSetSpeed(fsm.direction * Speed)
	case OrderReached: // When we reach an order, we switch the direction to brake, start the door timer and lights and go to at floor state
		drivers.ElevSetSpeed(-1 * fsm.direction * Speed)
		brake()
		doorTimer = time.After(time.Second * doorOpenDur)
		drivers.ElevSetDoorOpenLamp(1)
		fsm.state = fsm.atFloorState
	}
}

// At floor state
func (fsm *FSM) atFloorState(event Event) {
	switch event {
	case TimerFinished: // When the door timer is finished we turn of the door lights
		drivers.ElevSetDoorOpenLamp(0)
		if fsm.noOrders { // If there aren't any more order we go to the idle state
			fsm.state = fsm.idleState
		} else { // If there are more orders we set speed and go to the running state
			drivers.ElevSetSpeed(fsm.direction * Speed)
			fsm.state = fsm.runningState
		}
	}
}
