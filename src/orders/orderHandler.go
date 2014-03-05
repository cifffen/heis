package orders

import (
	"../drivers"
	"../network"
	"fmt"
	"time"
)

type Direction int

const (
	Down Direction = -1
	Up             = 1
	Stop           = 0
)

const SamplingTime = 1
const Floors = 4

///////////////////////////
/*   OrderMatrix Layout  */
//   Up // Down	// Panel//   <-- Button order
//------------------------
//    	//  	// 		//  	Lowest floor
//		//		//		//
//	...	//	..	//	...	//
//	...	//	..	//	...	//
//		//		//		//   	Highest floor (Floors)
//////////////////////////
var orderMatrix [Floors][3]int

var direction Direction // Keeps the last direction the elevator was heading
var prevFloor int       // Holds the previous floor the elevator past
var orderCount int      // Keeps track of the number of active orders.
var firstOrderFloor int // Keeps the floor where the first order came when the elevator was Idle
var atEndFloor bool     // True if the elevator is at the lowest or highest floor. Used to change direction in case it got "lost"

// Initializes the order module.
func InitOrderMod(floor int) {
	direction = Down
	prevFloor = floor
	orderCount = 0
	firstOrderFloor = -1
	atEndFloor = false
	go OrderHandler()
}

// Checks if the order matrix is empty
func IsOrderMatrixEmpty() bool {
	if orderCount == 0 {
		return true
	} else {
		return false
	}
}

//Returns current direction.
func ReturnDirection() Direction {
	if atEndFloor {
		return -1 * direction
	} else {
		return direction
	}
}

//Check for events in order module
func CheckForEvents(orderReachedEvent chan<- bool, newOrderEvent chan<- bool, atEndEvent chan<- bool) {
	for {
		select {
		case <-time.After(time.Millisecond * SamplingTime):
			orderReachedEvent <- AtOrder()
			newOrderEvent <- GetOrders()
			if atEndFloor {
				atEndEvent <- true
			}
		}
	}
}

//Handles orders both locally and over the network
func OrderHandler() {
	msgChan := make(chan network.ButtonMsg)
}

// Check for orders
func getOrders() bool {
	firstOrderEvent := false
	for i := 0; i < Floors; i++ {
		for j := 0; j < 3; j++ {
			if (i != 0 && i != Floors-1) || (i == 0 && j != 1) || (i == Floors-1 && j != 0) { // Statement that makes sure that we don't check the Down button at the groud floor and
				if drivers.ElevGetButtonSignal(j, i) == 1 && orderMatrix[i][j] != 1 { // the Up button at the top floor, as they don't exist.
					orderMatrix[i][j] = 1
					drivers.ElevSetButtonLamp(drivers.TagElevLampType(j), i, 1)
					if orderCount == 0 { //set  newOrderEvent if there is made an order to an empty orderMatrix
						firstOrderEvent = true
						firstOrderFloor = i //remember where to first order was made for. Might not be necessary with more elevators.
					}
					orderCount++ // count number of active orders.
				}
			}
		}
	}
	return firstOrderEvent
}

//Delete given orders at current floor
func deleteFloorOrders(floor int) {
	if orderMatrix[floor][2] == 1 { // Delete panel buttnon
		orderMatrix[floor][2] = 0
		drivers.ElevSetButtonLamp(drivers.TagElevLampType(2), floor, 0)
		orderCount--
	}
	switch direction {
	case Up:
		if orderMatrix[floor][0] == 1 {
			orderMatrix[floor][0] = 0
			drivers.ElevSetButtonLamp(drivers.TagElevLampType(0), floor, 0)
			orderCount--
		}
	case Down:
		if orderMatrix[floor][1] == 1 {
			drivers.ElevSetButtonLamp(drivers.TagElevLampType(1), floor, 0)
			orderMatrix[floor][1] = 0
			orderCount--
		}
	}
}

// Check if the elevator should stop at a floor it passes
func atOrder() bool {
	floor := drivers.ElevGetFloorSensorSignal()
	if floor != -1 {
		prevFloor = floor
		drivers.ElevSetFloorIndicator(floor) //Set floor indicator
		if floor == Floors-1 {               // If the elevator is at the top floor the direction is changed as it can't go further Upwards.
			direction = Down
			atEndFloor = true
		} else if floor == 0 { // If the elevator is at the bottom floor the direction is changed as it can't go further Downwards.
			direction = Up
			atEndFloor = true
		} else {
			atEndFloor = false
		}
		if orderMatrix[floor][2] == 1 || firstOrderFloor == floor { // Stop if an order from the inside panel has been made at the current floor.
			firstOrderFloor = -1
			DeleteFloorOrders(floor)
			return true
		} else if (direction == Up && orderMatrix[floor][0] == 1) || (direction == Down && orderMatrix[floor][1] == 1) { // Stop if an order from the direction button at the current floor has been made and the elevator is going in that direction.
			DeleteFloorOrders(floor)
			return true
		}
	}
	return false
}

func GetDir() Direction {
	if drivers.ElevGetFloorSensorSignal() == -1 {
		return direction
	}
	if orderCount == 0 { //If called and no orders exisits, just set direction to Stop and return Stop as FSM will go to Idle.
		return Stop
	}
	var ordersAtCur [3]bool //	Holds all orders on the current floor
	var ordersInDir [2]bool // [0] is true if there are orders further up, [1] is true if there is any up
	var currDir int         // Varable to hold the current direction to be used in orderInDir. 0 for up and 1 for down.

	for i := 0; i < Floors; i++ {
		for j := 0; j < 3; j++ {
			if orderMatrix[i][j] == 1 {
				if i == prevFloor { //check for orders at current floor
					ordersAtCur[j] = true
				} else if i > prevFloor { // check for orders upwards
					ordersInDir[0] = true
				} else if i < prevFloor { // check for orders downwards
					ordersInDir[1] = true
				}
			}
		}
	}
	switch direction {
	/*
		case Stop:
			if firstOrderFloor != -1 {
				if firstOrderFloor == prevFloor {
					return Stop
				} else if firstOrderFloor > prevFloor {
					direction = Up
				} else if firstOrderFloor < prevFloor {
					direction = Down
				}
				return direction
			}
	*/
	case Up:
		currDir = 0
	case Down:
		currDir = 1
	}
	if ordersAtCur[currDir] || ordersAtCur[2] { //Just stay put if there is an order at current floor from the panel or from outside in the same direction as travel
		DeleteFloorOrders(prevFloor)
		return Stop
	} else if ordersInDir[currDir] { //Return current direction if there is an order in that direction
		return direction
	} else if ordersAtCur[currDir+int(direction)] { //Just stay put if there is an order at current flor in opposite direction
		firstOrderFloor = prevFloor
		direction = -1 * direction
		DeleteFloorOrders(prevFloor)
		return Stop
	} else if ordersInDir[currDir+int(direction)] { //Go in opposit direction if there is an order there there
		direction = -1 * direction
		return direction
	}
	return direction //direction = Stop	// Stay put if the logic above fails (Yeah, right...)
}
