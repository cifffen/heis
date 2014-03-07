package orders

import (
        //"fmt"
		//"../drivers"
		//"../orders"
		"math"
)

const floorValue 	 = 1 	//Time to travel from one floor to the next
const waitValue      = 3 	//Time spent for each stop
const directionValue = 20 	//Extra cost if the order is in a conflicting direction


func cost(orderFloor int, orderType int) (cost int) {
	if IsLocOrdMatEmpty(){
		cost = getTravelCost(orderFloor)
		fmt.Printf("Cost:%d\n", cost)
		return
	} else{
		cost = getTravelCost(orderFloor)
		cost += getWaitCost(orderFloor, orderType)
		cost += getDirectionCost(orderFloor, orderType)
		fmt.Printf("Cost:%d\n", cost)
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
	for i :=0 range locOrdMat{
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

func getDirectionCost(orderFloor int, orderType int) (directionCost int) {
	if((orderType == UpButton && direction == Up)||(orderType == DownButton && direction == Down)){
		directionCost = 0
	} //else if(orderFloor == 
	return
}
