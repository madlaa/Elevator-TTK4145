//Authors: Mads Laastad & Tommy Berntzen

package state

import(
	"fmt"
	"time"
	"math/rand"
	"../global"
	"../driver"
)

func ElevatorInit(buttonDetectChan chan global.Button, floorDetectChan chan int, setDirectionChan chan global.Direction, setButtonLightChan chan global.Button, clearButtonLightChan chan global.Button, openDoorChan chan bool, localOrderChan chan global.Button, localOrderRemoveChan chan global.Button, externalOrderChan chan global.Button, globalOrderRemoveChan chan global.Button, lastKnownStateChan chan global.Button, doorClosedChan chan bool){
	driver.DriverInit()
	elevatorPollingStartup(buttonDetectChan, floorDetectChan, setDirectionChan, setButtonLightChan, clearButtonLightChan, openDoorChan, doorClosedChan, localOrderChan, externalOrderChan, lastKnownStateChan)
	
	rand.Seed(time.Now().Unix())
	randValue := global.Direction(rand.Intn(2))
	var initFloor int
	initFinished := false
	for{
		select{
		case initFloor = <- floorDetectChan:
			if initFloor == -1{
				setDirectionChan <- randValue
			} else if initFloor != -1{
				setDirectionChan <- global.NONE
				openDoorChan <- true
				initFinished = true
			}
		}
		if initFinished {
			fmt.Println("Initialization complete!")
			break
		}
		time.Sleep(global.RECIVING_SLEEP_TIME)
	}
	go elevatorHandler(initFloor, floorDetectChan, openDoorChan, setDirectionChan, setButtonLightChan, clearButtonLightChan, localOrderChan, localOrderRemoveChan, lastKnownStateChan, doorClosedChan)
}
