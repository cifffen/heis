package main

import (
	"../drivers"
	"../fsm"
)
func main() int{
	if fsm.InitElev()==0{
		fmt.Printf("Unable to initialize elevator hardware.\n")
		return 1
	}
	for {
		drivers.EventHandler()
	}
	
	
	return 0
}
