
package drivers

/*
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import "C"


type TagElevLampType int

const(
	ButtonCallUp TagElevLampType = iota
	ButtonCallDown
	ButtonCommand
)

func ElevInit() int{
	return int(C.elev_init())
}

func ElevSetSpeed(speed int){
	C.elev_set_speed(C.int(speed))
}

func ElevGetFloorSensorSignal() int{
	return int(C.elev_get_floor_sensor_signal())
}

func ElevGetButtonSignal(button int, floor int) int{
	return int(C.elev_get_button_signal(C.int(button),C.int(floor)))
}

func ElevGetStopSignal() int{
	return int(C.elev_get_stop_signal())
}

func ElevSetFloorIndicator(floor int){
	C.elev_set_floor_indicator(C.int(floor))
}

func ElevSetButtonLamp(button TagElevLampType, floor int, value int){
	C.elev_set_button_lamp(C.int(button),C.int(floor),C.int(value))
}

func ElevSetStopLamp(value int){
	C.elev_set_stop_lamp(C.int(value))
}

func ElevSetDoorOpenLamp(value int){
	C.elev_set_door_open_lamp(C.int(value))
}


