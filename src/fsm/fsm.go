package fsm

import (
	"../drivers"
	"../orderHandler"
)

type(
	Event int
	State int
	Direction int
)
const (
	Down Direction = -300
	Up = 300
	Stop = 0
)
const(
	OrderReached Event =iota
    TimerFinished
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
	orderHandler.InitOrderHandler(ElevGetFloorSensorSignal)
	drivers.ElevSetSpeed(Stop)
	state = Idle
	return 1
}
func FloorReach(event chan bool){
	if (floor == 1)
		event <- true
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
				case OrderReached:
				
				default:
			}
		case AtFloor:
			switch event{
				case TimerFinished:
				
				default:
			}
		}
}