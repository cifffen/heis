package orders

import (
      "fmt"
	  "math"
)

const floorValue 	 = 4 	//Time to travel from one floor to the next
const waitValue      = 3 	//Time spent for each stop
const directionValue = 20 	//Extra cost if the order is in a conflicting direction
const twoDirSameFloorValue = 7 //Extra cost if the floor already has an order on the floor, in the opposite direction

//Calculates a cost for orders recieved from outside elevator
func cost(orderFloor int, orderType int, locOrdMat [Floors][Buttons] int, prevFloor int, direction Direction) (cost int) {
	if isMatrixEmpty(locOrdMat){ //If the elevator is not running the only concern is the distance from the order
		cost = getTravelCost(orderFloor, prevFloor)
		return
	} else{
		cost = getTravelCost(orderFloor, prevFloor)
		cost += getWaitCost(orderFloor, orderType, locOrdMat, prevFloor, direction)
		cost += getDirectionCost(orderFloor, orderType, direction, prevFloor, locOrdMat)
        cost += getTwoDirSameFloorCost(orderFloor, orderType, locOrdMat)
		return
	}
}

//Adds to the cost based on distance from current position.
func getTravelCost(orderFloor int, prevFloor int) (travelCost int) {
		travelDistance := prevFloor - orderFloor
        travelCost = int(math.Abs(float64(travelDistance))*float64(floorValue))
        return
}

//Adds to the cost based on how many stops are between the order and current position.
func getWaitCost(orderFloor int, orderType int, locOrdMat [Floors][Buttons] int, prevFloor int, direction Direction) (waitCost int) {
	waitCount := 0
	for i := range locOrdMat{
		for j := range locOrdMat[i]{
			if(locOrdMat[i][j]==1){
				if(direction == Up && j != DownButton && orderFloor > prevFloor && i < orderFloor){
					waitCount++
					break
				} else if(direction == Down && j != UpButton && orderFloor < prevFloor){
					waitCount++
					break
				}
			}
		}
	}
	waitCost = waitCount*waitValue
	return
}

//Adds to the cost if the order is conflicting with the elevators current direction.
func getDirectionCost(orderFloor int, orderType int, direction Direction, prevFloor int, locOrdMat [Floors][Buttons] int) (directionCost int) {
	if(direction == Up){
	    if(onlyOppositeDirection(orderType, locOrdMat)){ //If we only have orders down despite going up, 
	        directionCost=directionValue                 //we'll consider down as current direction.
		    return
		}
		if((orderType == 0 || (orderType == 1 && orderFloor == Floors-1)) && orderFloor > prevFloor){ 
			directionCost = 0                           // End floors are treated as current direction.
			return
		}else{
			directionCost = directionValue
			return
		}
	}else{
		if(onlyOppositeDirection(orderType, locOrdMat)){ //If we only have orders up despite going down,
	        directionCost=directionValue				 //we'll consider up as current direction.
		    return
		}
		if((orderType == 1 || (orderType == 0 && orderFloor == 0)) && orderFloor < prevFloor){ // End floors are treated 
			directionCost = 0																   // as current direction.
			return
		}else{
			directionCost = directionValue
      	return
      }
   }
    fmt.Printf("Error: Failed to get directionCost \n")
	directionCost = 0
	return
}

//Checks whether or not our local matrix contains only orders in the opposite direction.
func onlyOppositeDirection(orderType int, locOrdMat [Floors][Buttons] int) bool {
	for i := range locOrdMat {
	    if locOrdMat[i][orderType]==1 {
	        return false
        }
    }
    return true
}

//If a floor already has the an order on the floor in the opposite direction,
//this function will add to the cost and encourage another elevator to take the order.
func getTwoDirSameFloorCost(orderFloor int, orderType int, locOrdMat [Floors][Buttons] int) int {
    if orderType == 0 {
        if locOrdMat[orderFloor][1] == 1{
            return twoDirSameFloorValue
        }
    }
    if locOrdMat[orderFloor][0] == 1{
        return twoDirSameFloorValue
    }
    return 0
}
    
     
               
        
        
        
        

