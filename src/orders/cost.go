package orders

import (
        //"fmt"
		"../drivers"
		"../orderMod"
		"math"
)

const floorValue = 1 \\Time to travel from one floor to the next
const waitValue = 3 \\Time spent for each stop
const directionValue = 20 \\Extra cost if the order is in a conflicting direction


func CostFunc(orderFloor int, orderType int) (cost int) {
	directionCost, waitCost, travelCost := 0, 0, 0
	prevFloor := drivers.ElevGetFloorSensorSignal
	direction := orderMod.GetDirection
	if(orderMod.IsOrderMatrixEmpty){
		cost = getTravelCost(orderFloor, prevFloor)
		return
	} else{
		cost = getTravelCost
		cost += getWaitCost
		cost += getDirectionCost
		return
	}
}

func GetTravelCost(orderFloor int, prevFloor int) (travelCost int) {
		travelDistance := prevFloor - orderFloor
        travelCost = Abs(travelDistance)*floorValue
        return
}

func GetWaitCost(orderFloor int, orderType int, prevFloor int, direction Direction)(waitCost int) {
	waitCount := 0
	for i:=0;i<Floors;i++{
		for j:=0;j<3;j++{
			if(orderMod.orderMatrix[i][j]==1){
				if(direction == Up && j != 1 && orderFloor > prevFloor && i < orderFloor){
					waitCount++
					break
				}else if(direction == Down && j != 0 && orderFloor < prevFloor){
					waitCount++
					break
				}
			}
		}
	}
	waitCost = waitCount*waitValue
	return
}

func GetDirectionCost(orderFloor int, orderType int, direction Direction, prevFloor int) (directionCost int) {
	if((orderType == 0 && direction == Up)||(orderType== 1 && direction = Down)){
		directionCost = 0
	} else if(orderFloor == 
	
	
	        return
}