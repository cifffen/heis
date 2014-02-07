package orderHandler



const Floors:=4

var orderMatrix[Floors][3] int
var direction Direction
var previousFloor int

func InitOrderHandler(floor int)(){
	direction = Stop
	previousFloor=floor
}

func GetOrders(eventChannel)(){  
	for{
		for (i:=0;i<Floors;i++){
			for( j:=0;j<3;j++){
				if((i!=0 && i!=Floors)||(i==0 && j!=1)||(i==Floors && j!=0)){   // Statement that makes sure that we don't check the Down button at the groud floor and 
					if ElevGetButtonSignal(j, i)==1 && orderMatrix[i][j]!=1{                 // the Up button at the top floor, as they don't exist.
						orderMatrix[i][j]=elevGetButtonSignal(j, i)
						temp := true
					}
				}
			}
		}
		if temp{
			eventChannel <- true
		}
	}
		
}
func DeleteFloorOrders(floor int)(){ 
	orderMatrix[floor][2]=0;
	switch direction{
		case Up:
			orderMatrix[floor][0]=0;  // Deletes the order from the Up button at the given floor.
		case Down:
			orderMatrix[floor][1]=0;  // Deletes the order from the Down button at the given floor.
	}
	
}

func AtOrder(eventChannel chan)(){ 
	for{
		floor := ElevGetFloorSensorSignal()
		if(floor!=-1){
			previousFloor=floor
			if(floor==Floors-1){                 // If the elevator is at the top floor the direction is changed as it can't go further Upwards.
				direction=Down 
			}
			else if(floor==0){                    // If the elevator is at the bottom floor the direction is changed as it can't go further Downwards.
				direction=Up
			}
			if(orderMatrix[floor][2]==1){                  // Stop if an order from the inside panel has been made at the current floor.
				DeleteFloorOrders(floor)
				eventChannel <- true
			}
			else if(direction==Up && orderMatrix[floor][0]==1){                // Stop if an order from the Up button at the current floor has been made and the elevator is going Up.
				DeleteFloorOrders(floor)
				eventChannel <- true
			}
			else if(direction==Down && orderMatrix[floor][1]==1){             // Stop if an order from the Down button at the current florr has been made and the elevator is going Down.
				DeleteFloorOrders(floor)
				eventChannel <- true
			}
		}
	}
}

func GetDirection() Direction{

	switch direction{
		case Up:
			for(i=previousFloor; i<Floors;i++){
				if((orderMatrix[i][0] || orderMatrix[i][2]) && previousFloor!=Floors-1){ // Go further Up if an order is made from the panel or Up button
					return Up													 // to a floor higher Up than the current floor.
				}
			}
			for(i=previousFloor; i>=0; i--){
				if(orderMatrix[i][1] || orderMatrix[i][2]){    // If there are no orders further Up, go Down if there are any orders made there.
					return Down
				}
			}return Stop;
		
		case Down:
			for(i=previousFloor; i>=0;i--){
				if((orderMatrix[i][1] || orderMatrix[i][2])&& previousFloor!=0){      // Go further Down if an order is made from the panel or Down button 
					return Down											// to a floor lower than the current floor.
				}
			}
			for(i=previousFloor; i<Floors; i++){  // If there are no orders further Down, go Up if there are any orders made there.
				if(orderMatrix[i][0] || orderMatrix[i][2]){
					return Up
				}
			}return Stop
		case Stop:                    
			if(firstOrderFloor!= -1){
				if(firstOrderFloor>previousFloor){
					return Up
				}
				else if(firstOrderFloor<previousFloor){
					return Down;
				}
				else if(firstOrderFloor==previousFloor){
					return Stop;
				}
			}
	}
}
