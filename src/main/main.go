package main

import (
	"fmt"
	"../fsm"
	"os"
)



func main(){
	if fsm.InitElev()==0{
		fmt.Printf("Unable to initialize elevator hardware.\n")
		os.Exit(0)
	}
	
	fsm.EventManager()
	
}
