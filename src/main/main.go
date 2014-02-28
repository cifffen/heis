package main

import (
	"fmt"
	"../fsm"
)



func main(){
	if fsm.InitElev()==0{
		fmt.Printf("Unable to initialize elevator hardware.\n")
	}
	
	fsm.EventManager()
	
}
