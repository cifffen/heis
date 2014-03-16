package pp
/*
Process pairs module. Will reboot the program "numberOfReboots" times if the porgram should shut down for 
some reason.
*/
import (
	"net"
	"os"
	"fmt"
	"time"
	"os/exec"
	"log"
	"strconv"
)
const numberOfReboots 	= 3		 
const ProPairsPort  	= ":1989" // Port used by processpairs
const BroadcastRate 	= 50      //How often you broadcast to the slave, in milliseconds
const HeartBeat 		= 400	  // Time check for a heartbeat [ms]

func StartSlave(number int)() { // Start a new program with the argument Slave so the program knows what it is.
	cmd := exec.Command("mate-terminal", "-x", "./../main/main", "Slave", strconv.Itoa(number) ) 
	err := cmd.Start()  
	if err != nil {
		log.Printf("Error: %v", err)
		fmt.Printf("Could not start slave. \n")
		return
	}	
	fmt.Printf("Slave number %d started\n", number)	
}	
// Listen to the heartbeat from the master. If an error occurs here we have to shut down to prevent multiple programs runnings at once.
func UdpListenToMaster(number chan<- int, sock **net.UDPConn)() {	
	addr, err := net.ResolveUDPAddr("udp", "localhost"+ProPairsPort)
	if err != nil {							
		log.Printf("Error: %v", err)
		fmt.Printf("Shuting down slave.\n")
		time.Sleep(time.Second*2)
		os.Exit(1)
	} 
	*sock, err = net.ListenUDP("udp", addr)
	if err != nil {							
		log.Printf("Error: %v", err)
		fmt.Printf("Shuting down slave.\n")
		time.Sleep(time.Second*2)
		os.Exit(1)
	}
	for {
		buf :=make([]byte,1024)
		_, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error: %v", err)	
			fmt.Printf("Shuting down slave.\n")
			time.Sleep(time.Second*2)
			os.Exit(1)
		}
		number<-int(buf[0])
	}
}

// Broadcast heartbeat so the slave on localhost. We shut down the elevator should an error occur here to prevent mulitple programs running at once.
func UdpHeartBeat(number int)(){
	for {
		select {
			case <-time.After(time.Millisecond*BroadcastRate):
				con,err:= net.Dial("udp", "localhost"+ProPairsPort)
				if err != nil {
					log.Printf("Error: %v ", err)
					fmt.Printf("Shuting down.\n")
					time.Sleep(time.Second*4)
					os.Exit(1)
				}
				buf :=[]byte(string(number))
				_, err = con.Write(buf)
				if err != nil {
					log.Printf("Error: %v ", err)
					fmt.Printf("Shuting down.\n")
					time.Sleep(time.Second*4)
					os.Exit(1)
				}
		}
	}
}

func ProcessPairs(args []string) int {
	if len(os.Args)==1{
		time.Sleep(time.Second*3)
		StartSlave(0)
		go UdpHeartBeat(0)
		return 1
	} else if os.Args[1] == "Slave" && len(os.Args)==3 {
		ticker := time.NewTicker(time.Second*1)
		numberChan := make(chan int)
		var sock *net.UDPConn
		go UdpListenToMaster(numberChan, &sock)
		num, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Printf("Error: %v", err)        //If we get a wrong input when starting the slave
			fmt.Printf("Slave shutting down.\n") // we kill the slave so the program will end
			time.Sleep(time.Second*2)			// if the master dies.
			os.Exit(1) 	
		}				
		for {
			select {
				case <-numberChan:
					ticker.Stop()
					ticker = time.NewTicker(time.Millisecond*HeartBeat)
				case <-ticker.C:  // If we don't here from the master in a given time
					sock.Close()  // We close the socket so the next slave can use it
					time.Sleep(time.Millisecond*200) // Wait from 200 ms to be sure the socket is closed before we start a slave
					num++							 // Keep count of number of reboots
					if num >= numberOfReboots {					// Shut down if there were to many reboots
						return 0
					} else {
						go UdpHeartBeat(num)  		// Start broadcasting as master
						time.Sleep(time.Second*2)   // Give us a little time in case some weird stuff happens that makes a fork bomb.
						StartSlave(num)
						return 1
					}
			}
		}
	} 
	fmt.Printf("Error: Wrong input. Running without processparis.")
	return 1
}
