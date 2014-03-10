package main

import (
	"fmt"
	"../fsm"
	"os"
	//"../pp"
	"time"
)

func main(){
	/*
	ppSuccess := pp.ProcessPairs(os.Args)    // Launch process pairs
	if ppSuccess==0{   // If the elevator program has crahsed too many times, process pairs will shut down and the program stops.
		fmt.Printf("Too many reboots. Elevator shutting down. \n") 
		go fsm.InitElev()         // Try to init the system so we can stop at a floor in case the elevator was runing during the last crash.
		time.Sleep(time.Second*4) // Sleep for 4 seconds so we can get the init done.
		os.Exit(1)	
	} 
	*/
	if fsm.InitElev()==0{ // If we fail to init the IO we exit the program. Process pairs will eventually start the program back up again
		fmt.Printf("Error: Unable to initialize elevator hardware. Shuting down.\n")
		os.Exit(1) 
	}
	fsm.EventManager()  // Start the Event manager
	
}
