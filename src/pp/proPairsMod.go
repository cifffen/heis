package pp

import (
	"net"
	"os"
	"fmt"
	"time"
	"os/exec"
	"log"
)

const ProPairsPort = ":1989"
const BroadcastRate = 50 //How often you broadcast to the slave, in milliseconds

func StartSlave()() {
	//cmd := exec.Command("cmd", "/C", "start", "processpairs.exe", "Slave")
	cmd := exec.Command("mate-terminal", "-x", "./../main/main", "Slave") 
	fmt.Printf("Slave started\n")	
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}	
}	

func UdpListenToMaster(number chan<- int, quit <-chan int)() {	
	addr, err := net.ResolveUDPAddr("udp", ProPairsPort)
	if err != nil {
		fmt.Println(err)
	}
	sock, err:= net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second*4)
	}
	fmt.Println("Connected")
	l:
		for {
			buf :=make([]byte,1024)
			_, _, err := sock.ReadFromUDP(buf)
			if err != nil {
				time.Sleep(time.Second*4)
				fmt.Println(err)	
			}
			select{
				case <-quit:
					fmt.Printf("Close \n")
					sock.Close()
					number<-10
					break l
				default:
					number<-int(buf[0])
			}
		}
}

func UdpHeartBeat(number int, )(){
	fmt.Printf("Heart Beat! \n")
	for {
		select{
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
		numOfRestarts := 0
		go UdpHeartBeat(numOfRestarts)
	}else if os.Args[1] == "Slave" {
		ticker := time.NewTicker(time.Second*1)
		numberChan := make(chan int)
		quitChan := make(chan int)
		go UdpListenToMaster(numberChan, quitChan)
		var num int	
		for{
			select{
				case num = <-numberChan:
					ticker.Stop()
					ticker = time.NewTicker(time.Second*1)
				case <-ticker.C:
					close(quitChan)
					num++
					go UdpHeartBeat(num)
					if num >=3 {
						return 0
					}else{
					StartSlave()
						return 1
					}
					
			}
		}
	}else{
		fmt.Printf("Error: Wrong input. Running without processparis.")
	}
	return 1 
}





