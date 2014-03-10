package pp

import (
	"net"
	"os"
	"fmt"
	"time"
	"os/exec"
	"log"
)

const ProPairsPort  = ":1989" // Port used by processpairs
const BroadcastRate = 50      //How often you broadcast to the slave, in milliseconds
const HeartBeat 	= 400	  // Time check for a heartbeat [ms]
func StartSlave()() {
	cmd := exec.Command("mate-terminal", "-x", "./../main/main", "Slave") // Start a new program with the argument Slave so the program knows what it is.
	fmt.Printf("Slave started\n")	
	err := cmd.Start()  
	if err != nil {
		log.Printf("Error: %v", err)
	}	
}	

func UdpListenToMaster(number chan<- int, sock *net.UDPConn)() {	
	addr, err := net.ResolveUDPAddr("udp", ProPairsPort)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	} 
	sock, err = net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second*4)
	}
	fmt.Println("Connected")
	for {
		buf :=make([]byte,1024)
		_, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			time.Sleep(time.Second*4)
			fmt.Println(err)	
		}
		number<-int(buf[0])
	}
}

func UdpHeartBeat(number int)(){
	for {
		select {
			case <-time.After(time.Millisecond*BroadcastRate):
				con,err:= net.Dial("udp", "localhost"+ProPairsPort)
				if err != nil {
					log.Printf("Error: %v ", err)
				}
				buf :=[]byte(string(number))
				_, err = con.Write(buf)
				if err != nil {
					log.Printf("Error: %v ", err)
				}
		}
	}
}

func ProcessPairs(args []string) int {
	if len(os.Args)==1{
		StartSlave()
		go UdpHeartBeat(0)
	} else if os.Args[1] == "Slave" {
		ticker := time.NewTicker(time.Second*1)
		numberChan := make(chan int)
		var sock *net.UDPConn
		go UdpListenToMaster(numberChan, sock)
		var num int	
		for {
			select {
				case num = <-numberChan:
					ticker.Stop()
					ticker = time.NewTicker(time.MilliSecond*HeartBeat)
				case <-ticker.C:
					sock.Close()
					num++
					go UdpHeartBeat(num)
					if num >=3 {
						return 0
					} else {
						StartSlave()
						return 1
					}
			}
		}
	} 
	fmt.Printf("Error: Wrong input. Running without processparis.")
	return 1 
}





