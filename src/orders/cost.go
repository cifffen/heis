package orders

import (
        //"fmt"
		//"../drivers"
		//"../orders"
		"math"
)

const floorValue = 1 //Time to travel from one floor to the next
const waitValue = 3 //Time spent for each stop
const directionValue = 20 //Extra cost if the order is in a conflicting direction


func cost(orderFloor int, orderType int) (cost int) {
	//directionCost, waitCost, travelCost := 0, 0, 0
	if IsLocOrdMatEmpty(){
		cost = getTravelCost(orderFloor)
		return
	} else{
		cost = getTravelCost(orderFloor)
		cost += getWaitCost(orderFloor, prevFloor)
		cost += getDirectionCost(orderFloor, prevFloor)
		return
	}
}

func getTravelCost(orderFloor int) (travelCost int) {
		travelDistance := prevFloor - orderFloor
        travelCost = int(math.Abs(float64(travelDistance))*float64(floorValue))
        return
}

func getWaitCost(orderFloor int, orderType int)(waitCost int) {
	waitCount := 0
	for i:=0;i<Floors;i++{
		for j:=0;j<3;j++{
			if(locOrdMat[i][j]==1){
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

func getDirectionCost(orderFloor int, orderType int) (directionCost int) {
	if((orderType == 0 && direction == Up)||(orderType == 1 && direction == Down)){
		directionCost = 0
	} //else if(orderFloor == 
	return
}
