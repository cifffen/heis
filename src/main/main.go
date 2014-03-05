package main

import (
	"fmt"
	"../fsm"
	"os"
	//"time"
	//"../pp"
)



func main(){
	//ppSuccess := pp.ProcessPairs(os.Args)
	if fsm.InitElev()==0{
		fmt.Printf("Error: Unable to initialize elevator hardware. Shuting down.\n")
		os.Exit(1)
	}
	/*
	if ppSuccess==0{
		fmt.Printf("Too many reboots. Elevator shutting down. \n") 
		go fsm.EventManager()
		time.Sleep(time.Second*4)
		os.Exit(1)	
	} 
	*/
	fsm.EventManager()
	
}
