package fsm

import (
	"../drivers"
)

type(
	Event int
	State int
)
const (
	Down = -300
	Up = 300
	Stop = 0
)
const(
	FloorReached Event =iota
    TimerOut
	NewOrder
)
const (
	Idle State =iota
	Runing
	AtFloor
)
var state State

func InitElev() int{
	if drivers.ElevInit() ==0 {  //IO init failed
		return 0
	}else{
	drivers.ElevSetSpeed(Down)
	for ElevGetFloorSensorSignal() !=1 {
	}
	drivers.ElevSetSpeed(Stop)
	state = Idle
	return 1
}

func EventHandler()(){
	
}
func StateMachine(event Event)(){
	switch state{
		case Idle:
			switch event{
				case NewOrder:
				
				default:
			}
		case Runing:
			switch event{
				case atOrder:
				
				default:
			}
		case AtFloor:
			switch event{
				case TimerOut:
				
				default:
			}
		}
}