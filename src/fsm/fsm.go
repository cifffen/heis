package fsm

import (
	"../drivers"
	"../orderMod"
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
	Runing
	AtFloor
)
var state State
var DoorTimer <-chan time.Time 
func InitElev() int{
	if drivers.ElevInit() ==0 {  //IO init failed
		return 0
	} else {
	drivers.ElevSetSpeed(int(orderMod.Down))
	for ElevGetFloorSensorSignal() !=1 {
	}
	orderMod.InitOrderMod(ElevGetFloorSensorSignal)
	drivers.ElevSetSpeed(int(orderMod.Stop))
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
					elevDrivers.ElevSetSpeed(int(orderMod.GetDir()))
					state = Runing
				default:
			}
		case Runing:
			switch event{
				case OrderReached:
					elevDrivers.ElevSetSpeed(int(orderMod.Stop))
					DoorTimer = time.After(time.Second*3)
					state=AtFloor
				default:
			}
		case AtFloor:
			switch event{
				case TimerFinished:
					direction = orderMod.GetDir()
					if(direction==orderMod.Stop){
						state=Idle
					} else {
						state=Runing
						elevDrivers.ElevSetSpeed(int(direction))
				default:
			}
		}
}