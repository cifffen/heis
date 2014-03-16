package orders

import (
      "fmt"
	  "math"
)

const floorValue 	 = 3 	//Time to travel from one floor to the next
const waitValue      = 3 	//Time spent for each stop
const directionValue = 20 	//Extra cost if the order is in a conflicting direction


func cost(orderFloor int, orderType int, locOrdMat [Floors][Buttons] int, prevFloor int, direction Direction) (cost int) {
	if isMatrixEmpty(locOrdMat){
		cost = getTravelCost(orderFloor, prevFloor)
		return
	} else{
		cost = getTravelCost(orderFloor, prevFloor)
		cost += getWaitCost(orderFloor, orderType, locOrdMat, prevFloor, direction)
		cost += getDirectionCost(orderFloor, orderType, direction, prevFloor, locOrdMat)

		return
	}
}

func getTravelCost(orderFloor int, prevFloor int) (travelCost int) {
		travelDistance := prevFloor - orderFloor
        travelCost = int(math.Abs(float64(travelDistance))*float64(floorValue))
        return
}

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

func getDirectionCost(orderFloor int, orderType int, direction Direction, prevFloor int, locOrdMat [Floors][Buttons] int) (directionCost int) {
	if(direction == Up){
	    if(onlyOppositeDirection(orderType, locOrdMat)){
	        directionCost=directionValue
		    return
		}
		if((orderType != 1 || (orderType == 1 && orderFloor == Floors-1)) && orderFloor > prevFloor){
			directionCost = 0
			return
		}else{
			directionCost = directionValue
			return
		}
	}else{
		if(onlyOppositeDirection(orderType, locOrdMat)){
	        directionCost=directionValue
		    return
		}
		if((orderType != 0 || (orderType == 0 && orderFloor == 0)) && orderFloor < prevFloor){
			directionCost = 0
			return
		}else{
			directionCost = directionValue
      	return
      }
   }
    fmt.Printf("Failed to get directionCost \n")
	directionCost = 0
	return
}

func onlyOppositeDirection(orderType int, locOrdMat [Floors][Buttons] int) bool {
	for i := range locOrdMat {
	    if locOrdMat[i][orderType]==1 {
	        return false
        }
    }
    return true
}
