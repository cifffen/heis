package orders




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


// Check for orders
func getOrders() bool {
	firstOrderEvent := false
	for i := range OrderMatrix {
		for j := range OrderMatrix[i] {
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

// Check if the elevator should stop at a floor it passes
func atOrder() bool {
	floor := drivers.ElevGetFloorSensorSignal()
	if floor != -1 {
		prevFloor = floor
		drivers.ElevSetFloorIndicator(floor) //Set floor indicator
		if floor == Floors-1 {               // If the elevator is at the top floor the direction is changed as it can't go further Upwards.
			//direction = Down
			atEndFloor = true
		} else if floor == 0 { // If the elevator is at the bottom floor the direction is changed as it can't go further Downwards.
			//direction = Up
			atEndFloor = true
		} else {
			atEndFloor = false
		}
		dir:= ReturnDirection()
		if orderMatrix[floor][2] == 1 || firstOrderFloor == floor { // Stop if an order from the inside panel has been made at the current floor.
			firstOrderFloor = -1
			DeleteFloorOrders(floor)
			return true
		} else if (dir == Up && orderMatrix[floor][0] == 1) || (dir == Down && orderMatrix[floor][1] == 1) { // Stop if an order from the direction button at the current floor has been made and the elevator is going in that direction.
			DeleteFloorOrders(floor)
			return true
		}
	}
	return false
}


