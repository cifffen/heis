package fsm

import (
	"../drivers"
	"../orderMod"
	"time"
	"fmt"
)

type(
	Event int
	State int
)

const(
	OrderReached Event =iota
    	TimerFinished
	NewOrder
)
const (
	Idle State =iota
	Running
	AtFloor
)
var state State
var DoorTimer <-chan time.Time 
func InitElev() int{
	if drivers.ElevInit() ==0 {  //IO init failed
		return 0
	} else {
		drivers.ElevSetSpeed(int(orderMod.Down))
		for drivers.ElevGetFloorSensorSignal() ==-1 {
		}
		orderMod.InitOrderMod(drivers.ElevGetFloorSensorSignal())
		drivers.ElevSetSpeed(int(orderMod.Stop))
		state = Idle
		fmt.Printf("Initialized\n")
		return 1
	}
}

func StateMachine(event Event)(){
	switch state{
		case Idle:
			switch event{
				case NewOrder:
					drivers.ElevSetSpeed(int(orderMod.GetDir()))
					fmt.Printf("New order\n")
					state = Running
				default:
			}
		case Running:
			switch event{
				case OrderReached:
					drivers.ElevSetSpeed(int(orderMod.Stop))
					DoorTimer = time.After(time.Second*3)
					state=AtFloor
				default:
			}
		case AtFloor:
			switch event{
				case TimerFinished:
					direction := orderMod.GetDir()
					if(direction==orderMod.Stop){
						state=Idle
					} else {
						state=Running
						drivers.ElevSetSpeed(int(direction))
					}
				default:
							
			}
		}
}
