
package orders

import (
	"../network"
	"time"
	"../types"	
	"../drivers"
	"fmt"
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

const TakeActiveTender	 = 1000 // Milliseconds
const TakeLostTender     = 20  // Seconds
const SamplingTime		 = 1   // Milliseconds
const Floors 			 = 4   // Number of floors
const Buttons			 = 3   // Number of buttons

func OrderHandler(orderReachedEvent chan<- bool, newOrderEvent chan<- bool, switchDirEvent chan<- Direction, noOrdersEvent chan<- bool) {
	var direction 		Direction 	// Keeps the last direction the elevator was heading. Can only be changed in atOrder() and GetDir()
	var prevDirection Direction
	var prevFloor 		int      	// Holds the previous floor the elevator past. Can only be changed at atOrder()
	var noOrders		bool		// True if we have no orders in locOrdMat
	var locOrdMat [Floors][3]int 	// Holds the orders that the elevator has accepted and will carry out
	var activeTenders map[types.OrderType] 	TenderType
	var lostTenders map[types.OrderType] 		time.Time
	
	//---- Start Init of variables and network--------//
	direction 		= Down
	noOrders		= true
	activeTenders 	= make(map[types.OrderType] TenderType)
	lostTenders 	= make(map[types.OrderType] time.Time)
	msgChan 		:= make(chan types.OrderMsg)
	go network.ListenOnNetwork(msgChan)
	//---- Init complete ------//
	for {
		select {
		case <-time.After(time.Millisecond * SamplingTime): 	//Only check for events bellow every Sampling time [ms]
			if newOrders, msgSlice := getOrders(&locOrdMat, activeTenders, lostTenders); newOrders {  // Check for new orders.
				for _, msg := range msgSlice {  // Go through all new orders and process them in msgHandler
					if newOrder:= msgHandler(msg, &locOrdMat, &activeTenders, &lostTenders, prevFloor, direction); newOrder{
						newOrderEvent <-true  	// New order from an empty order matrix has occured
						noOrders = false		// Must be set false as we now have an order in the order list
					}
				}
			}
			if matEmpty := isLocOrdMatEmpty(locOrdMat); matEmpty && !noOrders { // only launch the event once when we have no orders left
				noOrdersEvent <- true
				noOrders = true  	// We now have no orders in our order list
			}
			if orderReached, del, delOrders := atOrder(locOrdMat, direction, &prevFloor); orderReached { // Launch event if we reach an order	
				if del {							// If we have orders to delete, delete them
					for _ , msg := range delOrders {
						if newOrder:= msgHandler(msg, &locOrdMat, &activeTenders, &lostTenders, prevFloor, direction); newOrder{
						   newOrderEvent <-true  	// New order from an empty order matrix has occured
						   noOrders = false		// Must be set false as we now have an order in the order list
					   }
					}
				}
				orderReachedEvent <- true
			}
			if currDir := getDir(direction, prevFloor, locOrdMat); currDir !=  prevDirection{
				switchDirEvent <- currDir
				prevDirection = currDir
				fmt.Printf("direction here: %d \n", direction)
				if currDir != Stop {
					direction = currDir
				}
			}

			if tenderAction , tenderOrders := checkTenderMaps(activeTenders, lostTenders); tenderAction{ // If some times for the tenders on the tender lists have run out -
				for _, msg := range tenderOrders {
					if newOrder:= msgHandler(msg, &locOrdMat, &activeTenders, &lostTenders, prevFloor, direction); newOrder{ // we let msgHandler handle the messages/orders. Add them if they are from active tenders or start a new tender session over the network if they are from lost tenders
						newOrderEvent <-true  	// New order from an empty order matrix has occured
						noOrders = false		// Must be set false as we now have an order in the order list
					}  
				}															   
			}
			
		case msg:= <-msgChan:  // Received message on the network
			if newOrder := msgHandler(msg, &locOrdMat, &activeTenders, &lostTenders, prevFloor, direction); newOrder{
				newOrderEvent <-true  // New order from an empty order matrix has occured
				noOrders = false		// Must be set false as we now have an order in the order list
			}
		}
	}
}

//Handles orders both locally and over the network
func msgHandler(msg types.OrderMsg, locOrdMat *[Floors][Buttons] int, aTenders *map[types.OrderType] TenderType, lTenders *map[types.OrderType] time.Time, prevFloor int, direction Direction)(newOrder bool) {
	newOrder = false
	if checkMsg(msg) {     // Check if message is valid
		order := msg.Order
		floor, button := order.Floor, order.Button
		switch msg.Action {
			case types.NewOrder:
				if (*locOrdMat)[floor][button] == 0 { // Check if we have an order there already
					drivers.ElevSetButtonLamp(drivers.TagElevLampType(button), floor, 1) 
					if button == PanelButton {        		// If the order is from the inside panel we add the order directly	
						if isLocOrdMatEmpty(*locOrdMat) {   // Launch new order event if the order list is empty
							newOrder = true
						}
						(*locOrdMat)[floor][button]=1	
					} else {								// If the order is from the direction panel, -
						msg.Action = types.Tender					// we calculate our tender, add  to active tenders list and - 
						msg.TenderVal = cost(floor, button, *locOrdMat, prevFloor, direction) // start a tender session on the network. Lowest tender "wins" the order.
						(*aTenders)[order] = TenderType{time.Now(), msg.TenderVal}  // Add tender to active tenders
						network.BroadcastOnNet(msg)  
					}
				}			
			case types.DeleteOrder:            // Delete order
				delete(*aTenders, order)
				delete(*lTenders , order)
				drivers.ElevSetButtonLamp(drivers.TagElevLampType(button), floor, 0)
				if (*locOrdMat)[floor][button] == 1 { 	// If it is "our" order -
					(*locOrdMat)[floor][button]=0	   	// we delete it and -
					msg.Action = types.DeleteOrder			// tell the other elevators on the network to delete it from there tender lists
					network.BroadcastOnNet(msg)
				}	
			case types.Tender:
			   drivers.ElevSetButtonLamp(drivers.TagElevLampType(button), floor, 1)
				if tender, ok := (*aTenders)[order]; ok { // Check if we already have a tender there
					if tender.val > msg.TenderVal {		// If our tender is worse than the one received -
						delete(*aTenders, order)		// we delete it from active tenders -
						(*lTenders)[order] = time.Now()	// and add it to lost tenders 
					} 
				} else {																// If we don't already have a tender at that order, 
					if tenderVal := cost(floor, button, *locOrdMat, prevFloor, direction); tenderVal < msg.TenderVal {	// we calculate a tender for it and check if ours is better than there's
						msg.TenderVal = tenderVal										// If our tender is better -
						(*aTenders)[order] = TenderType{time.Now(), tenderVal}			// we add it to active tenders
						network.BroadcastOnNet(msg)  									//we send it out on the network
					} else {
						(*lTenders)[order] = time.Now() 			// If our tenders is worse, we add it to lost tenders
					}
				}
			case types.AddOrder:	// Directly add an order from active tenders if the time has run out
				delete(*aTenders, order)
				delete(*lTenders , order)
				if (*locOrdMat)[floor][button]  == 0 { // Make sure we already don't have that order (should not happen)
					drivers.ElevSetButtonLamp(drivers.TagElevLampType(button), floor, 1)  // Set order lamp	
					if isLocOrdMatEmpty(*locOrdMat) { // Launch new order event if the order list is empty
						newOrder = true
						fmt.Printf("New order \n")
					}
					(*locOrdMat)[floor][button] = 1	
				}
		}
	}
	return
}
// Check the tender maps. 
func checkTenderMaps(aTenders map[types.OrderType] TenderType, lTenders map[types.OrderType] time.Time)(tenderAction bool, tenderOrders []types.OrderMsg){
	var msg types.OrderMsg
	tenderAction = false
	for order, tenderTime := range lTenders {   
		if time.Since(tenderTime) > time.Second*TakeLostTender{  	// If the time for the lost tender has run out
				msg.Order     = order									// we delete the order from our lists
				msg.Action    = types.DeleteOrder							// and start a new tender session on the network for the order
				tenderOrders  = append(tenderOrders,msg)
				msg.Action 	  = types.NewOrder
				tenderOrders  = append(tenderOrders,msg)
				tenderAction  = true
		}
	}
	for order, value := range aTenders {
		if time.Since(value.time) > time.Millisecond*TakeActiveTender{  // If the time has passed, we have won the tender and can add it to our order list
				msg.Order     = order
				msg.Action    = types.AddOrder
				tenderOrders  = append(tenderOrders,msg)
				tenderAction  = true
		}
	}
	return
}

// Check if the elevator should stop at a floor it passes
func atOrder(locOrdMat[Floors][Buttons] int, prevDir Direction, prevFloor *int) (orderReached bool, del bool, delOrders []types.OrderMsg) {
	floor := drivers.ElevGetFloorSensorSignal()
	orderReached = false
	del = false
	if floor != -1 {
		*prevFloor = floor
		drivers.ElevSetFloorIndicator(floor) 	//Set floor indicator
		/*
		if isLocOrdMatEmpty(locOrdMat){  		// Stop at this floor if we have no orders in our order list (should not happen)
			orderReached = true
			return
		}
		*/
		var msg types.OrderMsg
		msg.Action = types.DeleteOrder
		ordersAtCur, ordersInDir := checkForOrders(locOrdMat, *prevFloor)
		var currDir int     // Varable to hold the current direction to be used in orderInDir. 0 for up and 1 for down.
		switch prevDir {
			case Up:
				currDir = UpButton
			case Down:
				currDir = DownButton
		}
		if  ordersAtCur[PanelButton] { //Stop here and delete order if there is an order at current floor from the panel button 
			order := types.OrderType{PanelButton, *prevFloor}
			msg.Order = order
			delOrders = append(delOrders, msg)
			del = true			// Mark that we have orders to delete
			orderReached = true
		}
		if  ordersAtCur[currDir] { //Stop here and delete order if there is an order at current floor from the the direction button outside in the same direction
			order := types.OrderType{currDir, *prevFloor}
			msg.Order = order
			delOrders = append(delOrders, msg)
			del = true			// Mark that we have orders to delete
			orderReached = true
		} 
		if ordersAtCur[currDir+int(prevDir)] && !ordersInDir[currDir]{ //If we have no further orders in the current direction and an order in the opposite. Stop here and delete the order
			order := types.OrderType{currDir+int(prevDir), *prevFloor}
			msg.Order = order
			delOrders = append(delOrders, msg)
			del = true			// Mark that we have orders to delete
			orderReached = true
		}
	}
	return 
}

//Checks that the message is valid
func checkMsg(msg types.OrderMsg) bool {
	switch msg.Action {
		case types.NewOrder, types.DeleteOrder, types.Tender, types.AddOrder:
			order 		  := msg.Order
			floor, button := order.Floor, order.Button
			if((floor != 0 && floor != Floors-1) || (floor == 0 && button != DownButton) || (floor == Floors-1 && button != UpButton)){
				if (floor>=0 && floor<Floors) && (button >=0 && button<Buttons) && msg.TenderVal>=0 { 
					return true
				}
			}
	}
	return false
}

func isLocOrdMatEmpty(locOrdMat [Floors][Buttons] int) bool {
	for i := range locOrdMat {
		for _, order := range locOrdMat[i] {
			if order == 1 {
				return false
			}
		}
	}
	return true
}

// Check for orders
func getOrders(locOrdMat *[Floors][Buttons] int, aTenders map[types.OrderType] TenderType, lTenders map[types.OrderType] time.Time )(newOrders bool, orders []types.OrderMsg ) {
	newOrders = false
	var msg types.OrderMsg
	msg.Action = types.NewOrder
	for i := range *locOrdMat {
		for j := range (*locOrdMat)[i] {
			if (i != 0 && i != Floors-1) || (i == 0 && j != 1) || (i == Floors-1 && j != 0) { // Statement that makes sure that we don't check the Down button at the groud floor and the Up button at the top floor, as they don't exist.
				if drivers.ElevGetButtonSignal(j, i) == 1 && (*locOrdMat)[i][j] == 0 {
					order := types.OrderType{j, i}
					_, lostOk   := lTenders[order]; 
					_, activeOk := aTenders[order]; 
					if !lostOk && !activeOk{ 	//Check that those order are not already active on the network, either as an active- or lost tender
						newOrders = true
						msg.Order = order
						orders = append(orders, msg)
					}
				}
			}
		}
	}
	return
}

func getDir(direction Direction, prevFloor int, locOrdMat[Floors][Buttons] int) Direction {
	if isLocOrdMatEmpty(locOrdMat){
		return Stop
	} else if prevFloor == Floors-1 {
		return Down
	} else if prevFloor == 0{
		return Up
	}
	ordersAtCur, ordersInDir := checkForOrders(locOrdMat, prevFloor)
	var currDir int     // Varable to hold the current direction to be used in orderInDir. 0 for up and 1 for down.
	switch direction {
		case Up:
			currDir = UpButton
		case Down:
			currDir = DownButton
	}
	if ordersAtCur[currDir] || ordersAtCur[2] { //Just stay put if there is an order at current floor from the panel or from outside in the same direction as travel
	   if drivers.ElevGetFloorSensorSignal() != -1{
		   return Stop
		 } else {
		    return direction
		 }
	} else if ordersInDir[currDir] { //Return current direction if there is an order in that direction
		return direction
	} else if ordersAtCur[currDir+int(direction)] { //Just stay put if there is an order at current flor in opposite direction
	   if drivers.ElevGetFloorSensorSignal() != -1{
		   return Stop
		 } else {
		    return -1*direction
		 }
	} else if ordersInDir[currDir+int(direction)] { //Go in opposit direction if there is an order there 
		direction = -1 * direction
		fmt.Printf("dir: %d \n", currDir+int(direction))
		return direction
	}
	return direction 	// Stay put if the logic above fails (Yeah, right...)
}

// Return all orders on the current floor in ordersAtCurFloor. OrderInDir's elements will be true if there is an order further in that direction. [0] is up, [1] is down
func checkForOrders(locOrdMat[Floors][Buttons] int, prevFloor int)(ordersAtCurFloor[Buttons] bool, ordersInDir[2] bool) {
	for i := range locOrdMat {
		for j := range locOrdMat[i] {
			if locOrdMat[i][j] == 1 {
				if i == prevFloor { 		//check for orders at current floor
					ordersAtCurFloor[j] = true
				} else if i > prevFloor { 	// check for orders upwards
					ordersInDir[UpButton] = true
				} else if i < prevFloor { 	// check for orders downwards
					ordersInDir[DownButton] = true
				}
			}
		}
	}
	return
}
	
