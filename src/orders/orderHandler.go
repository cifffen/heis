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

const (
	UpButton 	int = iota
	DownButton
	PanelButton
)

type TenderType struct{
	time 	time.Time
	val 	int
}

const TenderWon			 = 500 // Milliseconds
const TakeLostTenderTime = 20  // Seconds
const SamplingTime		 = 1   // Milliseconds
const Floors 			 = 4   // Number of floors

///////////////////////////
//   locOrdMat Layout  //   
//   Up // Down	// Panel//   <-- Button order
//------------------------
//    	//  	// 		//  	Lowest floor
//		//		//		//
//	...	//	..	//	...	//
//	...	//	..	//	...	//
//		//		//		//   	Highest floor (Floors)
//////////////////////////
var locOrdMat [Floors][3]int 	// Holds the orders that the elevator has accepted and will carry out
//var globOrdMat [Floors][2]int 	// Holds the direction orders on the network, both lost and non-desided tenders
var activeTenders map[network.OrderType] 	TenderType
var lostTenders map[network.OrderType] 		time.Time


//----Package variables-----//
//All are initialized in initOrderMod()and can only be changed by certain functions.
var direction 		Direction 	// Keeps the last direction the elevator was heading. Can only be changed in atOrder() and GetDir()
var prevFloor 		int      	// Holds the previous floor the elevator past. Can only be changed at atOrder()
var orderCount 		int      	// Keeps track of the number of active orders.
var firstOrderFloor int 		// Keeps the floor where the first order came from when the elevator was Idle
var atEndFloor 		bool     	// True if the elevator is at the lowest or highest floor. Used to change direction in case it got "lost"
var newOrder		bool   		// Set high to launch NewOrderEvent if an order is made and the orderMatrix is empyty
//--------------------------//

// Initializes the order module.
func InitOrderMod(floor int) {
	direction 		= Up
	prevFloor 		= floor
	orderCount 		= 0
	firstOrderFloor = -1
	atEndFloor 		= false
	newOrder 		= false
	activeTenders 	= make(map[network.OrderType] TenderType)
	lostTenders 	= make(map[network.OrderType] time.Time)
}
//Check for events in order module
func CheckForEvents(orderReachedEvent chan<- bool, newOrderEvent chan<- bool, atEndEvent chan<- bool) {
	msgChan := make(chan network.ButtonMsg)
	network.ListenOnNetwork(msgChan)
	for {
		select {
		case <-time.After(time.Millisecond * SamplingTime):
			orderReachedEvent <- atOrder()
			getOrders()
			checkTenderMaps()
			if newOrder {
				newOrderEvent <-true
				newOrder = false
			}
			if atEndFloor {
				atEndEvent <- true
			}
			
		case msg:= <-msgChan:
			orderHandler(msg)
		}
	}
}

func checkTenderMaps()(){
	for order, tenderTime := range lostTenders {
		if time.Since(tenderTime) > time.Second*TenderWon{
				var msg network.ButtonMsg
				msg.Action = network.AddOrder
				msg.Order = order
		}
	}
	for order, value := range activeTenders {
		if time.Since(value.time) > time.Millisecond*TakeLostTenderTime{
				var msg network.ButtonMsg
				msg.Action = network.AddOrder
				msg.Order = order
		}
	}
}

// Check for orders
func getOrders() {
	//firstOrderEvent := false
	var msg network.ButtonMsg
	msg.Action = network.NewOrder
	for i := range locOrdMat {
		for j := range locOrdMat[i] {
			if (i != 0 && i != Floors-1) || (i == 0 && j != 1) || (i == Floors-1 && j != 0) { // Statement that makes sure that we don't check the Down button at the groud floor and the Up button at the top floor, as they don't exist.
				if drivers.ElevGetButtonSignal(j, i) == 1 && locOrdMat[i][j] == 0 {
					msg.Order=network.OrderType{j, i}
					orderHandler(msg)
					//fmt.Printf("getorder \n")
					
				}
			}
		}
	}
	//return firstOrderEvent
}

// Check if the elevator should stop at a floor it passes
func atOrder() (orderReached bool) {
	floor := drivers.ElevGetFloorSensorSignal()
	orderReached = false
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
		dir:= ReturnDirection()
		var msg network.ButtonMsg
		msg.Action = network.DeleteOrder
		if locOrdMat[floor][PanelButton] == 1 || firstOrderFloor == floor { // Stop if an order from the inside panel has been made at the current floor.
			firstOrderFloor = -1
			msg.Order=network.OrderType{PanelButton, floor}
			orderHandler(msg)
			orderReached = true
		} 
		if (dir == Up && locOrdMat[floor][UpButton] == 1) { // Stop if an order from the direction button at the current floor has been made and the elevator is going in that direction.
			msg.Order=network.OrderType{UpButton, floor}
			orderHandler(msg)
			orderReached = true
		} else if (dir == Down && locOrdMat[floor][DownButton] == 1) {
			msg.Order=network.OrderType{DownButton, floor}
			orderHandler(msg)
			orderReached = true
		}	
	}
	return 
}

//Handles orders both locally and over the network
func orderHandler(msg network.ButtonMsg) {
	if checkMsg(msg) {
		fmt.Printf("msg %d \n", msg)
		order 		  := msg.Order
		floor, button := order.Floor, order.Button
		switch msg.Action {
			case network.NewOrder:
				if locOrdMat[floor][button] == 0 {
					drivers.ElevSetButtonLamp(drivers.TagElevLampType(button), floor, 1)
					if button == PanelButton {
						locOrdMat[floor][button]=1
						if IsLocOrdMatEmpty() {
							newOrder = true
							firstOrderFloor = floor
						}
						orderCount++
					} else {
						msg.Action = network.Tender
						msg.TenderVal = cost(floor, button)
						activeTenders[order] = TenderType{time.Now(), msg.TenderVal}
						//network.BroadcastOnNet(msg)  // Send tender for order on network
					}
				}			
			case network.DeleteOrder:
				delete(activeTenders, order)
				delete(lostTenders , order)
				drivers.ElevSetButtonLamp(drivers.TagElevLampType(button), floor, 0)
				if locOrdMat[floor][button] == 1 {
					locOrdMat[floor][button]=0
					orderCount--
				}	
			case network.Tender:
				if tender, ok := activeTenders[order]; ok { // Check if we already have a tender there
					if tender.val > msg.TenderVal {				// If our tender is worse than the one received, we delete it from active tenders and add it to lost tenders and let it go
						delete(activeTenders, order)
						lostTenders[order] = time.Now()
					} 
				} else {
					if tenderVal := cost(floor, button); tenderVal < msg.TenderVal {
						msg.TenderVal = tenderVal
						activeTenders[order] = TenderType{time.Now(), tenderVal}
						//network.BroadcastOnNet(msg)  // Send tender for order on network
					}
				}
			case network.AddOrder:
				delete(activeTenders, order)
				delete(lostTenders , order)
				if locOrdMat[floor][button] != 1 {
					drivers.ElevSetButtonLamp(drivers.TagElevLampType(button), floor, 1)
					locOrdMat[floor][button]=1
					if IsLocOrdMatEmpty() {
						newOrder = true
						firstOrderFloor = floor
					}
					orderCount++ 
				}
		}
	}
}

//Checks that the message is valid
func checkMsg(msg network.ButtonMsg) bool {
	switch msg.Action {
		case network.NewOrder, network.DeleteOrder, network.Tender, network.AddOrder:
			order 		  := msg.Order
			floor, button := order.Floor, order.Button
			//fmt.Printf("Floor %d, button %d \n", floor, button)
			if((floor != 0 && floor != Floors-1) || (floor == 0 && button != DownButton) || (floor == Floors-1 && button != UpButton)){
				if (floor>=0 && floor<Floors) && (button >=0 && button<3) && msg.TenderVal>=0 { 
				   fmt.Printf("Floor %d, button %d \n", floor, button)
					return true
				}
			}
	}
	return false
}

// Checks if the order matrix is empty
func IsLocOrdMatEmpty() bool {
	if orderCount == 0 {
		return true
	} else {
		return false
	}
}

//Returns current direction.
func ReturnDirection() Direction {
	fmt.Printf("dir top %d \n", direction)
	if atEndFloor {
		return direction
	} else {
		return direction
	}
}
func GetDir() Direction {
	if drivers.ElevGetFloorSensorSignal() == -1 {
		return direction
	}
	fmt.Printf("dir 1%d \n", direction)
	if orderCount == 0 { //If called and no orders exisits, just set direction to Stop and return Stop as FSM will go to Idle.
		return Stop
	}
	var ordersAtCur [3]bool //	Holds all orders on the current floor
	var ordersInDir [2]bool // [0] is true if there are orders further up, [1] is true if there is any up
	var currDir int     // Varable to hold the current direction to be used in orderInDir. 0 for up and 1 for down.
	var msg network.ButtonMsg
	msg.Action = network.DeleteOrder
	for i := range locOrdMat {
		for j := range locOrdMat[i] {
			if locOrdMat[i][j] == 1 {
				if i == prevFloor { //check for orders at current floor
					ordersAtCur[j] = true
				} else if i > prevFloor { // check for orders upwards
					ordersInDir[UpButton] = true
				} else if i < prevFloor { // check for orders downwards
					ordersInDir[1] = true
				}
			}
		}
	}
	switch direction {
		case Up:
			currDir = UpButton
		case Down:
			currDir = DownButton
	}
	fmt.Printf("dir2 %d \n", direction)
	if ordersAtCur[currDir] || ordersAtCur[2] { //Just stay put if there is an order at current floor from the panel or from outside in the same direction as travel
		/*
		msg.Order=network.OrderType{currDir, prevFloor}
		orderHandler(msg)
		msg.Order=network.OrderType{PanelButton, prevFloor}
		orderHandler(msg)
		*/
		fmt.Printf("dir3 %d \n", direction)
		return Stop
	} else if ordersInDir[currDir] { //Return current direction if there is an order in that direction
			fmt.Printf("dir 6%d \n", direction)
		return direction
	} else if ordersAtCur[currDir+int(direction)] { //Just stay put if there is an order at current flor in opposite direction
		firstOrderFloor = prevFloor
		direction = -1 * direction
		/*
		msg.Order=network.OrderType{currDir, prevFloor}
		orderHandler(msg)
		msg.Order=network.OrderType{PanelButton, prevFloor}
		orderHandler(msg)
		msg.Order=network.OrderType{currDir+int(direction), prevFloor}
		orderHandler(msg)
		*/
			fmt.Printf("dir 4%d \n", direction)
		return Stop
	} else if ordersInDir[currDir+int(direction)] { //Go in opposit direction if there is an order there there
		direction = -1 * direction
		return direction
	}
	fmt.Printf("dir 5%d \n", direction)
	return direction 	// Stay put if the logic above fails (Yeah, right...)
}
